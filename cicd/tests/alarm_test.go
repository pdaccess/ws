package tests

import (
	"net/http/httptest"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alarms API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler(mockSvc)
	})

	Context("GET /alarms", func() {
		It("should list alarms", func() {
			req := httptest.NewRequest("GET", "/v1/ws/alarms", nil)
			w := httptest.NewRecorder()
			handler.GetAlarms(w, req, external.GetAlarmsParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("GET /alarms/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/alarms/1", nil)
			w := httptest.NewRecorder()
			handler.GetAlarmsId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("POST /alarms/{id}/acknowledge", func() {
		It("should acknowledge alarm", func() {
			req := httptest.NewRequest("POST", "/v1/ws/alarms/1/acknowledge", nil)
			w := httptest.NewRecorder()
			handler.PostAlarmsIdAcknowledge(w, req, 1)
			Expect(w.Code).Should(Equal(200))
		})
	})
})
