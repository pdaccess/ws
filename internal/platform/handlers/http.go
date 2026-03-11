package handlers

import (
	"context"
	"time"

	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/platform/handlers/external"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type httpHandler struct {
	service ports.Service
}

func NewHttpHandler(service ports.Service) external.StrictServerInterface {
	return &httpHandler{service: service}
}

func NewHttpHandlerWithDefault() external.StrictServerInterface {
	return &httpHandler{}
}

func (h *httpHandler) AlarmIndex(ctx context.Context, request external.AlarmIndexRequestObject) (external.AlarmIndexResponseObject, error) {
	return external.AlarmIndex200JSONResponse{}, nil
}

func (h *httpHandler) CreateGroupAlarm(ctx context.Context, request external.CreateGroupAlarmRequestObject) (external.CreateGroupAlarmResponseObject, error) {
	return external.CreateGroupAlarm200JSONResponse{}, nil
}

func (h *httpHandler) CreateGroupMember(ctx context.Context, request external.CreateGroupMemberRequestObject) (external.CreateGroupMemberResponseObject, error) {
	return external.CreateGroupMember200JSONResponse{}, nil
}

func (h *httpHandler) CreateGroupMessage(ctx context.Context, request external.CreateGroupMessageRequestObject) (external.CreateGroupMessageResponseObject, error) {
	return external.CreateGroupMessage200JSONResponse{}, nil
}

func (h *httpHandler) DeleteGroup(ctx context.Context, request external.DeleteGroupRequestObject) (external.DeleteGroupResponseObject, error) {
	return nil, nil
}

func (h *httpHandler) DeleteGroupAlarm(ctx context.Context, request external.DeleteGroupAlarmRequestObject) (external.DeleteGroupAlarmResponseObject, error) {
	return external.DeleteGroupAlarm200JSONResponse{}, nil
}

func (h *httpHandler) DeleteGroupMembers(ctx context.Context, request external.DeleteGroupMembersRequestObject) (external.DeleteGroupMembersResponseObject, error) {
	return external.DeleteGroupMembers200JSONResponse{}, nil
}

func (h *httpHandler) DeleteGroupMessage(ctx context.Context, request external.DeleteGroupMessageRequestObject) (external.DeleteGroupMessageResponseObject, error) {
	return external.DeleteGroupMessage200JSONResponse{}, nil
}

func (h *httpHandler) DeleteService(ctx context.Context, request external.DeleteServiceRequestObject) (external.DeleteServiceResponseObject, error) {
	return nil, nil
}

func (h *httpHandler) DeleteSnippet(ctx context.Context, request external.DeleteSnippetRequestObject) (external.DeleteSnippetResponseObject, error) {
	return external.DeleteSnippet200JSONResponse{}, nil
}

func (h *httpHandler) FetchContext(ctx context.Context, request external.FetchContextRequestObject) (external.FetchContextResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized, returning empty response")
		return external.FetchContext200JSONResponse{}, nil
	}

	items, err := h.service.GetItemContext(ctx, domain.ItemContext(request.ConfigContext))
	if err != nil {
		log.Error().Err(err).Str("context", request.ConfigContext).Msg("failed to fetch config context")
		return nil, err
	}

	response := make(external.ConfigContextResponse, len(items))
	for i, item := range items {
		response[i] = external.ConfigContext{
			Entry: &item.Key,
			Value: &item.Value,
		}
	}

	return external.FetchContext200JSONResponse(response), nil
}

func (h *httpHandler) FetchGroupAlarms(ctx context.Context, request external.FetchGroupAlarmsRequestObject) (external.FetchGroupAlarmsResponseObject, error) {
	return external.FetchGroupAlarms200JSONResponse{}, nil
}

func (h *httpHandler) GetConfiguration(ctx context.Context, request external.GetConfigurationRequestObject) (external.GetConfigurationResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized, returning empty response")
		return external.GetConfiguration200JSONResponse{}, nil
	}

	groupID := uuid.UUID(request.GroupId)
	group, err := h.service.GetInventory(ctx, groupID)
	if err != nil {
		log.Error().Err(err).Str("groupId", groupID.String()).Msg("failed to get group")
		return nil, err
	}

	return external.GetConfiguration200JSONResponse{
		GroupId:     group.ID,
		Name:        group.Name,
		Description: group.Description,
		Parent:      group.ParentID,
	}, nil
}

func (h *httpHandler) GetSnippet(ctx context.Context, request external.GetSnippetRequestObject) (external.GetSnippetResponseObject, error) {
	return external.GetSnippet200JSONResponse{}, nil
}

