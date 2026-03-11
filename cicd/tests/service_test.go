package tests

import (
	"context"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service Endpoints", func() {
	Context("POST /api/v1/ws/service", func() {
		It("should handle valid request without error", func() {
			groupID := uuid.MustParse("4ff2d0ee-5479-433e-84c0-f6f8f44a8f50")
			port := int32(22)
			req := external.NewServiceRequestObject{
				Body: &external.CreateServiceRequest{
					Name:           "test-service",
					GroupId:        groupID,
					AccessProtocol: external.AccessProtocolSsh,
					AuthProtocol:   external.AuthProtocolRadius,
					Vendor:         "Redhat",
					Host:           "demo.pdaccess.com",
					Port:           &port,
				},
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.NewService(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("POST /api/v1/ws/service/index", func() {
		It("should handle search request without error", func() {
			limit := int32(10)
			offset := int32(0)
			req := external.SearchServiceRequestObject{
				Params: external.SearchServiceParams{
					Limit:  &limit,
					Offset: &offset,
				},
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.SearchService(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("GET /api/v1/ws/service/{service}", func() {
		It("should handle service lookup without error", func() {
			req := external.ServiceByIdRequestObject{
				Service: "non-existent-id",
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.ServiceById(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
