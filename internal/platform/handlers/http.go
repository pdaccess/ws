package handlers

import (
	"context"
	"fmt"
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

func NewHttpHandler(svc ports.Service) external.StrictServerInterface {
	return &httpHandler{svc: svc}
}

func (h *httpHandler) GetActivities(ctx context.Context, request external.GetActivitiesRequestObject) (external.GetActivitiesResponseObject, error) {
	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}

	opts := []domain.ActivitySearchOption{
		domain.WithActivityLimit(limit),
		domain.WithActivityOffset(offset),
	}

	activities, err := h.svc.SearchActivities(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var extActivities []external.Activity
	for _, a := range activities {
		id := int(a.ID.ID())
		extActivities = append(extActivities, external.Activity{
			Id:       &id,
			Message:  &a.Details,
			Severity: func() *external.ActivitySeverity { s := external.ActivitySeverityInfo; return &s }(),
			Source:   &a.Resource,
			Title:    &a.Action,
			Time:     func() *string { s := a.Time.Format("2006-01-02T15:04:05Z"); return &s }(),
		})
	}

	page := 1
	total := len(activities)
	totalPages := (total + limit - 1) / limit

	return external.GetActivities200JSONResponse(external.ActivityList{
		Data: &extActivities,
		Meta: &external.PaginationMeta{
			Limit:      &limit,
			Page:       &page,
			Total:      &total,
			TotalPages: &totalPages,
		},
	}), nil
}

func (h *httpHandler) GetActivitiesId(ctx context.Context, request external.GetActivitiesIdRequestObject) (external.GetActivitiesIdResponseObject, error) {
	activityID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid activity id")
	}

	activities, err := h.svc.GetActivitiesByResourceID(ctx, activityID, 1)
	if err != nil {
		return nil, err
	}

	if len(activities) == 0 {
		return nil, fmt.Errorf("activity not found")
	}

	a := activities[0]
	id := int(a.ID.ID())
	severity := external.ActivityDetailSeverityInfo
	timeStr := a.Time.Format("2006-01-02T15:04:05Z")
	return external.GetActivitiesId200JSONResponse{
		Id:       &id,
		Message:  &a.Details,
		Severity: &severity,
		Source:   &a.Resource,
		Title:    &a.Action,
		Time:     &timeStr,
	}, nil
}

func (h *httpHandler) GetAdminAuditLogs(ctx context.Context, request external.GetAdminAuditLogsRequestObject) (external.GetAdminAuditLogsResponseObject, error) {
	var data []external.AuditLog
	limit := 20
	page := 1
	total := 0
	totalPages := 0
	return external.GetAdminAuditLogs200JSONResponse(external.AuditLogList{
		Data: &data,
		Meta: &external.PaginationMeta{
			Limit:      &limit,
			Page:       &page,
			Total:      &total,
			TotalPages: &totalPages,
		},
	}), nil
}

func (h *httpHandler) GetAdminSettings(ctx context.Context, request external.GetAdminSettingsRequestObject) (external.GetAdminSettingsResponseObject, error) {
	sessionTimeout := 30
	maxLoginAttempts := 5
	passwordMinLength := 8
	passwordRequireMfa := false
	emailNotifications := true
	return external.GetAdminSettings200JSONResponse{
		SessionTimeout:     &sessionTimeout,
		MaxLoginAttempts:   &maxLoginAttempts,
		PasswordMinLength:  &passwordMinLength,
		PasswordRequireMfa: &passwordRequireMfa,
		AllowedIpRanges:    &[]string{},
		EmailNotifications: &emailNotifications,
	}, nil
}

func (h *httpHandler) PutAdminSettings(ctx context.Context, request external.PutAdminSettingsRequestObject) (external.PutAdminSettingsResponseObject, error) {
	return external.PutAdminSettings200JSONResponse{}, nil
}

func (h *httpHandler) GetAdminSystemHealth(ctx context.Context, request external.GetAdminSystemHealthRequestObject) (external.GetAdminSystemHealthResponseObject, error) {
	status := external.Healthy
	uptime := "24h"
	cpu := float32(0.25)
	memory := float32(0.45)
	disk := float32(0.60)
	services := &map[string]any{}
	return external.GetAdminSystemHealth200JSONResponse{
		Status:   &status,
		Uptime:   &uptime,
		Cpu:      &cpu,
		Memory:   &memory,
		Disk:     &disk,
		Services: services,
	}, nil
}

