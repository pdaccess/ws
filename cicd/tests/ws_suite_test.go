package tests

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pdaccess/ws/internal/platform/servers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func waitForDB(ctx context.Context, connStr string) (*pgxpool.Pool, error) {
	maxRetries := 30
	retryInterval := 2 * time.Second

	var pool *pgxpool.Pool
	var err error

	for range maxRetries {
		poolConfig, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			time.Sleep(retryInterval)
			continue
		}

		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			time.Sleep(retryInterval)
			continue
		}

		if err = pool.Ping(ctx); err != nil {
			pool.Close()
			time.Sleep(retryInterval)
			continue
		}

		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, err)
}

var (
	pgContainer *postgres.PostgresContainer
	dbPool      *pgxpool.Pool
	server      *httptest.Server
)

func Test_Stater(t *testing.T) {
	RegisterFailHandler(Fail)
	defer GinkgoRecover()
	RunSpecs(t, "ws suite")
}

var _ = BeforeSuite(func() {
	ginkgo.By("starting pgvector container", func() {
		ctx := context.Background()

		pgContainer, err := postgres.Run(ctx,
			"pgvector/pgvector:pg17",
			postgres.WithDatabase("testdb"),
			postgres.WithUsername("testuser"),
			postgres.WithPassword("testpass"),
		)
		Expect(err).ShouldNot(HaveOccurred())

		host, err := pgContainer.Host(ctx)
		Expect(err).ShouldNot(HaveOccurred())

		port, err := pgContainer.MappedPort(ctx, "5432")
		Expect(err).ShouldNot(HaveOccurred())

		connStr := fmt.Sprintf("postgres://testuser:testpass@%s:%d/testdb?sslmode=disable", host, port.Int())

		ginkgo.By("waiting for database to be ready")
		dbPool, err = waitForDB(ctx, connStr)
		Expect(err).ShouldNot(HaveOccurred())

		_, err = dbPool.Exec(ctx, "CREATE EXTENSION IF NOT EXISTS vector")
		Expect(err).ShouldNot(HaveOccurred())

		ginkgo.DeferCleanup(func() {
			if dbPool != nil {
				dbPool.Close()
			}
			if pgContainer != nil {
				pgContainer.Terminate(context.Background())
			}
		})
	})

	ginkgo.By("starting test HTTP server", func() {
		server = httptest.NewServer(servers.NewHttpServer())
		ginkgo.DeferCleanup(func() {
			server.Close()
		})
	})
})
