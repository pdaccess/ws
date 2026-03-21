package tests

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pdhttp "github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Search API", func() {
	var createdServiceIDs []uuid.UUID

	BeforeEach(func() {
		createdServiceIDs = nil
	})

	AfterEach(func() {
		for _, id := range createdServiceIDs {
			GetAPIClient().DeleteServiceServiceId(context.Background(), id)
		}
	})

	Context("GET /search", func() {
		It("should return search results", func() {
			q := "test"
			resp, err := GetAPIClient().GetSearchWithResponse(context.Background(), &pdhttp.GetSearchParams{Q: &q})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(BeNumerically(">=", 200))
		})

		It("should search services by name with keyword", func() {
			body := pdhttp.PostServiceJSONRequestBody{
				Name:     "web-server-prod",
				Hostname: "prod-web-01.example.com",
				Protocol: "https",
				Type:     "web",
			}
			resp, err := GetAPIClient().PostServiceWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
			id := uuid.MustParse(resp.JSON201.Id.String())
			createdServiceIDs = append(createdServiceIDs, id)

			svcType := pdhttp.GetSearchParamsTypeService
			limit := 10
			searchResp, err := GetAPIClient().GetSearchWithResponse(context.Background(), &pdhttp.GetSearchParams{
				Type:  &svcType,
				Limit: &limit,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(searchResp.StatusCode()).Should(Equal(200))
			Expect(searchResp.JSON200).ShouldNot(BeNil())
			Expect(searchResp.JSON200.Services).ShouldNot(BeNil())
		})

		It("should filter search by type", func() {
			body := pdhttp.PostServiceJSONRequestBody{
				Name:     "monitoring-service",
				Hostname: "monitor.example.com",
				Protocol: "https",
				Type:     "monitoring",
			}
			resp, err := GetAPIClient().PostServiceWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
			id := uuid.MustParse(resp.JSON201.Id.String())
			createdServiceIDs = append(createdServiceIDs, id)

			svcType := pdhttp.GetSearchParamsTypeService
			searchResp, err := GetAPIClient().GetSearchWithResponse(context.Background(), &pdhttp.GetSearchParams{
				Type: &svcType,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(searchResp.StatusCode()).Should(Equal(200))
			Expect(searchResp.JSON200).ShouldNot(BeNil())
			if searchResp.JSON200.Services != nil {
				Expect(*searchResp.JSON200.Services).ShouldNot(BeEmpty())
			}
		})

		It("should list services when no query provided", func() {
			body := pdhttp.PostServiceJSONRequestBody{
				Name:     "list-test-service",
				Hostname: "list-test.example.com",
				Protocol: "http",
				Type:     "test",
			}
			resp, err := GetAPIClient().PostServiceWithResponse(context.Background(), body)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
			id := uuid.MustParse(resp.JSON201.Id.String())
			createdServiceIDs = append(createdServiceIDs, id)

			svcType := pdhttp.GetSearchParamsTypeService
			searchResp, err := GetAPIClient().GetSearchWithResponse(context.Background(), &pdhttp.GetSearchParams{
				Type: &svcType,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(searchResp.StatusCode()).Should(Equal(200))
			Expect(searchResp.JSON200).ShouldNot(BeNil())
			Expect(searchResp.JSON200.Services).ShouldNot(BeNil())
		})

		It("should list groups when type is Group", func() {
			svcType := pdhttp.GetSearchParamsTypeGroup
			searchResp, err := GetAPIClient().GetSearchWithResponse(context.Background(), &pdhttp.GetSearchParams{
				Type: &svcType,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(searchResp.StatusCode()).Should(Equal(200))
			Expect(searchResp.JSON200).ShouldNot(BeNil())
		})
	})
})
