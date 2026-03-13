package tests

import (
	"bytes"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Paste API", func() {
	Context("GET /paste", func() {
		It("should list pastes", func() {
			resp, err := http.Get(GetBaseURL() + "/paste")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})

	Context("POST /paste", func() {
		It("should create paste", func() {
			body := []byte(`{"content": "test content"}`)
			resp, err := http.Post(GetBaseURL()+"/paste", "application/json", bytes.NewReader(body))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(201))
		})
	})

	Context("GET /paste/{id}", func() {
		It("should return 404 for non-existent paste", func() {
			resp, err := http.Get(GetBaseURL() + "/paste/999")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(404))
		})
	})

	Context("DELETE /paste/{id}", func() {
		It("should return 204 for delete (idempotent)", func() {
			req, err := http.NewRequest("DELETE", GetBaseURL()+"/paste/999", nil)
			Expect(err).ShouldNot(HaveOccurred())
			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(204))
		})
	})
})
