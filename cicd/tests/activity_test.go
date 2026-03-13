package tests

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Activities API", func() {
	Context("GET /activities", func() {
		It("should list activities", func() {
			resp, err := http.Get(GetBaseURL() + "/activities")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /activities/{id}", func() {
		It("should return 501", func() {
			resp, err := http.Get(GetBaseURL() + "/activities/1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(501))
		})
	})

	Context("GET /sessions/{id}/activities", func() {
		It("should return 404 (not implemented)", func() {
			resp, err := http.Get(GetBaseURL() + "/sessions/1/activities")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(404))
		})
	})
})
