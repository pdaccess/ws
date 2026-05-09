package tests

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pdhttp "github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Admin API", func() {
	Context("GET /admin/config", func() {
		It("should return admin config", func() {
			resp, err := GetAPIClient().GetAdminConfigWithResponse(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("PATCH /admin/config", func() {
		It("should update jwt_ttl", func() {
			resp, err := GetAPIClient().PatchAdminConfigWithResponse(context.Background(), pdhttp.PatchAdminConfigJSONRequestBody{
				JwtTtl: new(7200),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})

		It("should update network_whitelist", func() {
			resp, err := GetAPIClient().PatchAdminConfigWithResponse(context.Background(), pdhttp.PatchAdminConfigJSONRequestBody{
				NetworkWhitelist: new([]string{"10.0.0.0/8", "192.168.0.0/16"}),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})

		It("should update rebuild_policy_cache", func() {
			resp, err := GetAPIClient().PatchAdminConfigWithResponse(context.Background(), pdhttp.PatchAdminConfigJSONRequestBody{
				RebuildPolicyCache: new(true),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})

		It("should reject unknown config keys (invalid body)", func() {
			req, _ := http.NewRequest("PATCH", GetBaseURL()+"/admin/config",
				http.NoBody,
			)
			req.Header.Set("Authorization", "Bearer "+GetTestToken())
			req.Header.Set("Content-Type", "application/json")
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("GET /health", func() {
		It("should return 200", func() {
			req, _ := http.NewRequest("GET", GetBaseURL()+"/health", nil)
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})
})
