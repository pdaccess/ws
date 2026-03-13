package tests

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Policies API", func() {
	Context("GET /policies", func() {
		It("should list policies", func() {
			resp, err := http.Get(GetBaseURL() + "/policies")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("POST /policies", func() {
		It("should create policy", func() {
			resp, err := http.Post(GetBaseURL()+"/policies", "application/json", nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(201))
		})
	})

	Context("GET /policies/{id}", func() {
		It("should return 501", func() {
			resp, err := http.Get(GetBaseURL() + "/policies/1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(501))
		})
	})

	Context("PUT /policies/{id}", func() {
		It("should update policy", func() {
			req, err := http.NewRequest("PUT", GetBaseURL()+"/policies/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("DELETE /policies/{id}", func() {
		It("should delete policy", func() {
			req, err := http.NewRequest("DELETE", GetBaseURL()+"/policies/1", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(204))
		})
	})
})