func (h *httpHandler) GetAdminUsers(ctx context.Context, request external.GetAdminUsersRequestObject) (external.GetAdminUsersResponseObject, error) {
	var data []external.User
	limit := 20
	page := 1
	total := 0
	totalPages := 0
	return external.GetAdminUsers200JSONResponse{
		Data: &data,
		Meta: &external.PaginationMeta{
			Limit:      &limit,
			Page:       &page,
			Total:      &total,
			TotalPages: &totalPages,
		},
	}, nil
}

func (h *httpHandler) PostAdminUsers(ctx context.Context, request external.PostAdminUsersRequestObject) (external.PostAdminUsersResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}
	email := request.Body.Email
	id := 0
	return external.PostAdminUsers201JSONResponse(external.User{
		Email: &email,
		Id:    &id,
	}), nil
}

func (h *httpHandler) DeleteAdminUsersId(ctx context.Context, request external.DeleteAdminUsersIdRequestObject) (external.DeleteAdminUsersIdResponseObject, error) {
	return external.DeleteAdminUsersId204Response{}, nil
}

func (h *httpHandler) GetAdminUsersId(ctx context.Context, request external.GetAdminUsersIdRequestObject) (external.GetAdminUsersIdResponseObject, error) {
	userID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid user id")
	}

	user, err := h.svc.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	id := int(user.ID.ID())
	return external.GetAdminUsersId200JSONResponse{
		Email: &user.Email,
		Id:    &id,
	}, nil
}

func (h *httpHandler) PutAdminUsersId(ctx context.Context, request external.PutAdminUsersIdRequestObject) (external.PutAdminUsersIdResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}
	email := ""
	if request.Body.Email != nil {
		email = *request.Body.Email
	}
	id := int(request.Id)
	return external.PutAdminUsersId200JSONResponse(external.User{
		Email: &email,
		Id:    &id,
	}), nil
}

func (h *httpHandler) PutAdminUsersIdStatus(ctx context.Context, request external.PutAdminUsersIdStatusRequestObject) (external.PutAdminUsersIdStatusResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}
	return external.PutAdminUsersIdStatus200Response{}, nil
}

func (h *httpHandler) GetAlarms(ctx context.Context, request external.GetAlarmsRequestObject) (external.GetAlarmsResponseObject, error) {
	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}

	alarms, err := h.svc.SearchAlarms(ctx, uuid.Nil, limit, offset)
	if err != nil {
		return nil, err
	}

	var extAlarms []external.Alarm
	for _, a := range alarms {
		severity := external.AlarmSeverity(a.Severity)
		extAlarms = append(extAlarms, external.Alarm{
			Acknowledged: &a.Acknowledged,
			Message:      &a.Message,
			Severity:     &severity,
			Source:       &a.Source,
			Title:        &a.Title,
			Time:         func() *string { s := a.Time.Format("2006-01-02T15:04:05Z"); return &s }(),
		})
	}

	page := 1
	total := len(alarms)
	totalPages := (total + limit - 1) / limit

	return external.GetAlarms200JSONResponse(external.AlarmList{
		Data: &extAlarms,
		Meta: &external.PaginationMeta{
			Limit:      &limit,
			Page:       &page,
			Total:      &total,
			TotalPages: &totalPages,
		},
	}), nil
}

func (h *httpHandler) GetAlarmsId(ctx context.Context, request external.GetAlarmsIdRequestObject) (external.GetAlarmsIdResponseObject, error) {
	alarmID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid alarm id")
	}

	alarm, err := h.svc.GetAlarm(ctx, alarmID)
	if err != nil {
		return nil, err
	}

	if alarm == nil {
		return nil, fmt.Errorf("alarm not found")
	}

	severity := external.AlarmSeverity(alarm.Severity)
	timeStr := alarm.Time.Format("2006-01-02T15:04:05Z")
	return external.GetAlarmsId200JSONResponse{
		Acknowledged: &alarm.Acknowledged,
		Id:           func() *int { i := int(alarm.ID.ID()); return &i }(),
		Message:      &alarm.Message,
		Severity:     &severity,
		Source:       &alarm.Source,
		Title:        &alarm.Title,
		Time:         &timeStr,
	}, nil
}

func (h *httpHandler) PostAlarmsIdAcknowledge(ctx context.Context, request external.PostAlarmsIdAcknowledgeRequestObject) (external.PostAlarmsIdAcknowledgeResponseObject, error) {
	alarmID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid alarm id")
	}

	if err := h.svc.AcknowledgeAlarm(ctx, alarmID); err != nil {
		return nil, err
	}

	return external.PostAlarmsIdAcknowledge200Response{}, nil
}

