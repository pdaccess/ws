package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Groups API", func() {
	Context("POST /group", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/group", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("GET /group/{groupId}", func() {
		It("should return 400 for invalid group id", func() {
			resp, err := GetAPIClient().GetGroupGroupIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("DELETE /group/{groupId}", func() {
		It("should return 400 for invalid group id", func() {
			resp, err := GetAPIClient().DeleteGroupGroupIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("GET /group/{groupId}/members", func() {
		It("should return 400 for invalid group id", func() {
			resp, err := GetAPIClient().GetGroupGroupIdMembersWithResponse(context.Background(), uuid.Nil, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("POST /group/{groupId}/members", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/group/"+uuid.Nil.String()+"/members", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /group/{groupId}/members/{userId}", func() {
		It("should return 400 for invalid ids", func() {
			resp, err := GetAPIClient().DeleteGroupGroupIdMembersUserIdWithResponse(context.Background(), uuid.Nil, uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("GET /group/{groupId}/policy", func() {
		It("should return 400 for invalid group id", func() {
			resp, err := GetAPIClient().GetGroupGroupIdPolicyWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("POST /group/{groupId}/policy", func() {
		It("should return 400 for missing body", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/group/"+uuid.Nil.String()+"/policy", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("DELETE /group/{groupId}/policy/{policyId}", func() {
		It("should return 400 for invalid ids", func() {
			resp, err := GetAPIClient().DeleteGroupGroupIdPolicyPolicyIdWithResponse(context.Background(), uuid.Nil, uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("GET /group/{groupId}/activities", func() {
		It("should return 404", func() {
			resp, err := http.Get(GetBaseURL() + "/group/" + uuid.Nil.String() + "/activities")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(404))
		})
	})
})
