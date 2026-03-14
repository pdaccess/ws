package tests

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Activities API", func() {
	Context("GET /activities", func() {
		It("should list activities", func() {
			resp, err := GetAPIClient().GetActivities(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("GET /activities/{id}", func() {
		It("should return 500 (not implemented)", func() {
			resp, err := GetAPIClient().GetActivitiesId(context.Background(), 1)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(500))
		})
	})
})