func (h *httpHandler) PostGroups(ctx context.Context, request external.PostGroupsRequestObject) (external.PostGroupsResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}

	desc := ""
	if request.Body.Description != nil {
		desc = *request.Body.Description
	}

	group := &domain.Inventory{
		Name:        request.Body.Name,
		Description: desc,
		ItemType:    domain.ItemTypeGroup,
	}

	if err := h.svc.CreateInventory(ctx, group); err != nil {
		return nil, err
	}

	id := int(group.ID.ID())
	return external.PostGroups201JSONResponse(external.Group{
		Name:        &request.Body.Name,
		Description: &desc,
		Id:          &id,
	}), nil
}

func (h *httpHandler) DeleteGroupsId(ctx context.Context, request external.DeleteGroupsIdRequestObject) (external.DeleteGroupsIdResponseObject, error) {
	groupID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid group id")
	}

	if err := h.svc.DeleteInventory(ctx, groupID); err != nil {
		return nil, err
	}

	return external.DeleteGroupsId204Response{}, nil
}

func (h *httpHandler) GetGroupsId(ctx context.Context, request external.GetGroupsIdRequestObject) (external.GetGroupsIdResponseObject, error) {
	groupID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid group id")
	}

	inv, err := h.svc.GetInventory(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if inv == nil {
		return nil, fmt.Errorf("group not found")
	}

	id := int(inv.ID.ID())
	return external.GetGroupsId200JSONResponse{
		Description: &inv.Description,
		Id:          &id,
		Name:        &inv.Name,
	}, nil
}

func (h *httpHandler) GetGroupsIdMembers(ctx context.Context, request external.GetGroupsIdMembersRequestObject) (external.GetGroupsIdMembersResponseObject, error) {
	groupID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid group id")
	}

	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}

	members, err := h.svc.GetInventoryMembers(ctx, groupID, limit, offset)
	if err != nil {
		return nil, err
	}

	page := 1
	total := len(members)
	totalPages := (total + limit - 1) / limit

	var extUsers []external.User
	for _, m := range members {
		email := m.UserID.String()
		id := int(m.UserID.ID())
		extUsers = append(extUsers, external.User{
			Email: &email,
			Id:    &id,
		})
	}

	return external.GetGroupsIdMembers200JSONResponse(external.UserList{
		Data: &extUsers,
		Meta: &external.PaginationMeta{
			Limit:      &limit,
			Page:       &page,
			Total:      &total,
			TotalPages: &totalPages,
		},
	}), nil
}

func (h *httpHandler) PostGroupsIdMembers(ctx context.Context, request external.PostGroupsIdMembersRequestObject) (external.PostGroupsIdMembersResponseObject, error) {
	groupID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid group id")
	}

	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}

	userID := uuid.Nil
	if request.Body.UserId != 0 {
		userID = intToUUID(request.Body.UserId)
	}

	member := &domain.InventoryMember{
		InventoryID: groupID,
		UserID:      userID,
		Role:        "member",
	}

	if err := h.svc.AddInventoryMember(ctx, member); err != nil {
		return nil, err
	}

	return external.PostGroupsIdMembers201Response{}, nil
}

func (h *httpHandler) DeleteGroupsIdMembersUserId(ctx context.Context, request external.DeleteGroupsIdMembersUserIdRequestObject) (external.DeleteGroupsIdMembersUserIdResponseObject, error) {
	return external.DeleteGroupsIdMembersUserId204Response{}, nil
}

func (h *httpHandler) GetGroupsIdPolicies(ctx context.Context, request external.GetGroupsIdPoliciesRequestObject) (external.GetGroupsIdPoliciesResponseObject, error) {
	return external.GetGroupsIdPolicies200JSONResponse{}, nil
}

func (h *httpHandler) PostGroupsIdPolicies(ctx context.Context, request external.PostGroupsIdPoliciesRequestObject) (external.PostGroupsIdPoliciesResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}
	return external.PostGroupsIdPolicies201Response{}, nil
}

func (h *httpHandler) DeleteGroupsIdPoliciesPolicyId(ctx context.Context, request external.DeleteGroupsIdPoliciesPolicyIdRequestObject) (external.DeleteGroupsIdPoliciesPolicyIdResponseObject, error) {
	return external.DeleteGroupsIdPoliciesPolicyId204Response{}, nil
}

