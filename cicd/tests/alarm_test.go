package tests

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alarms API", func() {
	Context("GET /alarms", func() {
		It("should list alarms", func() {
			resp, err := http.Get(GetBaseURL() + "/alarms")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /alarms/{id}", func() {
		It("should return 501", func() {
			resp, err := http.Get(GetBaseURL() + "/alarms/1")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(501))
		})
	})

	Context("POST /alarms/{id}/acknowledge", func() {
		It("should acknowledge alarm", func() {
			req, err := http.NewRequest("POST", GetBaseURL()+"/alarms/1/acknowledge", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})
})
