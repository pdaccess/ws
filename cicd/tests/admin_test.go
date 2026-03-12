package tests

import (
	"net/http/httptest"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Admin API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler(mockSvc)
	})

	Context("GET /admin/users", func() {
		It("should list users", func() {
			req := httptest.NewRequest("GET", "/v1/ws/admin/users", nil)
			w := httptest.NewRecorder()
			handler.GetAdminUsers(w, req, external.GetAdminUsersParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("POST /admin/users", func() {
		It("should create user", func() {
			req := httptest.NewRequest("POST", "/v1/ws/admin/users", nil)
			w := httptest.NewRecorder()
			handler.PostAdminUsers(w, req)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("GET /admin/users/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/admin/users/1", nil)
			w := httptest.NewRecorder()
			handler.GetAdminUsersId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("PUT /admin/users/{id}", func() {
		It("should update user", func() {
			req := httptest.NewRequest("PUT", "/v1/ws/admin/users/1", nil)
			w := httptest.NewRecorder()
			handler.PutAdminUsersId(w, req, 1)
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("DELETE /admin/users/{id}", func() {
		It("should delete user", func() {
			req := httptest.NewRequest("DELETE", "/v1/ws/admin/users/1", nil)
			w := httptest.NewRecorder()
			handler.DeleteAdminUsersId(w, req, 1)
			Expect(w.Code).Should(Equal(204))
		})
	})

	Context("PUT /admin/users/{id}/status", func() {
		It("should update user status", func() {
			req := httptest.NewRequest("PUT", "/v1/ws/admin/users/1/status", nil)
			w := httptest.NewRecorder()
			handler.PutAdminUsersIdStatus(w, req, 1)
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("GET /admin/audit-logs", func() {
		It("should list audit logs", func() {
			req := httptest.NewRequest("GET", "/v1/ws/admin/audit-logs", nil)
			w := httptest.NewRecorder()
			handler.GetAdminAuditLogs(w, req, external.GetAdminAuditLogsParams{})
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

	Context("PUT /admin/settings", func() {
		It("should update settings", func() {
			req := httptest.NewRequest("PUT", "/v1/ws/admin/settings", nil)
			w := httptest.NewRecorder()
			handler.PutAdminSettings(w, req)
			Expect(w.Code).Should(Equal(200))
		})
	})
})
