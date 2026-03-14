package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Policies API", func() {
	Context("GET /policies", func() {
		It("should list policies", func() {
			resp, err := GetAPIClient().GetPoliciesWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
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

	Context("GET /policies/{policyId}", func() {
		It("should return 400 for invalid policy id", func() {
			resp, err := GetAPIClient().GetPoliciesPolicyIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("PUT /policies/{policyId}", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("PUT", GetBaseURL()+"/policies/"+uuid.Nil.String(), nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /policies/{policyId}", func() {
		It("should return 400 for invalid policy id", func() {
			resp, err := GetAPIClient().DeletePoliciesPolicyIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})
})
