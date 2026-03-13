package tests

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services API", func() {
	Context("GET /health", func() {
		It("should return 200", func() {
			resp, err := http.Get(GetBaseURL() + "/health")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("POST /services", func() {
		It("should create service", func() {
			req, err := http.NewRequest("POST", GetBaseURL()+"/services", strings.NewReader(`{"name":"test"}`))
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(201))
		})
	})

	Context("GET /services/{id}", func() {
		It("should return 501", func() {
			resp, err := http.Get(GetBaseURL() + "/services/1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(501))
		})
	})

	Context("PUT /services/{id}", func() {
		It("should update service", func() {
			req, err := http.NewRequest("PUT", GetBaseURL()+"/services/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("DELETE /services/{id}", func() {
		It("should delete service", func() {
			req, err := http.NewRequest("DELETE", GetBaseURL()+"/services/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})
})
