package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/platform/handlers/external"
)

func intToUUID(id int) uuid.UUID {
	s := fmt.Sprintf("00000000-0000-0000-0000-%012d", id)
	u, _ := uuid.Parse(s)
	return u
}

var dummyUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

func pointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
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
		return nil, fmt.Errorf("failed to search activities %v: %w", err, domain.ErrInternal)
	}

	var extActivities []external.Activity
	for _, a := range activities {
		extActivities = append(extActivities, external.Activity{
			Id:       &a.ID,
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

func (h *httpHandler) GetActivitiesActivityId(ctx context.Context, request external.GetActivitiesActivityIdRequestObject) (external.GetActivitiesActivityIdResponseObject, error) {
	activityID := request.ActivityId
	if activityID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid activity id", Code: domain.ErrCodeInvalidID}
	}

	activities, err := h.svc.GetActivitiesByResourceID(ctx, activityID, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity %v: %w", err, domain.ErrInternal)
	}

	if len(activities) == 0 {
		return nil, domain.NotFoundError{Resource: "activity", ID: request.ActivityId.String(), Code: domain.ErrCodeNotFound}
	}

	a := activities[0]
	severity := external.ActivityDetailSeverityInfo
	timeStr := a.Time.Format("2006-01-02T15:04:05Z")
	return external.GetActivitiesActivityId200JSONResponse{
		Id:       &a.ID,
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

func (h *httpHandler) GetAlarms(ctx context.Context, request external.GetAlarmsRequestObject) (external.GetAlarmsResponseObject, error) {
	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}

	alarms, err := h.svc.SearchAlarms(ctx, uuid.Nil, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search alarms: %w", err)
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

func (h *httpHandler) GetAlarmsAlarmId(ctx context.Context, request external.GetAlarmsAlarmIdRequestObject) (external.GetAlarmsAlarmIdResponseObject, error) {
	alarmID := request.AlarmId
	if alarmID == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid alarm id", Code: domain.ErrCodeInvalidID}
	}

	alarm, err := h.svc.GetAlarm(ctx, alarmID)
	if err != nil {
		return nil, fmt.Errorf("failed to get alarm: %w", err)
	}

	if alarm == nil {
		return nil, domain.NotFoundError{Resource: "alarm", ID: request.AlarmId.String(), Code: domain.ErrCodeNotFound}
	}

	severity := external.AlarmSeverity(alarm.Severity)
	timeStr := alarm.Time.Format("2006-01-02T15:04:05Z")
	return external.GetAlarmsAlarmId200JSONResponse{
		Acknowledged: &alarm.Acknowledged,
		Id:           &alarm.ID,
		Message:      &alarm.Message,
		Severity:     &severity,
		Source:       &alarm.Source,
		Title:        &alarm.Title,
		Time:         &timeStr,
	}, nil
}

func (h *httpHandler) PostAlarmsAlarmIdAcknowledge(ctx context.Context, request external.PostAlarmsAlarmIdAcknowledgeRequestObject) (external.PostAlarmsAlarmIdAcknowledgeResponseObject, error) {
	alarmID := request.AlarmId
	if alarmID == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid alarm id", Code: domain.ErrCodeInvalidID}
	}

	if err := h.svc.AcknowledgeAlarm(ctx, alarmID); err != nil {
		return nil, fmt.Errorf("failed to acknowledge alarm: %w", err)
	}

	return external.PostAlarmsAlarmIdAcknowledge200Response{}, nil
}

func (h *httpHandler) PostGroup(ctx context.Context, request external.PostGroupRequestObject) (external.PostGroupResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	desc := ""
	if request.Body.Description != nil {
		desc = *request.Body.Description
	}

	group := &domain.Group{
		Name:        request.Body.Name,
		Description: desc,
		RealmID:     dummyUUID,
	}

	if err := h.svc.CreateGroup(ctx, group, dummyUUID, dummyUUID); err != nil {
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	return external.PostGroup201JSONResponse(external.Group{
		Name:        &request.Body.Name,
		Description: &desc,
		Id:          &group.ID,
	}), nil
}

func (h *httpHandler) DeleteGroupGroupId(ctx context.Context, request external.DeleteGroupGroupIdRequestObject) (external.DeleteGroupGroupIdResponseObject, error) {
	groupID := request.GroupId
	if groupID == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid group id", Code: domain.ErrCodeInvalidID}
	}

	if err := h.svc.DeleteGroup(ctx, groupID, dummyUUID, dummyUUID); err != nil {
		return nil, fmt.Errorf("failed to delete group: %w", err)
	}

	return external.DeleteGroupGroupId204Response{}, nil
}

func (h *httpHandler) GetGroupGroupId(ctx context.Context, request external.GetGroupGroupIdRequestObject) (external.GetGroupGroupIdResponseObject, error) {
	groupID := request.GroupId
	if groupID == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid group id", Code: domain.ErrCodeInvalidID}
	}

	group, err := h.svc.GetGroup(ctx, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	if group == nil {
		return nil, domain.NotFoundError{Resource: "group", ID: request.GroupId.String(), Code: domain.ErrCodeNotFound}
	}

	return external.GetGroupGroupId200JSONResponse{
		Description: &group.Description,
		Id:          &group.ID,
		Name:        &group.Name,
	}, nil
}

func (h *httpHandler) GetGroupGroupIdMembers(ctx context.Context, request external.GetGroupGroupIdMembersRequestObject) (external.GetGroupGroupIdMembersResponseObject, error) {
	groupID := request.GroupId
	if groupID == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid group id", Code: domain.ErrCodeInvalidID}
	}

	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}

	members, err := h.svc.GetGroupMembers(ctx, groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}

	page := 1
	total := len(members)
	totalPages := (total + limit - 1) / limit

	var data []map[string]any
	for _, m := range members {
		data = append(data, map[string]any{
			"userId": m.UserID.String(),
			"role":   m.Role,
		})
	}

	return external.GetGroupGroupIdMembers200JSONResponse(external.PaginatedResponse{
		Data: &data,
		Meta: &external.PaginationMeta{
			Limit:      &limit,
			Page:       &page,
			Total:      &total,
			TotalPages: &totalPages,
		},
	}), nil
}

func (h *httpHandler) PostGroupGroupIdMembers(ctx context.Context, request external.PostGroupGroupIdMembersRequestObject) (external.PostGroupGroupIdMembersResponseObject, error) {
	groupID := request.GroupId
	if groupID == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid group id", Code: domain.ErrCodeInvalidID}
	}

	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	userID := request.Body.UserId

	role := "member"
	if request.Body.Role != nil {
		role = string(*request.Body.Role)
	}

	member := &domain.GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    role,
	}

	if err := h.svc.AddGroupMember(ctx, member, dummyUUID, dummyUUID); err != nil {
		return nil, fmt.Errorf("failed to add group member: %w", err)
	}

	return external.PostGroupGroupIdMembers201Response{}, nil
}

