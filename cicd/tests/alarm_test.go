package tests

import (
	"net/http"

	"github.com/google/uuid"
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

	Context("GET /alarms/{alarmId}", func() {
		It("should return 400 for invalid alarm id", func() {
			resp, err := http.Get(GetBaseURL() + "/alarms/" + uuid.Nil.String())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("POST /alarms/{alarmId}/acknowledge", func() {
		It("should return 400 for invalid alarm id", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/alarms/"+uuid.Nil.String()+"/acknowledge", nil)
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})
})
