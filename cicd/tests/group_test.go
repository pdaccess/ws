package tests

import (
	"net/http/httptest"

	"github.com/pdaccess/ws/internal/platform/handlers"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Groups API", func() {
	var handler external.ServerInterface

	BeforeEach(func() {
		handler = handlers.NewHttpHandler(mockSvc)
	})

	Context("POST /groups", func() {
		It("should create group", func() {
			req := httptest.NewRequest("POST", "/v1/ws/groups", nil)
			w := httptest.NewRecorder()
			handler.PostGroups(w, req)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("GET /groups/{id}", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/groups/1", nil)
			w := httptest.NewRecorder()
			handler.GetGroupsId(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("DELETE /groups/{id}", func() {
		It("should delete group", func() {
			req := httptest.NewRequest("DELETE", "/v1/ws/groups/1", nil)
			w := httptest.NewRecorder()
			handler.DeleteGroupsId(w, req, 1)
			Expect(w.Code).Should(Equal(204))
		})
	})

	Context("GET /groups/{id}/members", func() {
		It("should return members", func() {
			req := httptest.NewRequest("GET", "/v1/ws/groups/1/members", nil)
			w := httptest.NewRecorder()
			handler.GetGroupsIdMembers(w, req, 1, external.GetGroupsIdMembersParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})

	Context("POST /groups/{id}/members", func() {
		It("should add member", func() {
			req := httptest.NewRequest("POST", "/v1/ws/groups/1/members", nil)
			w := httptest.NewRecorder()
			handler.PostGroupsIdMembers(w, req, 1)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("DELETE /groups/{id}/members/{userId}", func() {
		It("should remove member", func() {
			req := httptest.NewRequest("DELETE", "/v1/ws/groups/1/members/1", nil)
			w := httptest.NewRecorder()
			handler.DeleteGroupsIdMembersUserId(w, req, 1, 1)
			Expect(w.Code).Should(Equal(204))
		})
	})

	Context("GET /groups/{id}/policies", func() {
		It("should return 501", func() {
			req := httptest.NewRequest("GET", "/v1/ws/groups/1/policies", nil)
			w := httptest.NewRecorder()
			handler.GetGroupsIdPolicies(w, req, 1)
			Expect(w.Code).Should(Equal(501))
		})
	})

	Context("POST /groups/{id}/policies", func() {
		It("should assign policy", func() {
			req := httptest.NewRequest("POST", "/v1/ws/groups/1/policies", nil)
			w := httptest.NewRecorder()
			handler.PostGroupsIdPolicies(w, req, 1)
			Expect(w.Code).Should(Equal(201))
		})
	})

	Context("DELETE /groups/{id}/policies/{policyId}", func() {
		It("should remove policy", func() {
			req := httptest.NewRequest("DELETE", "/v1/ws/groups/1/policies/1", nil)
			w := httptest.NewRecorder()
			handler.DeleteGroupsIdPoliciesPolicyId(w, req, 1, 1)
			Expect(w.Code).Should(Equal(204))
		})
	})

	Context("GET /groups/{id}/activities", func() {
		It("should return activities", func() {
			req := httptest.NewRequest("GET", "/v1/ws/groups/1/activities", nil)
			w := httptest.NewRecorder()
			handler.GetGroupsIdActivities(w, req, 1, external.GetGroupsIdActivitiesParams{})
			Expect(w.Code).Should(Equal(200))
		})
	})
})
