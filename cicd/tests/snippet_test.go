package tests

import (
	"context"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Snippet Endpoints", func() {
	Context("POST /api/v1/ws/snippet", func() {
		It("should handle snippet creation without error", func() {
			content := "test content"
			req := external.NewSnippetRequestObject{
				Body: &external.SnippetRequest{
					Content: &content,
				},
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.NewSnippet(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("GET /api/v1/ws/snippet", func() {
		It("should handle snippet listing without error", func() {
			limit := int32(10)
			offset := int32(0)
			req := external.UserSnippetsRequestObject{
				Params: external.UserSnippetsParams{
					Limit:  &limit,
					Offset: &offset,
				},
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.UserSnippets(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
