package tests

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pdhttp "github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Audit API", func() {
	Context("GET /audit/logs", func() {
		It("should list audit logs", func() {
			resp, err := GetAPIClient().GetAuditLogsWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
		})

		It("should filter audit logs by actor_id", func() {
			resp, err := GetAPIClient().GetAuditLogsWithResponse(context.Background(), &pdhttp.GetAuditLogsParams{
				ActorId: new(uuid.UUID(uuid.New())),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
		})

		It("should filter audit logs by resource_id", func() {
			// create an asset to get a real resource_id
			var spec pdhttp.CreateAssetRequest_Spec
			err := spec.FromServiceSpec(pdhttp.ServiceSpec{
				Endpoint: "https://audit-filter.example.com",
				Protocol: pdhttp.Https,
			})
			Expect(err).ShouldNot(HaveOccurred())

			createResp, err := GetAPIClient().PostAssetsWithResponse(context.Background(), pdhttp.CreateAssetRequest{
				Name: "service-for-audit-filter-" + uuid.NewString()[:8],
				Spec: spec,
				Type: pdhttp.CreateAssetRequestTypeService,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(createResp.StatusCode()).Should(Equal(201))
			Expect(createResp.JSON201).ShouldNot(BeNil())
			Expect(createResp.JSON201.Id).ShouldNot(BeNil())

			id := uuid.UUID(*createResp.JSON201.Id)
			auditResp, err := GetAPIClient().GetAuditLogsWithResponse(context.Background(), &pdhttp.GetAuditLogsParams{
				ResourceId: new(uuid.UUID(id)),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(auditResp.StatusCode()).Should(Equal(200))
			Expect(auditResp.JSON200).ShouldNot(BeNil())
		})

		It("should filter audit logs by actor and resource", func() {
			resp, err := GetAPIClient().GetAuditLogsWithResponse(context.Background(), &pdhttp.GetAuditLogsParams{
				ActorId:    new(uuid.UUID(uuid.New())),
				ResourceId: new(uuid.UUID(uuid.New())),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
		})
	})
})
