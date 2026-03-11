package tests

import (
	"context"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Group Endpoints", func() {
	Context("POST /api/v1/ws/group", func() {
		It("should handle group creation without error", func() {
			req := external.NewGroupRequestObject{
				Body: &external.CreateGroupRequest{
					Name:        "test-group",
					Description: "Test group description",
				},
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.NewGroup(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("POST /api/v1/ws/group/index", func() {
		It("should handle group search without error", func() {
			limit := int32(10)
			offset := int32(0)
			req := external.SearchGroupRequestObject{
				Params: external.SearchGroupParams{
					Limit:  &limit,
					Offset: &offset,
				},
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.SearchGroup(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("GET /api/v1/ws/group/{groupId}", func() {
		It("should handle group lookup without error", func() {
			invalidUUID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
			req := external.GetConfigurationRequestObject{
				GroupId: invalidUUID,
			}
			handler := handlers.NewHttpHandlerWithDefault()
			_, err := handler.GetConfiguration(context.Background(), req)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
