package tests

import (
	"context"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Alarm Endpoints", func() {
	Context("GET /api/v1/ws/alarm", func() {
		It("should handle alarm listing without error", func() {
			limit := int32(10)
			offset := int32(0)
			req := external.AlarmIndexRequestObject{
				Params: external.AlarmIndexParams{
					Limit:  &limit,
					Offset: &offset,
				},
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.AlarmIndex(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
