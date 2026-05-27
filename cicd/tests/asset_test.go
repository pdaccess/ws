package tests

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
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

	Context("POST /assets (policy re-evaluation)", func() {
		It("should re-evaluate effective policies when a policy is created after matching assets", func() {
			// Get JWT user ID for policy subject
			currentUser := true
			userResp, err := GetAPIClient().GetIdentityUsersWithResponse(context.Background(), &pdhttp.GetIdentityUsersParams{CurrentUser: &currentUser})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(userResp.JSON200).ShouldNot(BeNil())
			Expect(*userResp.JSON200).ShouldNot(BeEmpty())
			userID := uuid.UUID((*userResp.JSON200)[0].Id)

			// Get or create a group
			groupsResp, err := GetAPIClient().GetIdentityGroupsWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())

			groupID := uuid.New()
			if groupsResp.JSON200 != nil && len(*groupsResp.JSON200.Data) > 0 {
				groupID = uuid.UUID((*groupsResp.JSON200.Data)[0].Id)
			} else {
				groupResp, err := GetAPIClient().PostIdentityGroupsWithResponse(context.Background(), pdhttp.PostIdentityGroupsJSONRequestBody{
					Name: "policy-test-group-" + uuid.NewString()[:8],
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(groupResp.StatusCode()).Should(Equal(201))
				groupsResp2, err := GetAPIClient().GetIdentityGroupsWithResponse(context.Background(), nil)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(groupsResp2.JSON200).ShouldNot(BeNil())
				Expect(len(*groupsResp2.JSON200.Data)).Should(BeNumerically(">=", 1))
				groupID = uuid.UUID((*groupsResp2.JSON200.Data)[0].Id)
			}

			// Add JWT user to group
			_, err = GetAPIClient().PutIdentityUsersWithResponse(context.Background(), pdhttp.UpdateUserProfileRequest{
				GroupIds: &[]openapi_types.UUID{openapi_types.UUID(groupID)},
			})
			Expect(err).ShouldNot(HaveOccurred())

			// Create a service asset with matching tags (before policy creation)
			tags := "db, primary"
			var svcSpec pdhttp.CreateAssetRequest_Spec
			err = svcSpec.FromServiceSpec(pdhttp.ServiceSpec{
				Endpoint: "https://policy-test.example.com",
				Protocol: pdhttp.Https,
				Tags:     &tags,
			})
			Expect(err).ShouldNot(HaveOccurred())

			svcResp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "policy-target-" + uuid.NewString()[:8],
				Spec: svcSpec,
				Type: pdhttp.CreateAssetRequestTypeService,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(svcResp.StatusCode()).Should(Equal(201))

			// Create a policy asset — this triggers reconcilePolicyForAllAssets
			policySpec := map[string]any{
				"actions": []string{"read", "write"},
				"subjects": map[string]any{
					"users":  []string{userID.String()},
					"groups": []string{groupID.String()},
				},
				"objects": map[string]any{
					"tags": map[string]any{
						"tag": "db",
					},
				},
			}
			policyBytes, err := json.Marshal(policySpec)
			Expect(err).ShouldNot(HaveOccurred())

			var policyUnion pdhttp.CreateAssetRequest_Spec
			err = json.Unmarshal(policyBytes, &policyUnion)
			Expect(err).ShouldNot(HaveOccurred())

			policyResp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "test-policy-" + uuid.NewString()[:8],
				Spec: policyUnion,
				Type: pdhttp.CreateAssetRequestTypePolicy,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(policyResp.StatusCode()).Should(Equal(201))
		})

		It("should re-evaluate effective policies when an asset is created after a policy", func() {
			// Get JWT user ID
			currentUser := true
			userResp, err := GetAPIClient().GetIdentityUsersWithResponse(context.Background(), &pdhttp.GetIdentityUsersParams{CurrentUser: &currentUser})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(userResp.JSON200).ShouldNot(BeNil())
			Expect(*userResp.JSON200).ShouldNot(BeEmpty())
			userID := uuid.UUID((*userResp.JSON200)[0].Id)

			// Create a policy first
			policySpec := map[string]any{
				"actions": []string{"read"},
				"subjects": map[string]any{
					"users": []string{userID.String()},
				},
				"objects": map[string]any{
					"tags": map[string]any{
						"tag": "api",
					},
				},
			}
			policyBytes, err := json.Marshal(policySpec)
			Expect(err).ShouldNot(HaveOccurred())
			var policyUnion pdhttp.CreateAssetRequest_Spec
			err = json.Unmarshal(policyBytes, &policyUnion)
			Expect(err).ShouldNot(HaveOccurred())

			policyResp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "pre-policy-" + uuid.NewString()[:8],
				Spec: policyUnion,
				Type: pdhttp.CreateAssetRequestTypePolicy,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(policyResp.StatusCode()).Should(Equal(201))

			// Create a matching service asset — triggers reconcileEffectivePolicies (existing path)
			tags := "api, public"
			var svcSpec pdhttp.CreateAssetRequest_Spec
			err = svcSpec.FromServiceSpec(pdhttp.ServiceSpec{
				Endpoint: "https://api-test.example.com",
				Protocol: pdhttp.Https,
				Tags:     &tags,
			})
			Expect(err).ShouldNot(HaveOccurred())

			svcResp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "api-target-" + uuid.NewString()[:8],
				Spec: svcSpec,
				Type: pdhttp.CreateAssetRequestTypeService,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(svcResp.StatusCode()).Should(Equal(201))
		})

		It("should succeed creating a policy with no matching assets", func() {
			// Create a policy with no matching objects
			policySpec := map[string]any{
				"actions": []string{"admin"},
				"subjects": map[string]any{
					"users": []string{uuid.New().String()},
				},
				"objects": map[string]any{
					"tags": map[string]any{
						"tag": "nonexistent",
					},
				},
			}
			policyBytes, err := json.Marshal(policySpec)
			Expect(err).ShouldNot(HaveOccurred())
			var policyUnion pdhttp.CreateAssetRequest_Spec
			err = json.Unmarshal(policyBytes, &policyUnion)
			Expect(err).ShouldNot(HaveOccurred())

			policyResp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "noop-policy-" + uuid.NewString()[:8],
				Spec: policyUnion,
				Type: pdhttp.CreateAssetRequestTypePolicy,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(policyResp.StatusCode()).Should(Equal(201))
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
