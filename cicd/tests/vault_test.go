package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pdhttp "github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Vaults API", func() {
	var vaultID uuid.UUID

	Context("GET /vaults/{vaultId}/memberships", func() {
		It("should list vault members for existing vault", func() {
			// create a vault asset first
			var spec pdhttp.CreateAssetRequest_Spec
			err := spec.FromVaultSpec(pdhttp.VaultSpec{
				VaultType: "personal",
			})
			Expect(err).ShouldNot(HaveOccurred())

			createResp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "test-vault-members-" + uuid.NewString()[:8],
				Spec: spec,
				Type: pdhttp.CreateAssetRequestTypeVault,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createResp.StatusCode()).Should(Equal(201))
			Expect(createResp.JSON201).ShouldNot(BeNil())
			Expect(createResp.JSON201.Id).ShouldNot(BeNil())
			vaultID = uuid.UUID(*createResp.JSON201.Id)

			resp, err := GetAPIClient().GetVaultsVaultIdMembershipsWithResponse(context.Background(), uuid.UUID(vaultID))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
		})
	})

	Context("POST /vaults/{vaultId}/memberships", func() {
		It("should add a member to a vault", func() {
			if vaultID == uuid.Nil {
				Skip("no vault created")
			}
			resp, err := GetAPIClient().PostVaultsVaultIdMembershipsWithResponse(context.Background(),
				uuid.UUID(vaultID),
				pdhttp.PostVaultsVaultIdMembershipsJSONRequestBody{
					MemberId:   uuid.UUID(uuid.New()),
					MemberType: pdhttp.PostVaultsVaultIdMembershipsJSONBodyMemberTypeUser,
					Role:       pdhttp.Viewer,
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
		})

		It("should add a vault member with owner role", func() {
			if vaultID == uuid.Nil {
				Skip("no vault created")
			}
			resp, err := GetAPIClient().PostVaultsVaultIdMembershipsWithResponse(context.Background(),
				uuid.UUID(vaultID),
				pdhttp.PostVaultsVaultIdMembershipsJSONRequestBody{
					MemberId:   uuid.UUID(uuid.New()),
					MemberType: pdhttp.PostVaultsVaultIdMembershipsJSONBodyMemberTypeUser,
					Role:       pdhttp.Owner,
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
		})

		It("should reject vault membership with missing body", func() {
			if vaultID == uuid.Nil {
				Skip("no vault created")
			}
			req, _ := http.NewRequest("POST",
				GetBaseURL()+"/vaults/"+vaultID.String()+"/memberships",
				nil,
			)
			req.Header.Set("Authorization", "Bearer "+GetTestToken())
			req.Header.Set("Content-Type", "application/json")
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})
})