func (h *httpHandler) GetPaste(ctx context.Context, request external.GetPasteRequestObject) (external.GetPasteResponseObject, error) {
	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}

	opts := []domain.PasteSearchOption{
		domain.WithPasteLimit(limit),
		domain.WithPasteOffset(offset),
	}

	pastes, err := h.svc.SearchPastes(ctx, opts...)
	if err != nil {
		return nil, err
	}

	var extPastes []external.Paste
	for _, p := range pastes {
		extPastes = append(extPastes, external.Paste{
			Content:   &p.Content,
			CreatedAt: &p.CreatedAt,
			ExpiresAt: p.ExpiresAt,
			Id:        func() *string { s := p.ID.String(); return &s }(),
			Views:     &p.Views,
		})
	}

	page := 1
	total := len(pastes)
	totalPages := (total + limit - 1) / limit

	return external.GetPaste200JSONResponse(external.PasteList{
		Data: &extPastes,
		Meta: &external.PaginationMeta{
			Limit:      &limit,
			Page:       &page,
			Total:      &total,
			TotalPages: &totalPages,
		},
	}), nil
}

func (h *httpHandler) PostPaste(ctx context.Context, request external.PostPasteRequestObject) (external.PostPasteResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}

	paste := domain.Paste{
		UserID:  uuid.Nil,
		Content: request.Body.Content,
	}

	if request.Body.ExpiresIn != nil {
		expiresAt := time.Now().Add(time.Duration(*request.Body.ExpiresIn) * time.Hour)
		paste.ExpiresAt = &expiresAt
	}

	if err := h.svc.CreatePaste(ctx, &paste); err != nil {
		return nil, err
	}

	id := paste.ID.String()
	return external.PostPaste201JSONResponse(external.Paste{
		Content:   &paste.Content,
		CreatedAt: &paste.CreatedAt,
		ExpiresAt: paste.ExpiresAt,
		Id:        &id,
		Views:     &paste.Views,
	}), nil
}

func (h *httpHandler) DeletePasteId(ctx context.Context, request external.DeletePasteIdRequestObject) (external.DeletePasteIdResponseObject, error) {
	pasteID := intToUUID(int(request.Id))

	if err := h.svc.DeletePaste(ctx, pasteID); err != nil {
		return nil, err
	}

	return external.DeletePasteId204Response{}, nil
}

func (h *httpHandler) GetPasteId(ctx context.Context, request external.GetPasteIdRequestObject) (external.GetPasteIdResponseObject, error) {
	pasteID := intToUUID(int(request.Id))

	paste, err := h.svc.GetPaste(ctx, pasteID)
	if err != nil {
		return nil, err
	}

	if paste == nil {
		return nil, fmt.Errorf("paste not found")
	}

	id := paste.ID.String()
	return external.GetPasteId200JSONResponse(external.Paste{
		Content:   &paste.Content,
		CreatedAt: &paste.CreatedAt,
		ExpiresAt: paste.ExpiresAt,
		Id:        &id,
		Views:     &paste.Views,
	}), nil
}

func (h *httpHandler) GetPolicies(ctx context.Context, request external.GetPoliciesRequestObject) (external.GetPoliciesResponseObject, error) {
	var data []external.Policy
	limit := 20
	page := 1
	total := 0
	totalPages := 0
	return external.GetPolicies200JSONResponse{
		Data: &data,
		Meta: &external.PaginationMeta{
			Limit:      &limit,
			Page:       &page,
			Total:      &total,
			TotalPages: &totalPages,
		},
	}, nil
}

func (h *httpHandler) PostPolicies(ctx context.Context, request external.PostPoliciesRequestObject) (external.PostPoliciesResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}
	id := 0
	return external.PostPolicies201JSONResponse(external.Policy{
		Name: &request.Body.Name,
		Id:   &id,
	}), nil
}

func (h *httpHandler) DeletePoliciesId(ctx context.Context, request external.DeletePoliciesIdRequestObject) (external.DeletePoliciesIdResponseObject, error) {
	return external.DeletePoliciesId204Response{}, nil
}

func (h *httpHandler) GetPoliciesId(ctx context.Context, request external.GetPoliciesIdRequestObject) (external.GetPoliciesIdResponseObject, error) {
	id := int(request.Id)
	name := ""
	return external.GetPoliciesId200JSONResponse{
		Id:   &id,
		Name: &name,
	}, nil
}

func (h *httpHandler) PutPoliciesId(ctx context.Context, request external.PutPoliciesIdRequestObject) (external.PutPoliciesIdResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}
	name := ""
	if request.Body.Name != nil {
		name = *request.Body.Name
	}
	id := int(request.Id)
	return external.PutPoliciesId200JSONResponse(external.Policy{
		Name: &name,
		Id:   &id,
	}), nil
}

