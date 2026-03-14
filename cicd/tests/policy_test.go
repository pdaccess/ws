package tests

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Policies API", func() {
	Context("GET /policies", func() {
		It("should list policies", func() {
			resp, err := GetAPIClient().GetPolicies(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("POST /policies", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/policies", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("GET /policies/{id}", func() {
		It("should return 200", func() {
			resp, err := GetAPIClient().GetPoliciesId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("PUT /policies/{id}", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/policies/1", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /policies/{id}", func() {
		It("should delete policy", func() {
			resp, err := GetAPIClient().DeletePoliciesId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(204))
		})
	})
})
