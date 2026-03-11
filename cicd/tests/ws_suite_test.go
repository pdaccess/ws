package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.h2hsecure.com/core/ws/internal/handlers"
	"git.h2hsecure.com/core/ws/internal/handlers/external"
	"git.h2hsecure.com/core/ws/internal/servers"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	pgContainer *postgres.PostgresContainer
	dbPool      *pgxpool.Pool
	server      *httptest.Server
)

func Test_Stater(t *testing.T) {
	RegisterFailHandler(Fail)
	defer GinkgoRecover()
}

var _ = BeforeSuite(func() {
	ginkgo.By("starting pgvector container", func() {
		ctx := context.Background()

		pgContainer, err := postgres.RunContainer(ctx,
			testcontainers.WithImage("pgvector/pgvector:latest"),
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

		poolConfig, err := pgxpool.ParseConfig(connStr)
		Expect(err).ShouldNot(HaveOccurred())

		dbPool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		Expect(err).ShouldNot(HaveOccurred())

		err = dbPool.Ping(ctx)
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

var _ = Describe("OpenAPI Endpoints", func() {

	Describe("Service Endpoints", func() {
		Context("POST /api/v1/ws/service", func() {
			It("should return 200 for valid request", func() {
				groupID := uuid.MustParse("4ff2d0ee-5479-433e-84c0-f6f8f44a8f50")
				port := int32(22)
				req := external.NewServiceRequestObject{
					Body: &external.CreateServiceRequest{
						Name:           "test-service",
						GroupId:        groupID,
						AccessProtocol: external.AccessProtocolSsh,
						AuthProtocol:   external.AuthProtocolRadius,
						Vendor:         "Redhat",
						Host:           "demo.pdaccess.com",
						Port:           &port,
					},
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.NewService(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.NewService200JSONResponse{}))
			})

			It("should return 400 for invalid request", func() {
				req := external.NewServiceRequestObject{
					Body: &external.CreateServiceRequest{
						Name: "",
					},
				}
				handler := handlers.NewHttpHandler()
				_, err := handler.NewService(context.Background(), req)
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("POST /api/v1/ws/service/index", func() {
			It("should return 200 for search request", func() {
				limit := int32(10)
				offset := int32(0)
				req := external.SearchServiceRequestObject{
					Params: external.SearchServiceParams{
						Limit:  &limit,
						Offset: &offset,
					},
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.SearchService(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.SearchService200JSONResponse{}))
			})
		})

		Context("GET /api/v1/ws/service/{service}", func() {
			It("should handle service lookup", func() {
				req := external.ServiceByIdRequestObject{
					Service: "non-existent-id",
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.ServiceById(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				_, ok := resp.(*external.ServiceById404JSONResponse)
				if !ok {
					Expect(resp).Should(BeAssignableToTypeOf(&external.ServiceById200JSONResponse{}))
				}
			})
		})
	})

	Describe("Group Endpoints", func() {
		Context("POST /api/v1/ws/group", func() {
			It("should return 200 for valid group creation", func() {
				req := external.NewGroupRequestObject{
					Body: &external.CreateGroupRequest{
						Name:        "test-group",
						Description: "Test group description",
					},
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.NewGroup(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.NewGroup200JSONResponse{}))
			})
		})

		Context("POST /api/v1/ws/group/index", func() {
			It("should return 200 for group search", func() {
				limit := int32(10)
				offset := int32(0)
				req := external.SearchGroupRequestObject{
					Params: external.SearchGroupParams{
						Limit:  &limit,
						Offset: &offset,
					},
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.SearchGroup(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.SearchGroup200JSONResponse{}))
			})
		})

		Context("GET /api/v1/ws/group/{groupId}", func() {
			It("should handle group lookup", func() {
				invalidUUID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
				req := external.GetConfigurationRequestObject{
					GroupId: invalidUUID,
				}
				handler := handlers.NewHttpHandler()
				_, err := handler.GetConfiguration(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("Snippet Endpoints", func() {
		Context("POST /api/v1/ws/snippet", func() {
			It("should return 200 for valid snippet creation", func() {
				content := "test content"
				req := external.NewSnippetRequestObject{
					Body: &external.SnippetRequest{
						Content: &content,
					},
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.NewSnippet(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.NewSnippet200JSONResponse{}))
			})
		})

		Context("GET /api/v1/ws/snippet", func() {
			It("should return 200 for snippet listing", func() {
				limit := int32(10)
				offset := int32(0)
				req := external.UserSnippetsRequestObject{
					Params: external.UserSnippetsParams{
						Limit:  &limit,
						Offset: &offset,
					},
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.UserSnippets(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.UserSnippets200JSONResponse{}))
			})
		})
	})

	Describe("Config Endpoints", func() {
		Context("GET /api/v1/ws/config/{realm}/{configContext}", func() {
			It("should return 200 for valid config fetch", func() {
				req := external.FetchContextRequestObject{
					Realm:         "test-realm",
					ConfigContext: "test-context",
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.FetchContext(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.FetchContext200JSONResponse{}))
			})
		})

		Context("POST /api/v1/ws/config/{realm}/{configContext}", func() {
			It("should return 200 for valid config upsert", func() {
				req := external.UpsertContextRequestObject{
					Realm:         "test-realm",
					ConfigContext: "test-context",
					Body:          &external.ConfigContextRequest{},
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.UpsertContext(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.UpsertContext200JSONResponse{}))
			})
		})
	})

	Describe("Alarm Endpoints", func() {
		Context("GET /api/v1/ws/alarm", func() {
			It("should return 200 for alarm listing", func() {
				limit := int32(10)
				offset := int32(0)
				req := external.AlarmIndexRequestObject{
					Params: external.AlarmIndexParams{
						Limit:  &limit,
						Offset: &offset,
					},
				}
				handler := handlers.NewHttpHandler()
				resp, err := handler.AlarmIndex(context.Background(), req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp).Should(BeAssignableToTypeOf(&external.AlarmIndex200JSONResponse{}))
			})
		})
	})

	Describe("Database Integration", func() {
		Context("pgvector extension", func() {
			It("should have vector extension enabled", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				var extExists bool
				err := dbPool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'vector')").Scan(&extExists)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(extExists).Should(BeTrue())
			})

			It("should be able to create vector column", func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				_, err := dbPool.Exec(ctx, "CREATE TABLE IF NOT EXISTS test_vectors (id serial primary key, embedding vector(3))")
				Expect(err).ShouldNot(HaveOccurred())

				_, err = dbPool.Exec(ctx, "INSERT INTO test_vectors (embedding) VALUES ('[1,2,3]')")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Describe("HTTP Server Integration", func() {
		Context("Health endpoint", func() {
			It("should return 200 for health check", func() {
				resp, err := http.Get(server.URL + "/health")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(200))
			})
		})
	})
})
