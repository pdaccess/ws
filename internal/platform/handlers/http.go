package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/platform/handlers/external"
)

func intToUUID(id int) uuid.UUID {
	s := fmt.Sprintf("00000000-0000-0000-0000-%012d", id)
	u, _ := uuid.Parse(s)
	return u
}

type httpHandler struct {
	svc ports.Service
}

func NewHttpHandler(svc ports.Service) external.ServerInterface {
	return &httpHandler{svc: svc}
}

func (h *httpHandler) GetActivities(w http.ResponseWriter, r *http.Request, params external.GetActivitiesParams) {
	ctx := r.Context()

	limit := 20
	offset := 0
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	opts := []domain.ActivitySearchOption{
		domain.WithActivityLimit(limit),
		domain.WithActivityOffset(offset),
	}

	activities, err := h.svc.SearchActivities(ctx, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"data": activities,
		"meta": map[string]any{
			"page":       1,
			"limit":      limit,
			"total":      len(activities),
			"totalPages": (len(activities) + limit - 1) / limit,
		},
	})
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
	ctx := r.Context()

	limit := 20
	offset := 0
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	alarms, err := h.svc.SearchAlarms(ctx, uuid.Nil, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"data": alarms,
		"meta": map[string]any{
			"page":       1,
			"limit":      limit,
			"total":      len(alarms),
			"totalPages": (len(alarms) + limit - 1) / limit,
		},
	})
}

func (h *httpHandler) GetAlarmsId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) PostAlarmsIdAcknowledge(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	ctx := r.Context()

	alarmID, err := uuid.Parse(fmt.Sprint(id))
	if err != nil {
		http.Error(w, "invalid alarm id", http.StatusBadRequest)
		return
	}

	if err := h.svc.AcknowledgeAlarm(ctx, alarmID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *httpHandler) PostGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group := &domain.Inventory{
		Name:        req.Name,
		Description: req.Description,
		ItemType:    domain.ItemTypeGroup,
	}

	if err := h.svc.CreateInventory(ctx, group); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeleteGroupsId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	ctx := r.Context()

	groupID, err := uuid.Parse(fmt.Sprint(id))
	if err != nil {
		http.Error(w, "invalid group id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteInventory(ctx, groupID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetGroupsId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) GetGroupsIdMembers(w http.ResponseWriter, r *http.Request, id external.IdParam, params external.GetGroupsIdMembersParams) {
	ctx := r.Context()

	groupID, err := uuid.Parse(fmt.Sprint(id))
	if err != nil {
		http.Error(w, "invalid group id", http.StatusBadRequest)
		return
	}

	limit := 20
	offset := 0
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	members, err := h.svc.GetInventoryMembers(ctx, groupID, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"data": members,
		"meta": map[string]any{
			"page":       1,
			"limit":      limit,
			"total":      len(members),
			"totalPages": (len(members) + limit - 1) / limit,
		},
	})
}

func (h *httpHandler) PostGroupsIdMembers(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	ctx := r.Context()

	groupID, err := uuid.Parse(fmt.Sprint(id))
	if err != nil {
		http.Error(w, "invalid group id", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID string `json:"userId"`
		Role   string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	member := &domain.InventoryMember{
		InventoryID: groupID,
		UserID:      userID,
		Role:        req.Role,
	}

	if err := h.svc.AddInventoryMember(ctx, member); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	ctx := r.Context()

	limit := 20
	offset := 0
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	opts := []domain.PasteSearchOption{
		domain.WithPasteLimit(limit),
		domain.WithPasteOffset(offset),
	}

	pastes, err := h.svc.SearchPastes(ctx, opts...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"data": pastes,
		"meta": map[string]any{
			"page":       1,
			"limit":      limit,
			"total":      len(pastes),
			"totalPages": (len(pastes) + limit - 1) / limit,
		},
	})
}

func (h *httpHandler) PostPaste(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req external.PostPasteJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paste := domain.Paste{
		UserID:  uuid.Nil,
		Content: req.Content,
	}

	if req.ExpiresIn != nil {
		expiresAt := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Hour)
		paste.ExpiresAt = &expiresAt
	}

	if err := h.svc.CreatePaste(ctx, &paste); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(paste)
}

func (h *httpHandler) DeletePasteId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	ctx := r.Context()

	pasteID := intToUUID(int(id))

	if err := h.svc.DeletePaste(ctx, pasteID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetPasteId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	ctx := r.Context()

	pasteID := intToUUID(int(id))

	paste, err := h.svc.GetPaste(ctx, pasteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if paste == nil {
		http.Error(w, "paste not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paste)
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
	ctx := r.Context()

	query := ""
	if params.Q != nil {
		query = *params.Q
	}

	limit := 10
	offset := 0

	if params.Limit != nil {
		limit = int(*params.Limit)
	}
	if params.Page != nil {
		offset = (int(*params.Page) - 1) * limit
	}

	results := map[string]any{
		"users":      []any{},
		"services":   []any{},
		"groups":     []any{},
		"activities": []any{},
	}

	if query != "" {
		svcItems, err := h.svc.SearchInventory(ctx,
			domain.WithItemType(domain.ItemTypeService),
			domain.WithFilter(query),
			domain.WithLimit(limit),
			domain.WithOffset(offset),
		)
		if err == nil {
			results["services"] = svcItems
		}

		groupItems, err := h.svc.SearchInventory(ctx,
			domain.WithItemType(domain.ItemTypeGroup),
			domain.WithFilter(query),
			domain.WithLimit(limit),
			domain.WithOffset(offset),
		)
		if err == nil {
			results["groups"] = groupItems
		}
	} else if params.Type != "" {
		switch params.Type {
		case "Service":
			svcItems, _ := h.svc.SearchInventory(ctx,
				domain.WithItemType(domain.ItemTypeService),
				domain.WithLimit(limit),
				domain.WithOffset(offset),
			)
			results["services"] = svcItems
		case "Group":
			groupItems, _ := h.svc.SearchInventory(ctx,
				domain.WithItemType(domain.ItemTypeGroup),
				domain.WithLimit(limit),
				domain.WithOffset(offset),
			)
			results["groups"] = groupItems
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (h *httpHandler) PostServices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		RealmID     string `json:"realmId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	realmID, err := uuid.Parse(req.RealmID)
	if err != nil {
		realmID = uuid.Nil
	}

	inv := &domain.Inventory{
		Name:        req.Name,
		Description: req.Description,
		RealmID:     realmID,
		ItemType:    domain.ItemTypeService,
	}

	if err := h.svc.CreateInventory(ctx, inv); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *httpHandler) DeleteServicesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	ctx := r.Context()

	inventoryID, err := uuid.Parse(fmt.Sprint(id))
	if err != nil {
		http.Error(w, "invalid inventory id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteInventory(ctx, inventoryID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *httpHandler) GetServicesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *httpHandler) PutServicesId(w http.ResponseWriter, r *http.Request, id external.IdParam) {
	w.WriteHeader(http.StatusOK)
}
