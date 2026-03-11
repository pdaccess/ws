package handlers

import (
	"net/http"

	"github.com/pdaccess/ws/internal/platform/handlers/external"
)

type httpHandler struct{}

func NewHttpHandler() external.ServerInterface {
	return &httpHandler{}
}

func (h *httpHandler) GetActivities(w http.ResponseWriter, r *http.Request, params external.GetActivitiesParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"data":[],"meta":{"page":1,"limit":20,"total":0,"totalPages":0}}`))
}

func (h *httpHandler) GetActivitiesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) GetAdminAuditLogs(w http.ResponseWriter, r *http.Request, params external.GetAdminAuditLogsParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"data":[],"meta":{"page":1,"limit":20,"total":0,"totalPages":0}}`))
}

func (h *httpHandler) GetAdminSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"sessionTimeout":30,"maxLoginAttempts":5,"passwordMinLength":8,"passwordRequireMfa":false,"allowedIpRanges":[],"emailNotifications":true}`))
}

func (h *httpHandler) PutAdminSettings(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) GetAdminSystemHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy","uptime":"24h","cpu":0.25,"memory":0.45,"disk":0.60,"services":{}}`))
}

func (h *httpHandler) GetAdminUsers(w http.ResponseWriter, r *http.Request, params external.GetAdminUsersParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"data":[],"meta":{"page":1,"limit":20,"total":0,"totalPages":0}}`))
}

func (h *httpHandler) PostAdminUsers(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeleteAdminUsersId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetAdminUsersId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) PutAdminUsersId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) PutAdminUsersIdStatus(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) GetAlarms(w http.ResponseWriter, r *http.Request, params external.GetAlarmsParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"data":[],"meta":{"page":1,"limit":20,"total":0,"totalPages":0}}`))
}

func (h *httpHandler) GetAlarmsId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) PostAlarmsIdAcknowledge(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) PostGroups(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeleteGroupsId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetGroupsId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) GetGroupsIdActivities(w http.ResponseWriter, r *http.Request, id external.IdParam, params external.GetGroupsIdActivitiesParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"data":[],"meta":{"page":1,"limit":20,"total":0,"totalPages":0}}`))
}

func (h *httpHandler) GetGroupsIdMembers(w http.ResponseWriter, r *http.Request, id external.IdParam, params external.GetGroupsIdMembersParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"data":[],"meta":{"page":1,"limit":20,"total":0,"totalPages":0}}`))
}

func (h *httpHandler) PostGroupsIdMembers(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeleteGroupsIdMembersUserId(w http.ResponseWriter, r *http.Request, id external.IdParam, userId int) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetGroupsIdPolicies(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) PostGroupsIdPolicies(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeleteGroupsIdPoliciesPolicyId(w http.ResponseWriter, r *http.Request, id external.IdParam, policyId int) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetPaste(w http.ResponseWriter, r *http.Request, params external.GetPasteParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"data":[],"meta":{"page":1,"limit":20,"total":0,"totalPages":0}}`))
}

func (h *httpHandler) PostPaste(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeletePasteId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetPasteId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) GetPolicies(w http.ResponseWriter, r *http.Request, params external.GetPoliciesParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"data":[],"meta":{"page":1,"limit":20,"total":0,"totalPages":0}}`))
}

func (h *httpHandler) PostPolicies(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeletePoliciesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetPoliciesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) PutPoliciesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) GetSearch(w http.ResponseWriter, r *http.Request, params external.GetSearchParams) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"users":[],"services":[],"groups":[],"activities":[]}`))
}

func (h *httpHandler) PostServices(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeleteServicesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetServicesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) PutServicesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) GetSessionsIdActivities(w http.ResponseWriter, r *http.Request, id int) {
	w.WriteHeader(http.StatusNotImplemented)
}
