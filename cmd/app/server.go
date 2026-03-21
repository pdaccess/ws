package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/pdaccess/ws/internal/adapters"
	"github.com/pdaccess/ws/internal/core/service"
	"github.com/pdaccess/ws/internal/database"
	"github.com/pdaccess/ws/internal/platform/servers"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "WebService server for managing inventory and user groups",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if err := RunWebServiceServer(); err != nil {
			panic(fmt.Errorf("server app: %w", err))
		}
	},
}

func RunWebServiceServer() error {
	config, err := ParseConfig()

	if err != nil {
		log.Err(err).Msg("read config")
		return fmt.Errorf("read config: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Info().Msgf("Server is starting")

	connStr := config.DatabaseConfig.Url

	db, err := database.New(connStr)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		return fmt.Errorf("database migrations failed: %w", err)
	}

	vg, err := adapters.NewVectorGenerator()
	if err != nil {
		return fmt.Errorf("vector generator: %w", err)
	}

	svc := service.New(db.InventoryRepo(), db.UserRepo(), db.ActivityRepo(), db.PasteRepo(), db.UserGroupRepo(), db.ServiceSettingsRepo(), vg)

	routers := servers.NewHttpServer(svc)
	server := &http.Server{Addr: config.HttpListenAddr, Handler: routers}

	errChan := make(chan error)

	go func() {
		log.Info().
			Str("port", config.HttpListenAddr).
			Msg("Starting http server...")
		listener, err := net.Listen("tcp", config.HttpListenAddr)
		if err != nil {
			errChan <- fmt.Errorf("can't listen http addr : %s  %w", config.HttpListenAddr, err)
			return
		}

		err = server.Serve(listener)
		if err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("Stopping server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return fmt.Errorf("[error] %w", err)
		}

	case err := <-errChan:
		return fmt.Errorf("returned error: %w", err)
	}

	return nil
}