func (h *httpHandler) DeleteGroupGroupIdMembersUserId(ctx context.Context, request external.DeleteGroupGroupIdMembersUserIdRequestObject) (external.DeleteGroupGroupIdMembersUserIdResponseObject, error) {
	if request.GroupId == uuid.Nil || request.UserId == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid id", Code: domain.ErrCodeInvalidID}
	}

	if err := h.svc.RemoveGroupMembers(ctx, request.GroupId, []uuid.UUID{request.UserId}, dummyUUID, dummyUUID); err != nil {
		return nil, fmt.Errorf("failed to remove group member: %w", err)
	}

	return external.DeleteGroupGroupIdMembersUserId204Response{}, nil
}

func (h *httpHandler) GetGroupGroupIdPolicy(ctx context.Context, request external.GetGroupGroupIdPolicyRequestObject) (external.GetGroupGroupIdPolicyResponseObject, error) {
	if request.GroupId == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid group id", Code: domain.ErrCodeInvalidID}
	}
	return external.GetGroupGroupIdPolicy200JSONResponse{}, nil
}

func (h *httpHandler) PostGroupGroupIdPolicy(ctx context.Context, request external.PostGroupGroupIdPolicyRequestObject) (external.PostGroupGroupIdPolicyResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	return external.PostGroupGroupIdPolicy201Response{}, nil
}

