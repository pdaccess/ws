package tests

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pdhttp "github.com/pdaccess/ws/pkg/http"
)

var _ = Describe("Identity API", func() {
	var groupID uuid.UUID
	var userID uuid.UUID

	// --- Users ---

	Context("GET /identity/user", func() {
		It("should return current user from JWT (summary)", func() {
			resp, err := GetAPIClient().GetIdentityUserWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
			Expect(resp.JSON200.Id).ShouldNot(BeNil())
		})

		It("should return current user with full profile", func() {
			resp, err := GetAPIClient().GetIdentityUserWithResponse(context.Background(), &pdhttp.GetIdentityUserParams{
				View: ptr(pdhttp.Full),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
		})

		It("should return 404 for non-existent user_id", func() {
			unknownID := uuid.New()
			resp, err := GetAPIClient().GetIdentityUserWithResponse(context.Background(), &pdhttp.GetIdentityUserParams{
				UserId: new(uuid.UUID(unknownID)),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(404))
		})

		It("should return user by specific user_id after creation", func() {
			// create a user first
			username := "test-user-" + uuid.NewString()[:8]
			onboardResp, err := GetAPIClient().PostIdentityUsersWithResponse(context.Background(), pdhttp.PostIdentityUsersJSONRequestBody{
				Username: username,
				Email:    username + "@test.com",
				Password: "testpass123",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(onboardResp.StatusCode()).Should(Equal(201))

			// list users to find the created one
			listResp, err := GetAPIClient().GetIdentityUsersWithResponse(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listResp.JSON200).ShouldNot(BeNil())
			userID = uuid.UUID((*listResp.JSON200)[0].Id)

			// get the user by id
			resp, err := GetAPIClient().GetIdentityUserWithResponse(context.Background(), &pdhttp.GetIdentityUserParams{
				UserId: new(uuid.UUID(userID)),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
		})
	})

	Context("POST /identity/user", func() {
		It("should update current user profile fields", func() {
			displayName := "Updated Name"
			firstName := "Updated"
			lastName := "User"
			email := "updated@test.com"

			resp, err := GetAPIClient().PostIdentityUserWithResponse(context.Background(), pdhttp.UpdateUserProfileRequest{
				DisplayName: &displayName,
				FirstName:   &firstName,
				LastName:    &lastName,
				Email:       &email,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
			Expect(resp.JSON200.DisplayName).ShouldNot(BeNil())
			Expect(*resp.JSON200.DisplayName).Should(Equal(displayName))
			Expect(*resp.JSON200.FirstName).Should(Equal(firstName))
			Expect(*resp.JSON200.LastName).Should(Equal(lastName))
			Expect(resp.JSON200.Email).Should(Equal(email))
		})

		It("should update notification settings", func() {
			notif := map[string]any{
				"email_notifications": true,
				"slack_webhook":       "https://hooks.slack.com/test",
			}
			resp, err := GetAPIClient().PostIdentityUserWithResponse(context.Background(), pdhttp.UpdateUserProfileRequest{
				NotificationSettings: &notif,
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
			Expect(resp.JSON200.NotificationSettings).ShouldNot(BeNil())
		})
	})

	Context("GET /identity/users", func() {
		It("should list users", func() {
			resp, err := GetAPIClient().GetIdentityUsersWithResponse(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
			Expect(len(*resp.JSON200)).Should(BeNumerically(">=", 1))
		})
	})

	Context("POST /identity/users", func() {
		It("should onboard a new user", func() {
			username := "onboard-" + uuid.NewString()[:8]
			resp, err := GetAPIClient().PostIdentityUsersWithResponse(context.Background(), pdhttp.PostIdentityUsersJSONRequestBody{
				Username: username,
				Email:    username + "@pdaccess.io",
				Password: "strongpass123",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
		})

		It("should reject onboarding with missing fields", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/identity/users", nil)
			req.Header.Set("Authorization", "Bearer "+GetTestToken())
			req.Header.Set("Content-Type", "application/json")
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	// --- Groups ---

	Context("GET /identity/groups", func() {
		It("should list groups with default pagination", func() {
			resp, err := GetAPIClient().GetIdentityGroupsWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
			Expect(resp.JSON200.Meta).ShouldNot(BeNil())
		})

		It("should list groups with custom pagination", func() {
			limit := 5
			resp, err := GetAPIClient().GetIdentityGroupsWithResponse(context.Background(), &pdhttp.GetIdentityGroupsParams{
				Limit:  &limit,
				Offset: new(0),
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
		})
	})

	Context("POST /identity/groups", func() {
		It("should create a new group", func() {
			resp, err := GetAPIClient().PostIdentityGroupsWithResponse(context.Background(), pdhttp.PostIdentityGroupsJSONRequestBody{
				Name: "test-group-" + uuid.NewString()[:8],
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
		})

		It("should return 400 for missing name", func() {
			req, _ := http.NewRequest("POST", GetBaseURL()+"/identity/groups", nil)
			req.Header.Set("Authorization", "Bearer "+GetTestToken())
			req.Header.Set("Content-Type", "application/json")
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})

	// --- Group Memberships ---

	Context("GET /identity/groups/{groupId}/memberships", func() {
		It("should list group members", func() {
			listResp, err := GetAPIClient().GetIdentityGroupsWithResponse(context.Background(), nil)
			Expect(err).ShouldNot(HaveOccurred())
			if listResp.JSON200 == nil || len(*listResp.JSON200.Data) == 0 {
				Skip("no groups available")
			}
			gid := uuid.UUID((*listResp.JSON200.Data)[0].Id)
			groupID = gid

			resp, err := GetAPIClient().GetIdentityGroupsGroupIdMembershipsWithResponse(context.Background(), uuid.UUID(gid), nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(200))
			Expect(resp.JSON200).ShouldNot(BeNil())
			Expect(resp.JSON200.Meta).ShouldNot(BeNil())
		})
	})

	Context("POST /identity/groups/{groupId}/memberships", func() {
		It("should assign user to group", func() {
			if groupID == uuid.Nil {
				Skip("no group available")
			}
			// onboard a real user
			username := "member-" + uuid.NewString()[:8]
			onboardResp, err := GetAPIClient().PostIdentityUsersWithResponse(context.Background(), pdhttp.PostIdentityUsersJSONRequestBody{
				Username: username,
				Email:    username + "@test.com",
				Password: "pass123",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(onboardResp.StatusCode()).Should(Equal(201))

			// get the user from the list to get the real ID
			listResp, err := GetAPIClient().GetIdentityUsersWithResponse(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listResp.JSON200).ShouldNot(BeNil())
			Expect(len(*listResp.JSON200)).Should(BeNumerically(">=", 1))
			uid := uuid.UUID((*listResp.JSON200)[0].Id)

			resp, err := GetAPIClient().PostIdentityGroupsGroupIdMembershipsWithResponse(context.Background(),
				uuid.UUID(groupID),
				pdhttp.PostIdentityGroupsGroupIdMembershipsJSONRequestBody{
					UserId: uuid.UUID(uid),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(201))
		})
	})

	Context("DELETE /identity/groups/{groupId}/memberships", func() {
		It("should remove user from group", func() {
			if groupID == uuid.Nil {
				Skip("no group available")
			}
			// onboard a real user for removal
			username := "remove-member-" + uuid.NewString()[:8]
			onboardResp, err := GetAPIClient().PostIdentityUsersWithResponse(context.Background(), pdhttp.PostIdentityUsersJSONRequestBody{
				Username: username,
				Email:    username + "@test.com",
				Password: "pass123",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(onboardResp.StatusCode()).Should(Equal(201))

			listResp, err := GetAPIClient().GetIdentityUsersWithResponse(context.Background())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(listResp.JSON200).ShouldNot(BeNil())
			uid := uuid.UUID((*listResp.JSON200)[0].Id)

			// first add the user
			_, _ = GetAPIClient().PostIdentityGroupsGroupIdMembershipsWithResponse(context.Background(),
				uuid.UUID(groupID),
				pdhttp.PostIdentityGroupsGroupIdMembershipsJSONRequestBody{
					UserId: uuid.UUID(uid),
				},
			)

			resp, err := GetAPIClient().DeleteIdentityGroupsGroupIdMembershipsWithResponse(context.Background(),
				uuid.UUID(groupID),
				&pdhttp.DeleteIdentityGroupsGroupIdMembershipsParams{
					UserId: uuid.UUID(uid),
				},
			)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode()).Should(Equal(204))
		})

		It("should return 400 for missing user_id param", func() {
			if groupID == uuid.Nil {
				Skip("no group available")
			}
			req, _ := http.NewRequest("DELETE",
				GetBaseURL()+"/identity/groups/"+groupID.String()+"/memberships",
				nil,
			)
			req.Header.Set("Authorization", "Bearer "+GetTestToken())
			resp, err := GetHTTPClient().Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(400))
		})
	})
})
