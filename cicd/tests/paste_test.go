package tests

import (
	"net/http/httptest"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Paste API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler(mockSvc)
	})

	Context("GET /paste", func() {
		It("should list pastes", func() {
			req := httptest.NewRequest("GET", "/v1/ws/paste", nil)
			w := httptest.NewRecorder()
			handler.GetPaste(w, req, external.GetPasteParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("POST /paste", func() {
		It("should create paste", func() {
			req := httptest.NewRequest("POST", "/v1/ws/paste", nil)
			w := httptest.NewRecorder()
			handler.PostPaste(w, req)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("GET /paste/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/paste/1", nil)
			w := httptest.NewRecorder()
			handler.GetPasteId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("DELETE /paste/{id}", func() {
		It("should delete paste", func() {
			req := httptest.NewRequest("DELETE", "/v1/ws/paste/1", nil)
			w := httptest.NewRecorder()
			handler.DeletePasteId(w, req, 1)
			Expect(w.Code).Should(Equal(204))
		})
	})
})
