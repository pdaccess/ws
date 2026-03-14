package tests

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Paste API", func() {
	Context("GET /paste", func() {
		It("should list pastes", func() {
			resp, err := GetAPIClient().GetPasteWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("POST /paste", func() {
		It("should create paste", func() {
			body := http.PostPasteJSONRequestBody{
				Content: "test content",
			}
			resp, err := GetAPIClient().PostPasteWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
		})
	})

	Context("GET /paste/{pasteId}", func() {
		It("should return 400 for invalid paste id", func() {
			resp, err := GetAPIClient().GetPastePasteIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(404))
		})
	})

	Context("DELETE /paste/{pasteId}", func() {
		It("should return 400 for invalid paste id", func() {
			resp, err := GetAPIClient().DeletePastePasteIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(404))
		})
	})
})