func (h *httpHandler) GroupMembers(ctx context.Context, request external.GroupMembersRequestObject) (external.GroupMembersResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized, returning empty response")
		return external.GroupMembers200JSONResponse{}, nil
	}

	groupID, err := uuid.Parse(request.Group)
	if err != nil {
		return nil, err
	}

	limit := 10
	offset := 0
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}
	if request.Params.Offset != nil {
		offset = int(*request.Params.Offset)
	}

	members, err := h.service.GetInventoryMembers(ctx, groupID, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("groupId", groupID.String()).Msg("failed to get group members")
		return nil, err
	}

	response := make(external.GroupMemebersResponse, len(members))
	for i, m := range members {
		userID := m.UserID
		role := m.Role
		membershipTime := m.MembershipTime.Format(time.RFC3339)
		response[i] = external.GroupMemeberItem{
			UserId:         &userID,
			Role:           &role,
			MembershipTime: &membershipTime,
		}
	}

	return external.GroupMembers200JSONResponse(response), nil
}

func (h *httpHandler) NewGroup(ctx context.Context, request external.NewGroupRequestObject) (external.NewGroupResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized")
		return nil, nil
	}

	if request.Body == nil || request.Body.Name == "" {
		return nil, nil
	}

	group := &domain.Inventory{
		ID:          uuid.New(),
		Name:        request.Body.Name,
		Description: request.Body.Description,
		ItemType:    domain.ItemTypeGroup,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if request.Body.Parent != nil {
		parentID, err := uuid.Parse(*request.Body.Parent)
		if err == nil {
			group.ParentID = &parentID
		}
	}

	err := h.service.CreateInventory(ctx, group)
	if err != nil {
		log.Error().Err(err).Msg("failed to create group")
		return nil, err
	}

	return external.NewGroup200JSONResponse{
		GroupId: &group.ID,
	}, nil
}

func (h *httpHandler) NewService(ctx context.Context, request external.NewServiceRequestObject) (external.NewServiceResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized")
		return nil, nil
	}

	if request.Body == nil || request.Body.Name == "" {
		return nil, nil
	}

	groupID := uuid.UUID(request.Body.GroupId)
	inventory := &domain.Inventory{
		ID:          uuid.New(),
		ParentID:    &groupID,
		Name:        request.Body.Name,
		Description: getString(request.Body.Description),
		ItemType:    domain.ItemTypeService,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := h.service.CreateInventory(ctx, inventory)
	if err != nil {
		log.Error().Err(err).Msg("failed to create service")
		return nil, err
	}

	return external.NewService200JSONResponse{
		ServiceId: &inventory.ID,
	}, nil
}

func (h *httpHandler) NewSnippet(ctx context.Context, request external.NewSnippetRequestObject) (external.NewSnippetResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized")
		return nil, nil
	}

	if request.Body == nil || request.Body.Content == nil || *request.Body.Content == "" {
		return nil, nil
	}

	snippet := domain.Snippet{
		Content:   *request.Body.Content,
		CreatedAt: time.Now(),
	}

	err := h.service.CreateSnippet(ctx, snippet)
	if err != nil {
		log.Error().Err(err).Msg("failed to create snippet")
		return nil, err
	}

	return external.NewSnippet200JSONResponse{
		Content: request.Body.Content,
	}, nil
}

func (h *httpHandler) SearchGroup(ctx context.Context, request external.SearchGroupRequestObject) (external.SearchGroupResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized, returning empty response")
		return external.SearchGroup200JSONResponse{}, nil
	}

	opts := []domain.InventorySearchOption{
		domain.WithItemType(domain.ItemTypeGroup),
		domain.WithDeleted(false),
	}

	if request.Params.Limit != nil {
		opts = append(opts, domain.WithLimit(int(*request.Params.Limit)))
	}
	if request.Params.Offset != nil {
		opts = append(opts, domain.WithOffset(int(*request.Params.Offset)))
	}
	if request.Params.Filter != nil && *request.Params.Filter != "" {
		opts = append(opts, domain.WithFilter(*request.Params.Filter))
	}

	groups, err := h.service.SearchInventory(ctx, opts...)
	if err != nil {
		log.Error().Err(err).Msg("failed to search groups")
		return nil, err
	}

	response := make(external.GroupSearchResponse, len(groups))
	for i, g := range groups {
		response[i] = external.Group{
			GroupId:     g.ID,
			Name:        g.Name,
			Description: g.Description,
			Parent:      g.ParentID,
		}
	}

	return external.SearchGroup200JSONResponse(response), nil
}

func (h *httpHandler) SearchMessage(ctx context.Context, request external.SearchMessageRequestObject) (external.SearchMessageResponseObject, error) {
	return external.SearchMessage200JSONResponse{}, nil
}

