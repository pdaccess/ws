package tests

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alarms API", func() {
	Context("GET /alarms", func() {
		It("should list alarms", func() {
			resp, err := GetAPIClient().GetAlarms(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /alarms/{id}", func() {
		It("should return 500 (not implemented)", func() {
			resp, err := GetAPIClient().GetAlarmsId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(500))
		})
	})

	Context("POST /alarms/{id}/acknowledge", func() {
		It("should acknowledge alarm", func() {
			resp, err := GetAPIClient().PostAlarmsIdAcknowledge(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})
})
