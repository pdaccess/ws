package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"git.h2hsecure.com/core/ws/internal/database"
	"git.h2hsecure.com/core/ws/internal/servers"
	"github.com/rs/zerolog/log"
)

func StartServer(config *ServerConfig, signalChan chan os.Signal) error {
	log.Info().Msgf("Server is starting")

	dbCfg := database.Config{
		Host:     config.DatabaseConfig.Host,
		Port:     config.DatabaseConfig.Port,
		Username: config.DatabaseConfig.Username,
		Password: config.DatabaseConfig.Password,
		DB:       config.DatabaseConfig.DB,
	}

	db, err := database.New(dbCfg)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		return fmt.Errorf("database migrations failed: %w", err)
	}

	routers := servers.NewHttpServer()
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
	case <-signalChan:
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