func (h *httpHandler) DeleteGroupGroupIdPolicyPolicyId(ctx context.Context, request external.DeleteGroupGroupIdPolicyPolicyIdRequestObject) (external.DeleteGroupGroupIdPolicyPolicyIdResponseObject, error) {
	if request.GroupId == uuid.Nil || request.PolicyId == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid id", Code: domain.ErrCodeInvalidID}
	}
	return external.DeleteGroupGroupIdPolicyPolicyId204Response{}, nil
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
		return nil, fmt.Errorf("failed to search pastes: %w", err)
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
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
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
		return nil, fmt.Errorf("failed to create paste: %w", err)
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

func (h *httpHandler) DeletePastePasteId(ctx context.Context, request external.DeletePastePasteIdRequestObject) (external.DeletePastePasteIdResponseObject, error) {
	pasteID := request.PasteId

	if err := h.svc.DeletePaste(ctx, pasteID); err != nil {
		return nil, fmt.Errorf("failed to delete paste: %w", err)
	}

	return external.DeletePastePasteId204Response{}, nil
}

func (h *httpHandler) GetPastePasteId(ctx context.Context, request external.GetPastePasteIdRequestObject) (external.GetPastePasteIdResponseObject, error) {
	pasteID := request.PasteId

	paste, err := h.svc.GetPaste(ctx, pasteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get paste: %w", err)
	}

	if paste == nil {
		return nil, domain.NotFoundError{Resource: "paste", ID: request.PasteId.String(), Code: domain.ErrCodeNotFound}
	}

	id := paste.ID.String()
	return external.GetPastePasteId200JSONResponse(external.Paste{
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
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	var id openapi_types.UUID
	return external.PostPolicies201JSONResponse(external.Policy{
		Name: &request.Body.Name,
		Id:   &id,
	}), nil
}

func (h *httpHandler) DeletePoliciesPolicyId(ctx context.Context, request external.DeletePoliciesPolicyIdRequestObject) (external.DeletePoliciesPolicyIdResponseObject, error) {
	if request.PolicyId == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid policy id", Code: domain.ErrCodeInvalidID}
	}
	return external.DeletePoliciesPolicyId204Response{}, nil
}

func (h *httpHandler) GetPoliciesPolicyId(ctx context.Context, request external.GetPoliciesPolicyIdRequestObject) (external.GetPoliciesPolicyIdResponseObject, error) {
	if request.PolicyId == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid policy id", Code: domain.ErrCodeInvalidID}
	}
	id := request.PolicyId
	name := ""
	return external.GetPoliciesPolicyId200JSONResponse{
		Id:   &id,
		Name: &name,
	}, nil
}

