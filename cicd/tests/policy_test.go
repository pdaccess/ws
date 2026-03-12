package tests

import (
	"net/http/httptest"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Policies API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler(mockSvc)
	})

	Context("GET /policies", func() {
		It("should list policies", func() {
			req := httptest.NewRequest("GET", "/v1/ws/policies", nil)
			w := httptest.NewRecorder()
			handler.GetPolicies(w, req, external.GetPoliciesParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("POST /policies", func() {
		It("should create policy", func() {
			req := httptest.NewRequest("POST", "/v1/ws/policies", nil)
			w := httptest.NewRecorder()
			handler.PostPolicies(w, req)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("GET /policies/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/policies/1", nil)
			w := httptest.NewRecorder()
			handler.GetPoliciesId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("PUT /policies/{id}", func() {
		It("should update policy", func() {
			req := httptest.NewRequest("PUT", "/v1/ws/policies/1", nil)
			w := httptest.NewRecorder()
			handler.PutPoliciesId(w, req, 1)
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("DELETE /policies/{id}", func() {
		It("should delete policy", func() {
			req := httptest.NewRequest("DELETE", "/v1/ws/policies/1", nil)
			w := httptest.NewRecorder()
			handler.DeletePoliciesId(w, req, 1)
			Expect(w.Code).Should(Equal(204))
		})
	})
})