func (h *httpHandler) GetSearch(ctx context.Context, request external.GetSearchRequestObject) (external.GetSearchResponseObject, error) {
	query := ""
	if request.Params.Q != nil {
		query = *request.Params.Q
	}

	limit := 10
	offset := 0

	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}
	if request.Params.Page != nil {
		offset = (int(*request.Params.Page) - 1) * limit
	}

	var services []external.Service
	var groups []external.Group
	var users []external.User
	var activities []external.Activity

	if query != "" {
		svcItems, err := h.svc.SearchInventory(ctx,
			domain.WithItemType(domain.ItemTypeService),
			domain.WithFilter(query),
			domain.WithLimit(limit),
			domain.WithOffset(offset),
		)
		if err == nil {
			for _, i := range svcItems {
				id := int(i.ID.ID())
				services = append(services, external.Service{
					Id:       &id,
					Name:     &i.Name,
					Hostname: &i.Name,
				})
			}
		}

		groupItems, err := h.svc.SearchInventory(ctx,
			domain.WithItemType(domain.ItemTypeGroup),
			domain.WithFilter(query),
			domain.WithLimit(limit),
			domain.WithOffset(offset),
		)
		if err == nil {
			for _, i := range groupItems {
				id := int(i.ID.ID())
				groups = append(groups, external.Group{
					Id:   &id,
					Name: &i.Name,
				})
			}
		}
	} else if request.Params.Type != "" {
		switch request.Params.Type {
		case external.GetSearchParamsTypeService:
			svcItems, _ := h.svc.SearchInventory(ctx,
				domain.WithItemType(domain.ItemTypeService),
				domain.WithLimit(limit),
				domain.WithOffset(offset),
			)
			for _, i := range svcItems {
				id := int(i.ID.ID())
				services = append(services, external.Service{
					Id:       &id,
					Name:     &i.Name,
					Hostname: &i.Name,
				})
			}
		case external.GetSearchParamsTypeGroup:
			groupItems, _ := h.svc.SearchInventory(ctx,
				domain.WithItemType(domain.ItemTypeGroup),
				domain.WithLimit(limit),
				domain.WithOffset(offset),
			)
			for _, i := range groupItems {
				id := int(i.ID.ID())
				groups = append(groups, external.Group{
					Id:   &id,
					Name: &i.Name,
				})
			}
		}
	}

	return external.GetSearch200JSONResponse(external.SearchResults{
		Activities: &activities,
		Groups:     &groups,
		Services:   &services,
		Users:      &users,
	}), nil
}

func (h *httpHandler) PostServices(ctx context.Context, request external.PostServicesRequestObject) (external.PostServicesResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}

	inv := &domain.Inventory{
		Name:     request.Body.Name,
		ItemType: domain.ItemTypeService,
	}

	if err := h.svc.CreateInventory(ctx, inv); err != nil {
		return nil, err
	}

	id := int(inv.ID.ID())
	return external.PostServices201JSONResponse(external.Service{
		Id:   &id,
		Name: &request.Body.Name,
	}), nil
}

func (h *httpHandler) DeleteServicesId(ctx context.Context, request external.DeleteServicesIdRequestObject) (external.DeleteServicesIdResponseObject, error) {
	inventoryID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid inventory id")
	}

	if err := h.svc.DeleteInventory(ctx, inventoryID); err != nil {
		return nil, err
	}

	return external.DeleteServicesId204Response{}, nil
}

func (h *httpHandler) GetServicesId(ctx context.Context, request external.GetServicesIdRequestObject) (external.GetServicesIdResponseObject, error) {
	svcID, err := uuid.Parse(fmt.Sprint(request.Id))
	if err != nil {
		return nil, fmt.Errorf("invalid service id")
	}

	inv, err := h.svc.GetInventory(ctx, svcID)
	if err != nil {
		return nil, err
	}

	if inv == nil {
		return nil, fmt.Errorf("service not found")
	}

	id := int(inv.ID.ID())
	return external.GetServicesId200JSONResponse{
		Hostname: &inv.Name,
		Id:       &id,
		Name:     &inv.Name,
	}, nil
}

func (h *httpHandler) PutServicesId(ctx context.Context, request external.PutServicesIdRequestObject) (external.PutServicesIdResponseObject, error) {
	if request.Body == nil {
		return nil, fmt.Errorf("missing request body")
	}
	id := int(request.Id)
	return external.PutServicesId200JSONResponse(external.Service{
		Id: &id,
	}), nil
}
