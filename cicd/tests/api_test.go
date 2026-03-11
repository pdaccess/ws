package tests

import (
	"net/http/httptest"
	"testing"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"
	"github.com/pdaccess/ws/internal/platform/servers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var server *httptest.Server

func Test_API(t *testing.T) {
	RegisterFailHandler(Fail)
	defer GinkgoRecover()
	RunSpecs(t, "API Suite")
}

var _ = BeforeSuite(func() {
	server = httptest.NewServer(servers.NewHttpServer())
	DeferCleanup(func() {
		server.Close()
	})
})

var _ = Describe("Services API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler()
	})

	Context("POST /services", func() {
		It("should create service", func() {
			req := httptest.NewRequest("POST", "/v1/ws/services", nil)
			w := httptest.NewRecorder()
			handler.PostServices(w, req)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("GET /services/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/services/1", nil)
			w := httptest.NewRecorder()
			handler.GetServicesId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})
})

var _ = Describe("Groups API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler()
	})

	Context("POST /groups", func() {
		It("should create group", func() {
			req := httptest.NewRequest("POST", "/v1/ws/groups", nil)
			w := httptest.NewRecorder()
			handler.PostGroups(w, req)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("GET /groups/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/groups/1", nil)
			w := httptest.NewRecorder()
			handler.GetGroupsId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("GET /groups/{id}/members", func() {
		It("should return members", func() {
			req := httptest.NewRequest("GET", "/v1/ws/groups/1/members", nil)
			w := httptest.NewRecorder()
			handler.GetGroupsIdMembers(w, req, 1, external.GetGroupsIdMembersParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})
})

var _ = Describe("Policies API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler()
	})

	Context("GET /policies", func() {
		It("should list policies", func() {
			req := httptest.NewRequest("GET", "/v1/ws/policies", nil)
			w := httptest.NewRecorder()
			handler.GetPolicies(w, req, external.GetPoliciesParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("POST /policies", func() {
		It("should create policy", func() {
			req := httptest.NewRequest("POST", "/v1/ws/policies", nil)
			w := httptest.NewRecorder()
			handler.PostPolicies(w, req)
			Expect(w.Code).Should(Equal(201))
		})
	})
})

var _ = Describe("Admin API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler()
	})

	Context("GET /admin/users", func() {
		It("should list users", func() {
			req := httptest.NewRequest("GET", "/v1/ws/admin/users", nil)
			w := httptest.NewRecorder()
			handler.GetAdminUsers(w, req, external.GetAdminUsersParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("GET /admin/system-health", func() {
		It("should return health status", func() {
			req := httptest.NewRequest("GET", "/v1/ws/admin/system-health", nil)
			w := httptest.NewRecorder()
			handler.GetAdminSystemHealth(w, req)
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("GET /admin/settings", func() {
		It("should return settings", func() {
			req := httptest.NewRequest("GET", "/v1/ws/admin/settings", nil)
			w := httptest.NewRecorder()
			handler.GetAdminSettings(w, req)
			Expect(w.Code).Should(Equal(200))
		})
	})
})

var _ = Describe("Alarms API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler()
	})

	Context("GET /alarms", func() {
		It("should list alarms", func() {
			req := httptest.NewRequest("GET", "/v1/ws/alarms", nil)
			w := httptest.NewRecorder()
			handler.GetAlarms(w, req, external.GetAlarmsParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})
})

var _ = Describe("Activities API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler()
	})

	Context("GET /activities", func() {
		It("should list activities", func() {
			req := httptest.NewRequest("GET", "/v1/ws/activities", nil)
			w := httptest.NewRecorder()
			handler.GetActivities(w, req, external.GetActivitiesParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})
})

var _ = Describe("Search API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler()
	})

	Context("GET /search", func() {
		It("should return search results", func() {
			req := httptest.NewRequest("GET", "/v1/ws/search", nil)
			w := httptest.NewRecorder()
			handler.GetSearch(w, req, external.GetSearchParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})
})

var _ = Describe("Paste API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler()
	})

	Context("GET /paste", func() {
		It("should list pastes", func() {
			req := httptest.NewRequest("GET", "/v1/ws/paste", nil)
			w := httptest.NewRecorder()
			handler.GetPaste(w, req, external.GetPasteParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})
})

var _ = Describe("HTTP Server", func() {
	Context("Health endpoint", func() {
		It("should return 200", func() {
			resp, err := server.Client().Get(server.URL + "/health")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})
})
