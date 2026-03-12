package tests

import (
	"net/http/httptest"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler(mockSvc)
	})

	Context("POST /services", func() {
		It("should create service", func() {
			req := httptest.NewRequest("POST", "/v1/ws/services", nil)
			w := httptest.NewRecorder()
			handler.PostServices(w, req)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("GET /services/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/services/1", nil)
			w := httptest.NewRecorder()
			handler.GetServicesId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("PUT /services/{id}", func() {
		It("should update service", func() {
			req := httptest.NewRequest("PUT", "/v1/ws/services/1", nil)
			w := httptest.NewRecorder()
			handler.PutServicesId(w, req, 1)
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("DELETE /services/{id}", func() {
		It("should delete service", func() {
			req := httptest.NewRequest("DELETE", "/v1/ws/services/1", nil)
			w := httptest.NewRecorder()
			handler.DeleteServicesId(w, req, 1)
			Expect(w.Code).Should(Equal(204))
		})
	})
})
