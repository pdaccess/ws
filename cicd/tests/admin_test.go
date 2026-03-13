package tests

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Admin API", func() {
	Context("GET /admin/users", func() {
		It("should list users", func() {
			resp, err := http.Get(GetBaseURL() + "/admin/users")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("POST /admin/users", func() {
		It("should create user", func() {
			resp, err := http.Post(GetBaseURL()+"/admin/users", "application/json", nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(201))
		})
	})

	Context("GET /admin/users/{id}", func() {
		It("should return 501", func() {
			resp, err := http.Get(GetBaseURL() + "/admin/users/1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(501))
		})
	})

	Context("PUT /admin/users/{id}", func() {
		It("should update user", func() {
			req, err := http.NewRequest("PUT", GetBaseURL()+"/admin/users/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("DELETE /admin/users/{id}", func() {
		It("should delete user", func() {
			req, err := http.NewRequest("DELETE", GetBaseURL()+"/admin/users/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(204))
		})
	})

	Context("PUT /admin/users/{id}/status", func() {
		It("should update user status", func() {
			req, err := http.NewRequest("PUT", GetBaseURL()+"/admin/users/1/status", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /admin/audit-logs", func() {
		It("should list audit logs", func() {
			resp, err := http.Get(GetBaseURL() + "/admin/audit-logs")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /admin/system-health", func() {
		It("should return health status", func() {
			resp, err := http.Get(GetBaseURL() + "/admin/system-health")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /admin/settings", func() {
		It("should return settings", func() {
			resp, err := http.Get(GetBaseURL() + "/admin/settings")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("PUT /admin/settings", func() {
		It("should update settings", func() {
			req, err := http.NewRequest("PUT", GetBaseURL()+"/admin/settings", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})
})
