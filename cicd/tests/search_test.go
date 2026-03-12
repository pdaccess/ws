package tests

import (
	"net/http/httptest"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Search API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler(mockSvc)
	})

	Context("GET /search", func() {
		It("should return search results", func() {
			req := httptest.NewRequest("GET", "/v1/ws/search", nil)
			w := httptest.NewRecorder()
			handler.GetSearch(w, req, external.GetSearchParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})
})
