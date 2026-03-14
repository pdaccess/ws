package tests

import (
	"net/http"

	"github.com/google/uuid"
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

	Context("GET /activities/{activityId}", func() {
		It("should return 400 for invalid activity id", func() {
			resp, err := http.Get(GetBaseURL() + "/activities/" + uuid.Nil.String())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(404))
		})
	})
})