func (h *httpHandler) PutPoliciesPolicyId(ctx context.Context, request external.PutPoliciesPolicyIdRequestObject) (external.PutPoliciesPolicyIdResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	name := ""
	if request.Body.Name != nil {
		name = *request.Body.Name
	}
	id := request.PolicyId
	return external.PutPoliciesPolicyId200JSONResponse(external.Policy{
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

	if query != "" {
		svcItems, err := h.svc.SearchServicesWithQuery(ctx, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to search services: %w", err)
		}
		for _, i := range svcItems {
			services = append(services, external.Service{
				Id:       &i.ID,
				Name:     &i.Name,
				Hostname: &i.Name,
			})
		}

		groupItems, err := h.svc.SearchGroupsWithQuery(ctx, query, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to search groups: %w", err)
		}
		for _, i := range groupItems {
			groups = append(groups, external.Group{
				Id:   &i.ID,
				Name: &i.Name,
			})
		}
	} else if request.Params.Type != nil {
		switch *request.Params.Type {
		case external.GetSearchParamsTypeService:
			svcItems, err := h.svc.SearchServices(ctx,
				domain.WithServiceLimit(limit),
				domain.WithServiceOffset(offset),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to search services: %w", err)
			}
			for _, i := range svcItems {
				services = append(services, external.Service{
					Id:       &i.ID,
					Name:     &i.Name,
					Hostname: &i.Name,
				})
			}
		case external.GetSearchParamsTypeGroup:
			groupItems, err := h.svc.SearchGroups(ctx,
				domain.WithGroupLimit(limit),
				domain.WithGroupOffset(offset),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to search groups: %w", err)
			}
			for _, i := range groupItems {
				groups = append(groups, external.Group{
					Id:   &i.ID,
					Name: &i.Name,
				})
			}
		}
	} else {
		svcItems, err := h.svc.SearchServices(ctx,
			domain.WithServiceLimit(limit),
			domain.WithServiceOffset(offset),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to search services: %w", err)
		}
		for _, i := range svcItems {
			services = append(services, external.Service{
				Id:       &i.ID,
				Name:     &i.Name,
				Hostname: &i.Name,
			})
		}

		groupItems, err := h.svc.SearchGroups(ctx,
			domain.WithGroupLimit(limit),
			domain.WithGroupOffset(offset),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to search groups: %w", err)
		}
		for _, i := range groupItems {
			groups = append(groups, external.Group{
				Id:   &i.ID,
				Name: &i.Name,
			})
		}
	}

	return external.GetSearch200JSONResponse(external.SearchResults{
		Groups:   &groups,
		Services: &services,
	}), nil
}

func (h *httpHandler) PostService(ctx context.Context, request external.PostServiceRequestObject) (external.PostServiceResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	svc := &domain.Service{
		Name:    request.Body.Name,
		RealmID: dummyUUID,
	}

	if err := h.svc.CreateService(ctx, svc, dummyUUID, dummyUUID); err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return external.PostService201JSONResponse(external.Service{
		Id:   &svc.ID,
		Name: &request.Body.Name,
	}), nil
}

func (h *httpHandler) DeleteServiceServiceId(ctx context.Context, request external.DeleteServiceServiceIdRequestObject) (external.DeleteServiceServiceIdResponseObject, error) {
	serviceID := request.ServiceId
	if serviceID == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid service id", Code: domain.ErrCodeInvalidID}
	}

	if err := h.svc.DeleteService(ctx, serviceID, dummyUUID, dummyUUID); err != nil {
		return nil, fmt.Errorf("failed to delete service: %w", err)
	}

	return external.DeleteServiceServiceId204Response{}, nil
}

func (h *httpHandler) GetServiceServiceId(ctx context.Context, request external.GetServiceServiceIdRequestObject) (external.GetServiceServiceIdResponseObject, error) {
	serviceID := request.ServiceId
	if serviceID == uuid.Nil {
		return nil, domain.InvalidIDError{Message: "invalid service id", Code: domain.ErrCodeInvalidID}
	}

	svc, err := h.svc.GetService(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	if svc == nil {
		return nil, domain.NotFoundError{Resource: "service", ID: request.ServiceId.String(), Code: domain.ErrCodeNotFound}
	}

	return external.GetServiceServiceId200JSONResponse{
		Hostname: &svc.Name,
		Id:       &svc.ID,
		Name:     &svc.Name,
	}, nil
}

func (h *httpHandler) PutServiceServiceId(ctx context.Context, request external.PutServiceServiceIdRequestObject) (external.PutServiceServiceIdResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	id := request.ServiceId
	return external.PutServiceServiceId200JSONResponse(external.Service{
		Id: &id,
	}), nil
}

func (h *httpHandler) GetGroupGroupIdCredential(ctx context.Context, request external.GetGroupGroupIdCredentialRequestObject) (external.GetGroupGroupIdCredentialResponseObject, error) {
	groupID := request.GroupId
	if groupID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid group id", Code: domain.ErrCodeInvalidID}
	}

	creds, err := h.svc.SearchCredentials(ctx, domain.WithCredentialGroupID(groupID))
	if err != nil {
		return nil, fmt.Errorf("failed to search credentials: %w", err)
	}

	var extCreds []external.Credential
	for _, c := range creds {
		ct := external.CredentialType(c.Type)
		extCreds = append(extCreds, external.Credential{
			Id:       &c.ID,
			GroupId:  &c.GroupID,
			Name:     &c.Name,
			Type:     &ct,
			IsActive: &c.IsActive,
		})
	}

	return external.GetGroupGroupIdCredential200JSONResponse(external.CredentialList{
		Data: &extCreds,
	}), nil
}

func (h *httpHandler) PostGroupGroupIdCredential(ctx context.Context, request external.PostGroupGroupIdCredentialRequestObject) (external.PostGroupGroupIdCredentialResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	groupID := request.GroupId
	if groupID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid group id", Code: domain.ErrCodeInvalidID}
	}

	credType := domain.CredentialType(request.Body.Type)
	cred := &domain.Credential{
		GroupID:  groupID,
		Name:     request.Body.Name,
		Type:     credType,
		IsActive: true,
	}

	if err := h.svc.CreateCredential(ctx, cred); err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}

	return external.PostGroupGroupIdCredential201JSONResponse(external.Credential{
		Id:       &cred.ID,
		GroupId:  &cred.GroupID,
		Name:     &cred.Name,
		Type:     func() *external.CredentialType { t := external.CredentialType(cred.Type); return &t }(),
		IsActive: &cred.IsActive,
	}), nil
}

