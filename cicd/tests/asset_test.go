package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pdhttp "github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Assets API", func() {
	var assetID uuid.UUID

	Context("GET /assets", func() {
		It("should list assets with default pagination", func() {
			resp, err := GetAPIClient().GetAssetsWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
			Expect(resp.JSON200.Meta).ShouldNot(BeNil())
		})

		It("should filter assets by type", func() {
			resp, err := GetAPIClient().GetAssetsWithResponse(context.Background(), &pdhttp.GetAssetsParams{
				Type: new("service"),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("GET /assets with search query", func() {
		It("should search assets by query string", func() {
			resp, err := GetAPIClient().GetAssetsWithResponse(context.Background(), &pdhttp.GetAssetsParams{
				Q: new("test"),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("POST /assets", func() {
		It("should create a service asset", func() {
			var spec pdhttp.CreateAssetRequest_Spec
			err := spec.FromServiceSpec(pdhttp.ServiceSpec{
				Endpoint: "https://test.example.com",
				Protocol: pdhttp.Https,
			})
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "test-service-" + uuid.NewString()[:8],
				Spec: spec,
				Type: pdhttp.CreateAssetRequestTypeService,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
			Expect(resp.JSON201).ShouldNot(BeNil())
			Expect(resp.JSON201.Id).ShouldNot(BeNil())
			assetID = uuid.UUID(*resp.JSON201.Id)
		})

		It("should create a vault asset", func() {
			var spec pdhttp.CreateAssetRequest_Spec
			err := spec.FromVaultSpec(pdhttp.VaultSpec{
				VaultType: "personal",
			})
			Expect(err).ShouldNot(HaveOccurred())

			resp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "test-vault-" + uuid.NewString()[:8],
				Spec: spec,
				Type: pdhttp.CreateAssetRequestTypeVault,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
		})

		It("should reject asset creation with empty body", func() {
			// force invalid request via raw HTTP
			req, _ := http.NewRequest("POST", GetBaseURL()+"/assets", nil)
			req.Header.Set("Authorization", "Bearer "+GetTestToken())
			req.Header.Set("Content-Type", "application/json")
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	Context("POST /assets/{id}/actions", func() {
		It("should perform check_out action on existing asset", func() {
			if assetID == uuid.Nil {
				Skip("no asset created in previous test")
			}
			resp, err := GetAPIClient().PostAssetsIdActionsWithResponse(context.Background(),
				uuid.UUID(assetID),
				pdhttp.PostAssetsIdActionsJSONRequestBody{
					ActionType: pdhttp.PostAssetsIdActionsJSONBodyActionTypeCheckOut,
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("GET /health", func() {
		It("should return 200", func() {
			req, _ := http.NewRequest("GET", GetBaseURL()+"/health", nil)
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(200))
		})
	})
})
