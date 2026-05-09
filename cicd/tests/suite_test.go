package tests

import (
	"context"
	"fmt"
	"net"
	stdhttp "net/http"
	"os"
	"testing"
	"time"

	"github.com/pdaccess/commons/pkg/domain"
	"github.com/pdaccess/ws/cmd/app"
	pdhttp "github.com/pdaccess/ws/pkg/http"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	dbContainer *postgres.PostgresContainer
	baseURL     string
	client      *pdhttp.Client
	apiClient   *pdhttp.ClientWithResponses
	httpClient  *stdhttp.Client
	testToken   string
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
		defer GinkgoRecover()
		err := app.RunWebServiceServer()
		Expect(err).To(BeNil())
	}()

	baseURL = fmt.Sprintf("http://%s", addr)

	claims := &domain.PdaccessClaims{
		UserId: "test-user-id",
		Urk:    "test-user-urk",
		Role:   "admin",
		Realm:  "test-realm",
		Iat:    time.Now().Unix(),
		Exp:    time.Now().Add(time.Hour).Unix(),
	}
	testToken, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("test-secret"))

	client, err = pdhttp.NewClient(baseURL,
		pdhttp.WithRequestEditorFn(func(ctx context.Context, req *stdhttp.Request) error {
			req.Header.Set("Authorization", "Bearer "+testToken)
			return nil
		}),
	)
	Expect(err).ShouldNot(HaveOccurred())

	apiClient, err = pdhttp.NewClientWithResponses(baseURL,
		pdhttp.WithRequestEditorFn(func(ctx context.Context, req *stdhttp.Request) error {
			req.Header.Set("Authorization", "Bearer "+testToken)
			return nil
		}),
	)
	Expect(err).ShouldNot(HaveOccurred())

	httpClient = &stdhttp.Client{}

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

func GetClient() *pdhttp.Client {
	return client
}

func GetAPIClient() *pdhttp.ClientWithResponses {
	return apiClient
}

func GetHTTPClient() *stdhttp.Client {
	return httpClient
}

func GetTestToken() string {
	return testToken
}