func (h *httpHandler) GetGroupGroupIdCredentialCredentialId(ctx context.Context, request external.GetGroupGroupIdCredentialCredentialIdRequestObject) (external.GetGroupGroupIdCredentialCredentialIdResponseObject, error) {
	credentialID := request.CredentialId
	if credentialID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid credential id", Code: domain.ErrCodeInvalidID}
	}

	cred, err := h.svc.GetCredential(ctx, credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	if cred == nil {
		return nil, domain.NotFoundError{Resource: "credential", ID: credentialID.String(), Code: domain.ErrCodeNotFound}
	}

	ct := external.CredentialType(cred.Type)
	return external.GetGroupGroupIdCredentialCredentialId200JSONResponse(external.Credential{
		Id:       &cred.ID,
		GroupId:  &cred.GroupID,
		Name:     &cred.Name,
		Type:     &ct,
		IsActive: &cred.IsActive,
	}), nil
}

func (h *httpHandler) PutGroupGroupIdCredentialCredentialId(ctx context.Context, request external.PutGroupGroupIdCredentialCredentialIdRequestObject) (external.PutGroupGroupIdCredentialCredentialIdResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	credentialID := request.CredentialId
	if credentialID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid credential id", Code: domain.ErrCodeInvalidID}
	}

	cred, err := h.svc.GetCredential(ctx, credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}

	if cred == nil {
		return nil, domain.NotFoundError{Resource: "credential", ID: credentialID.String(), Code: domain.ErrCodeNotFound}
	}

	if request.Body.Name != nil {
		cred.Name = *request.Body.Name
	}
	if request.Body.Type != nil {
		cred.Type = domain.CredentialType(*request.Body.Type)
	}
	if request.Body.IsActive != nil {
		cred.IsActive = *request.Body.IsActive
	}

	if err := h.svc.UpdateCredential(ctx, cred); err != nil {
		return nil, fmt.Errorf("failed to update credential: %w", err)
	}

	ct := external.CredentialType(cred.Type)
	return external.PutGroupGroupIdCredentialCredentialId200JSONResponse(external.Credential{
		Id:       &cred.ID,
		GroupId:  &cred.GroupID,
		Name:     &cred.Name,
		Type:     &ct,
		IsActive: &cred.IsActive,
	}), nil
}

func (h *httpHandler) DeleteGroupGroupIdCredentialCredentialId(ctx context.Context, request external.DeleteGroupGroupIdCredentialCredentialIdRequestObject) (external.DeleteGroupGroupIdCredentialCredentialIdResponseObject, error) {
	credentialID := request.CredentialId
	if credentialID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid credential id", Code: domain.ErrCodeInvalidID}
	}

	if err := h.svc.DeleteCredential(ctx, credentialID); err != nil {
		return nil, fmt.Errorf("failed to delete credential: %w", err)
	}

	return external.DeleteGroupGroupIdCredentialCredentialId204Response{}, nil
}

func (h *httpHandler) GetGroupGroupIdCredentialCredentialIdSecret(ctx context.Context, request external.GetGroupGroupIdCredentialCredentialIdSecretRequestObject) (external.GetGroupGroupIdCredentialCredentialIdSecretResponseObject, error) {
	credentialID := request.CredentialId
	if credentialID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid credential id", Code: domain.ErrCodeInvalidID}
	}

	secret, err := h.svc.GetCredentialSecret(ctx, credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential secret: %w", err)
	}

	if secret == nil {
		return nil, domain.NotFoundError{Resource: "credential secret", ID: credentialID.String(), Code: domain.ErrCodeNotFound}
	}

	return external.GetGroupGroupIdCredentialCredentialIdSecret200JSONResponse(external.CredentialSecret{
		Id:           &secret.ID,
		CredentialId: &secret.CredentialID,
		Username:     &secret.Username,
		PrivateKey:   &secret.PrivateKey,
		PublicKey:    &secret.PublicKey,
		ApiKey:       &secret.APIKey,
		Certificate:  &secret.Certificate,
		ExpiresAt:    secret.ExpiresAt,
		LastRotated:  secret.LastRotated,
	}), nil
}

