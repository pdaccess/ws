package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Admin API", func() {
	Context("GET /admin/users", func() {
		It("should list users", func() {
			resp, err := GetAPIClient().GetAdminUsersWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("POST /admin/users", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/admin/users", nil)
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("GET /admin/users/{userId}", func() {
		It("should return 400 for invalid user id", func() {
			resp, err := GetAPIClient().GetAdminUsersUserIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("PUT /admin/users/{userId}", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/admin/users/"+uuid.Nil.String(), nil)
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /admin/users/{userId}", func() {
		It("should return 400 for invalid user id", func() {
			resp, err := GetAPIClient().DeleteAdminUsersUserIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("PUT /admin/users/{userId}/status", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/admin/users/"+uuid.Nil.String()+"/status", nil)
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("GET /admin/audit-logs", func() {
		It("should list audit logs", func() {
			resp, err := GetAPIClient().GetAdminAuditLogsWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("GET /admin/system-health", func() {
		It("should return health status", func() {
			resp, err := GetAPIClient().GetAdminSystemHealthWithResponse(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("GET /admin/settings", func() {
		It("should return settings", func() {
			resp, err := GetAPIClient().GetAdminSettingsWithResponse(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("PUT /admin/settings", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/admin/settings", nil)
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})
})
