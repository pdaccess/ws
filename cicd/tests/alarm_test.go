package tests

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alarms API", func() {
	Context("GET /alarms", func() {
		It("should list alarms", func() {
			resp, err := GetAPIClient().GetAlarmsWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("GET /alarms/{alarmId}", func() {
		It("should return 400 for invalid alarm id", func() {
			resp, err := GetAPIClient().GetAlarmsAlarmIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})

	Context("POST /alarms/{alarmId}/acknowledge", func() {
		It("should return 400 for invalid alarm id", func() {
			resp, err := GetAPIClient().PostAlarmsAlarmIdAcknowledgeWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(400))
		})
	})
})
