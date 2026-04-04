package tests

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Admin API", func() {
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
