package tests

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Admin API", func() {
	Context("GET /admin/users", func() {
		It("should list users", func() {
			resp, err := GetAPIClient().GetAdminUsers(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("POST /admin/users", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/admin/users", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("GET /admin/users/{id}", func() {
		It("should return 500 (not implemented)", func() {
			resp, err := GetAPIClient().GetAdminUsersId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(500))
		})
	})

	Context("PUT /admin/users/{id}", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/admin/users/1", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /admin/users/{id}", func() {
		It("should delete user", func() {
			resp, err := GetAPIClient().DeleteAdminUsersId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(204))
		})
	})

	Context("PUT /admin/users/{id}/status", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/admin/users/1/status", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("GET /admin/audit-logs", func() {
		It("should list audit logs", func() {
			resp, err := GetAPIClient().GetAdminAuditLogs(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /admin/system-health", func() {
		It("should return health status", func() {
			resp, err := GetAPIClient().GetAdminSystemHealth(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /admin/settings", func() {
		It("should return settings", func() {
			resp, err := GetAPIClient().GetAdminSettings(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("PUT /admin/settings", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/admin/settings", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})
})
