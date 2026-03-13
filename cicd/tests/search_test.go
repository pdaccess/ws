package tests

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Search API", func() {
	Context("GET /search", func() {
		It("should return search results", func() {
			resp, err := http.Get(GetBaseURL() + "/search?query=test")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(BeNumerically(">=", 200))
		})
	})
})
