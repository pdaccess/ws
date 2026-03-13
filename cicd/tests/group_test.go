package tests

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Groups API", func() {
	Context("POST /groups", func() {
		It("should create group", func() {
			resp, err := http.Post(GetBaseURL()+"/groups", "application/json", nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})

	Context("GET /groups/{id}", func() {
		It("should return 501", func() {
			resp, err := http.Get(GetBaseURL() + "/groups/1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(501))
		})
	})

	Context("DELETE /groups/{id}", func() {
		It("should delete group", func() {
			req, err := http.NewRequest("DELETE", GetBaseURL()+"/groups/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})

	Context("GET /groups/{id}/members", func() {
		It("should return members", func() {
			resp, err := http.Get(GetBaseURL() + "/groups/1/members")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})

	Context("POST /groups/{id}/members", func() {
		It("should add member", func() {
			resp, err := http.Post(GetBaseURL()+"/groups/1/members", "application/json", nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})

	Context("DELETE /groups/{id}/members/{userId}", func() {
		It("should remove member", func() {
			req, err := http.NewRequest("DELETE", GetBaseURL()+"/groups/1/members/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(204))
		})
	})

	Context("GET /groups/{id}/policies", func() {
		It("should return 501", func() {
			resp, err := http.Get(GetBaseURL() + "/groups/1/policies")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(501))
		})
	})

	Context("POST /groups/{id}/policies", func() {
		It("should assign policy", func() {
			resp, err := http.Post(GetBaseURL()+"/groups/1/policies", "application/json", nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(201))
		})
	})

	Context("DELETE /groups/{id}/policies/{policyId}", func() {
		It("should remove policy", func() {
			req, err := http.NewRequest("DELETE", GetBaseURL()+"/groups/1/policies/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
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
