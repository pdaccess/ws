package tests

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/pdaccess/ws/cmd/app"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	dbContainer *postgres.PostgresContainer
	baseURL     string
)

func Test_API(t *testing.T) {
	RegisterFailHandler(Fail)
	defer GinkgoRecover()
	RunSpecs(t, "API Suite")
}

var _ = BeforeSuite(func() {
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"pgvector/pgvector:pg17",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		postgres.BasicWaitStrategies(),
	)
	Expect(err).ShouldNot(HaveOccurred())
	dbContainer = container

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	Expect(err).ShouldNot(HaveOccurred())

	os.Setenv("DATABASE_URL", connStr)
	os.Setenv("LOG_LEVEL", "debug")

	listener, err := net.Listen("tcp", "localhost:0")
	Expect(err).ShouldNot(HaveOccurred())

	addr := listener.Addr().String()
	listener.Close()

	os.Setenv("HTTP_LISTEN_ADDR", addr)

	go func() {
		err := app.RunWebServiceServer()
		if err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	baseURL = fmt.Sprintf("http://%s", addr)

	time.Sleep(500 * time.Millisecond)
})

var _ = AfterSuite(func() {
	if dbContainer != nil {
		dbContainer.Terminate(context.Background())
	}
})

func GetBaseURL() string {
	return baseURL
}
