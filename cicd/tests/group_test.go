package tests

import (
	"context"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Groups API", func() {
	Context("POST /groups", func() {
		It("should create group", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/groups", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})

	Context("GET /groups/{id}", func() {
		It("should return 500 (not implemented)", func() {
			resp, err := GetAPIClient().GetGroupsId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(500))
		})
	})

	Context("DELETE /groups/{id}", func() {
		It("should delete group", func() {
			resp, err := GetAPIClient().DeleteGroupsId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})

	Context("GET /groups/{id}/members", func() {
		It("should return members", func() {
			resp, err := GetAPIClient().GetGroupsIdMembers(context.Background(), 1, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})

	Context("POST /groups/{id}/members", func() {
		It("should add member", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/groups/1/members", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})

	Context("DELETE /groups/{id}/members/{userId}", func() {
		It("should remove member", func() {
			resp, err := GetAPIClient().DeleteGroupsIdMembersUserId(context.Background(), 1, 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(204))
		})
	})

	Context("GET /groups/{id}/policies", func() {
		It("should return 200", func() {
			resp, err := GetAPIClient().GetGroupsIdPolicies(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("POST /groups/{id}/policies", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/groups/1/policies", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /groups/{id}/policies/{policyId}", func() {
		It("should remove policy", func() {
			resp, err := GetAPIClient().DeleteGroupsIdPoliciesPolicyId(context.Background(), 1, 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(204))
		})
	})

	Context("GET /groups/{id}/activities", func() {
		It("should return 404 (not implemented)", func() {
			resp, err := http.Get(GetBaseURL() + "/groups/1/activities")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(404))
		})
	})
})