func (h *httpHandler) PostGroupGroupIdCredentialCredentialIdSecret(ctx context.Context, request external.PostGroupGroupIdCredentialCredentialIdSecretRequestObject) (external.PostGroupGroupIdCredentialCredentialIdSecretResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	credentialID := request.CredentialId
	if credentialID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid credential id", Code: domain.ErrCodeInvalidID}
	}

	secret := &domain.CredentialSecret{
		CredentialID:   credentialID,
		Username:       pointerToString(request.Body.Username),
		Password:       pointerToString(request.Body.Password),
		PrivateKey:     pointerToString(request.Body.PrivateKey),
		PublicKey:      pointerToString(request.Body.PublicKey),
		APIKey:         pointerToString(request.Body.ApiKey),
		APIsecret:      pointerToString(request.Body.ApiSecret),
		Certificate:    pointerToString(request.Body.Certificate),
		PrivateKeyPass: pointerToString(request.Body.PrivateKeyPass),
	}

	if request.Body.ExpiresAt != nil {
		secret.ExpiresAt = request.Body.ExpiresAt
	}

	if err := h.svc.CreateCredentialSecret(ctx, secret); err != nil {
		return nil, fmt.Errorf("failed to create credential secret: %w", err)
	}

	return external.PostGroupGroupIdCredentialCredentialIdSecret201JSONResponse(external.CredentialSecret{
		Id:           &secret.ID,
		CredentialId: &secret.CredentialID,
		Username:     &secret.Username,
		ExpiresAt:    secret.ExpiresAt,
	}), nil
}

func (h *httpHandler) PutGroupGroupIdCredentialCredentialIdSecret(ctx context.Context, request external.PutGroupGroupIdCredentialCredentialIdSecretRequestObject) (external.PutGroupGroupIdCredentialCredentialIdSecretResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	credentialID := request.CredentialId
	if credentialID.String() == "" {
		return nil, domain.InvalidIDError{Message: "invalid credential id", Code: domain.ErrCodeInvalidID}
	}

	secret, err := h.svc.GetCredentialSecret(ctx, credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential secret: %w", err)
	}

	if secret == nil {
		return nil, domain.NotFoundError{Resource: "credential secret", ID: credentialID.String(), Code: domain.ErrCodeNotFound}
	}

	if request.Body.Username != nil {
		secret.Username = *request.Body.Username
	}
	if request.Body.Password != nil {
		secret.Password = *request.Body.Password
	}
	if request.Body.PrivateKey != nil {
		secret.PrivateKey = *request.Body.PrivateKey
	}
	if request.Body.PublicKey != nil {
		secret.PublicKey = *request.Body.PublicKey
	}
	if request.Body.ApiKey != nil {
		secret.APIKey = *request.Body.ApiKey
	}
	if request.Body.ApiSecret != nil {
		secret.APIsecret = *request.Body.ApiSecret
	}
	if request.Body.Certificate != nil {
		secret.Certificate = *request.Body.Certificate
	}
	if request.Body.PrivateKeyPass != nil {
		secret.PrivateKeyPass = *request.Body.PrivateKeyPass
	}
	if request.Body.ExpiresAt != nil {
		secret.ExpiresAt = request.Body.ExpiresAt
	}

	if err := h.svc.UpdateCredentialSecret(ctx, secret); err != nil {
		return nil, fmt.Errorf("failed to update credential secret: %w", err)
	}

	return external.PutGroupGroupIdCredentialCredentialIdSecret200JSONResponse(external.CredentialSecret{
		Id:           &secret.ID,
		CredentialId: &secret.CredentialID,
		Username:     &secret.Username,
	}), nil
}
