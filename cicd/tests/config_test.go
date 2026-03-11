package tests

import (
	"context"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config Endpoints", func() {
	Context("GET /api/v1/ws/config/{realm}/{configContext}", func() {
		It("should handle config fetch without error", func() {
			req := external.FetchContextRequestObject{
				Realm:         "test-realm",
				ConfigContext: "test-context",
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.FetchContext(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("POST /api/v1/ws/config/{realm}/{configContext}", func() {
		It("should handle config upsert without error", func() {
			req := external.UpsertContextRequestObject{
				Realm:         "test-realm",
				ConfigContext: "test-context",
				Body:          &external.ConfigContextRequest{},
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.UpsertContext(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