func (h *httpHandler) SearchService(ctx context.Context, request external.SearchServiceRequestObject) (external.SearchServiceResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized, returning empty response")
		return external.SearchService200JSONResponse{}, nil
	}

	opts := []domain.InventorySearchOption{
		domain.WithItemType(domain.ItemTypeService),
		domain.WithDeleted(false),
	}

	if request.Params.Limit != nil {
		opts = append(opts, domain.WithLimit(int(*request.Params.Limit)))
	}
	if request.Params.Offset != nil {
		opts = append(opts, domain.WithOffset(int(*request.Params.Offset)))
	}
	if request.Params.Filter != nil && *request.Params.Filter != "" {
		opts = append(opts, domain.WithFilter(*request.Params.Filter))
	}

	services, err := h.service.SearchInventory(ctx, opts...)
	if err != nil {
		log.Error().Err(err).Msg("failed to search services")
		return nil, err
	}

	var serviceID *string
	var name *string
	var desc *string
	if len(services) > 0 {
		s := services[0]
		serviceID = new(s.ID.String())
		name = new(s.Name)
		desc = new(s.Description)
	}

	return external.SearchService200JSONResponse{
		ServiceId:   serviceID,
		Name:        name,
		Description: desc,
	}, nil
}

func (h *httpHandler) ServiceActiveMessage(ctx context.Context, request external.ServiceActiveMessageRequestObject) (external.ServiceActiveMessageResponseObject, error) {
	msg := "Service is active"
	return external.ServiceActiveMessage200JSONResponse{
		Message: &msg,
	}, nil
}

func (h *httpHandler) ServiceById(ctx context.Context, request external.ServiceByIdRequestObject) (external.ServiceByIdResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized, returning empty response")
		return external.ServiceById200JSONResponse{}, nil
	}

	serviceID, err := uuid.Parse(request.Service)
	if err != nil {
		return nil, err
	}

	inventory, err := h.service.GetInventory(ctx, serviceID)
	if err != nil {
		log.Error().Err(err).Str("serviceId", serviceID.String()).Msg("failed to get service")
		return nil, err
	}

	serviceIDStr := inventory.ID.String()
	return external.ServiceById200JSONResponse{
		ServiceId: &serviceIDStr,
	}, nil
}

func (h *httpHandler) ServiceMessage(ctx context.Context, request external.ServiceMessageRequestObject) (external.ServiceMessageResponseObject, error) {
	msg := "Service message"
	return external.ServiceMessage200JSONResponse{
		Message: &msg,
	}, nil
}

func (h *httpHandler) ServiceUpdate(ctx context.Context, request external.ServiceUpdateRequestObject) (external.ServiceUpdateResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized")
		return nil, nil
	}

	serviceID, err := uuid.Parse(request.Service)
	if err != nil {
		return nil, err
	}

	inventory, err := h.service.GetInventory(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if request.Body != nil {
		if request.Body.Name != nil {
			inventory.Name = *request.Body.Name
		}
		if request.Body.Description != nil {
			inventory.Description = *request.Body.Description
		}
	}

	inventory.UpdatedAt = time.Now()
	err = h.service.UpdateInventory(ctx, inventory)
	if err != nil {
		log.Error().Err(err).Msg("failed to update service")
		return nil, err
	}

	return external.ServiceUpdate200JSONResponse{}, nil
}

func (h *httpHandler) UpsertContext(ctx context.Context, request external.UpsertContextRequestObject) (external.UpsertContextResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized")
		return nil, nil
	}

	if request.Body == nil {
		return nil, nil
	}

	items := make([]domain.ConfigItem, len(*request.Body))
	for i, item := range *request.Body {
		if item.Entry != nil && item.Value != nil {
			items[i] = domain.ConfigItem{
				Key:   *item.Entry,
				Value: *item.Value,
			}
		}
	}

	err := h.service.UpsertItemContext(ctx, domain.ItemContext(request.ConfigContext), items)
	if err != nil {
		log.Error().Err(err).Str("context", request.ConfigContext).Msg("failed to upsert config context")
		return nil, err
	}

	return external.UpsertContext200JSONResponse{}, nil
}

func (h *httpHandler) UserSnippets(ctx context.Context, request external.UserSnippetsRequestObject) (external.UserSnippetsResponseObject, error) {
	if h.service == nil {
		log.Warn().Msg("service not initialized, returning empty response")
		return external.UserSnippets200JSONResponse{}, nil
	}

	snippets, err := h.service.UserSnippets(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to get snippets")
		return nil, err
	}

	response := make(external.SearchSnippetResponse, len(snippets))
	for i, s := range snippets {
		response[i] = external.SnippetResponse{
			Content: &s.Content,
		}
	}

	return external.UserSnippets200JSONResponse(response), nil
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

//go:fix inline
func ptr(s string) *string {
	return new(s)
}
