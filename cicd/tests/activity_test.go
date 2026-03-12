package tests

import (
	"net/http/httptest"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Activities API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler(mockSvc)
	})

	Context("GET /activities", func() {
		It("should list activities", func() {
			req := httptest.NewRequest("GET", "/v1/ws/activities", nil)
			w := httptest.NewRecorder()
			handler.GetActivities(w, req, external.GetActivitiesParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("GET /activities/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/activities/1", nil)
			w := httptest.NewRecorder()
			handler.GetActivitiesId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("GET /sessions/{id}/activities", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/sessions/1/activities", nil)
			w := httptest.NewRecorder()
			handler.GetSessionsIdActivities(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})
})
