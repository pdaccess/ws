package tests

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pdhttp "github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Activities API", func() {
	var createdServiceIDs []uuid.UUID
	var createdGroupIDs []uuid.UUID

	BeforeEach(func() {
		createdServiceIDs = nil
		createdGroupIDs = nil
	})

	AfterEach(func() {
		for _, id := range createdServiceIDs {
			GetAPIClient().DeleteServiceServiceId(context.Background(), id)
		}
		for _, id := range createdGroupIDs {
			GetAPIClient().DeleteGroupGroupId(context.Background(), id)
		}
	})

	Context("GET /activities", func() {
		It("should list activities", func() {
			resp, err := GetAPIClient().GetActivitiesWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
		})

		It("should return activities after creating a service", func() {
			body := pdhttp.PostServiceJSONRequestBody{
				Name:     "test-service-for-activity",
				Hostname: "test-activity.example.com",
				Protocol: "https",
				Type:     "web",
			}
			resp, err := GetAPIClient().PostServiceWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
			id := uuid.MustParse(resp.JSON201.Id.String())
			createdServiceIDs = append(createdServiceIDs, id)

			activitiesResp, err := GetAPIClient().GetActivitiesWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(activitiesResp.StatusCode()).Should(Equal(200))
			Expect(activitiesResp.JSON200).ShouldNot(BeNil())
		})

		It("should return activities after creating a group", func() {
			body := pdhttp.PostGroupJSONRequestBody{
				Name: "test-group-for-activity",
			}
			resp, err := GetAPIClient().PostGroupWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
			id := uuid.MustParse(resp.JSON201.Id.String())
			createdGroupIDs = append(createdGroupIDs, id)

			activitiesResp, err := GetAPIClient().GetActivitiesWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(activitiesResp.StatusCode()).Should(Equal(200))
			Expect(activitiesResp.JSON200).ShouldNot(BeNil())
		})

		It("should filter activities by service id", func() {
			body := pdhttp.PostServiceJSONRequestBody{
				Name:     "service-for-activity-filter",
				Hostname: "filter-activity.example.com",
				Protocol: "https",
				Type:     "api",
			}
			resp, err := GetAPIClient().PostServiceWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
			id := uuid.MustParse(resp.JSON201.Id.String())
			createdServiceIDs = append(createdServiceIDs, id)

			activitiesResp, err := GetAPIClient().GetActivitiesWithResponse(context.Background(), &pdhttp.GetActivitiesParams{
				ServiceId: &id,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(activitiesResp.StatusCode()).Should(Equal(200))
		})

		It("should filter activities by group id", func() {
			body := pdhttp.PostGroupJSONRequestBody{
				Name: "group-for-activity-filter",
			}
			resp, err := GetAPIClient().PostGroupWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
			id := uuid.MustParse(resp.JSON201.Id.String())
			createdGroupIDs = append(createdGroupIDs, id)

			activitiesResp, err := GetAPIClient().GetActivitiesWithResponse(context.Background(), &pdhttp.GetActivitiesParams{
				GroupId: &id,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(activitiesResp.StatusCode()).Should(Equal(200))
		})
	})

	Context("GET /activities/{activityId}", func() {
		It("should return 404 for invalid activity id", func() {
			resp, err := GetAPIClient().GetActivitiesActivityIdWithResponse(context.Background(), uuid.Nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(404))
		})
	})
})
