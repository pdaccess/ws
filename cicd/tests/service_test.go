package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pdhttp "github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Services API", func() {
	Context("GET /health", func() {
		It("should return 200", func() {
			resp, err := http.Get(GetBaseURL() + "/health")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("POST /service", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/service", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("GET /service/{serviceId}", func() {
		It("should return 400 for invalid service id", func() {
			resp, err := GetAPIClient().GetServiceServiceIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("PUT /service/{serviceId}", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/service/"+uuid.Nil.String(), nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /service/{serviceId}", func() {
		It("should return 400 for invalid service id", func() {
			resp, err := GetAPIClient().DeleteServiceServiceIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("POST /service", func() {
		It("should create service", func() {
			body := pdhttp.PostServiceJSONRequestBody{
				Name: "test",
			}
			resp, err := GetAPIClient().PostServiceWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
		})
	})
})
