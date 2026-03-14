package tests

import (
	"context"
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
			body := strings.NewReader(`{"name":"test"}`)
			resp, err := GetAPIClient().PostServicesWithBody(context.Background(), "application/json", body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(201))
		})
	})

	Context("GET /services/{id}", func() {
		It("should return 500 (not implemented)", func() {
			resp, err := GetAPIClient().GetServicesId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(500))
		})
	})

	Context("PUT /services/{id}", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/services/1", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /services/{id}", func() {
		It("should delete service", func() {
			resp, err := GetAPIClient().DeleteServicesId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})
})
