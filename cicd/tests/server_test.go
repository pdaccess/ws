package tests

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTP Server Integration", func() {
	Context("Health endpoint", func() {
		It("should return 200 for health check", func() {
			resp, err := http.Get(server.URL + "/health")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})
})
