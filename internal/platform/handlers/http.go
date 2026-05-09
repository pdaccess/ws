package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/pdaccess/commons/pkg/middleware"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/platform/handlers/external"
)

type httpHandler struct {
	svc ports.Service
}

func NewHttpHandler(svc ports.Service) external.StrictServerInterface {
	return &httpHandler{svc: svc}
}

// --- Admin Config ---

func (h *httpHandler) GetAdminConfig(ctx context.Context, request external.GetAdminConfigRequestObject) (external.GetAdminConfigResponseObject, error) {
	return external.GetAdminConfig200Response{}, nil
}

func (h *httpHandler) PatchAdminConfig(ctx context.Context, request external.PatchAdminConfigRequestObject) (external.PatchAdminConfigResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	configs := map[string]string{}
	if request.Body.JwtTtl != nil {
		configs["jwt_ttl"] = fmt.Sprintf("%d", *request.Body.JwtTtl)
	}
	if request.Body.NetworkWhitelist != nil {
		b, _ := json.Marshal(*request.Body.NetworkWhitelist)
		configs["network_whitelist"] = string(b)
	}
	if request.Body.RebuildPolicyCache != nil {
		configs["rebuild_policy_cache"] = fmt.Sprintf("%t", *request.Body.RebuildPolicyCache)
	}
	if err := h.svc.PatchConfig(ctx, configs); err != nil {
		return nil, err
	}
	return external.PatchAdminConfig200Response{}, nil
}

// --- Assets ---

func (h *httpHandler) GetAssets(ctx context.Context, request external.GetAssetsRequestObject) (external.GetAssetsResponseObject, error) {
	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}
	if request.Params.Offset != nil {
		offset = *request.Params.Offset
	}

	var assets []domain.Asset
	var err error

	if request.Params.Q != nil && *request.Params.Q != "" {
		assets, err = h.svc.HybridSearchAssets(ctx, *request.Params.Q, limit, offset)
	} else {
		var assetType string
		if request.Params.Type != nil {
			assetType = *request.Params.Type
		}
		var parentID *uuid.UUID
		if request.Params.ParentId != nil {
			id := uuid.UUID(*request.Params.ParentId)
			parentID = &id
		}
		assets, err = h.svc.SearchAssets(ctx, assetType, parentID, limit, offset)
	}
	if err != nil {
		return nil, err
	}

	data := make([]external.Asset, 0, len(assets))
	for _, a := range assets {
		data = append(data, domainAssetToExternal(a))
	}
	total := len(assets)
	return external.GetAssets200JSONResponse(external.AssetList{
		Data: &data,
		Meta: &external.PaginationMeta{
			Limit:  &limit,
			Offset: &offset,
			Total:  &total,
		},
	}), nil
}

func (h *httpHandler) PostAssets(ctx context.Context, request external.PostAssetsRequestObject) (external.PostAssetsResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	asset := &domain.Asset{
		Name:    request.Body.Name,
		Type:    domain.AssetType(request.Body.Type),
		OwnerID: userIDFromCtx(ctx),
	}

	specBytes, err := json.Marshal(request.Body.Spec)
	if err != nil {
		return nil, domain.ValidationError{Field: "spec", Message: "invalid spec", Code: domain.ErrCodeValidation}
	}
	asset.Spec = specBytes

	if err := h.svc.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return external.PostAssets201JSONResponse(domainAssetToExternal(*asset)), nil
}

func (h *httpHandler) PostAssetsIdActions(ctx context.Context, request external.PostAssetsIdActionsRequestObject) (external.PostAssetsIdActionsResponseObject, error) {
	return external.PostAssetsIdActions200Response{}, nil
}

// --- Audit Logs ---

func (h *httpHandler) GetAuditLogs(ctx context.Context, request external.GetAuditLogsRequestObject) (external.GetAuditLogsResponseObject, error) {
	var actorID, resourceID *uuid.UUID
	var from *string

	if request.Params.ActorId != nil {
		id := uuid.UUID(*request.Params.ActorId)
		actorID = &id
	}
	if request.Params.ResourceId != nil {
		id := uuid.UUID(*request.Params.ResourceId)
		resourceID = &id
	}
	if request.Params.From != nil {
		s := request.Params.From.Format("2006-01-02T15:04:05Z")
		from = &s
	}

	entries, err := h.svc.SearchAudits(ctx, actorID, resourceID, from)
	if err != nil {
		return nil, err
	}

	result := make(external.GetAuditLogs200JSONResponse, 0, len(entries))
	for _, e := range entries {
		result = append(result, domainAuditEntryToExternal(e))
	}
	return result, nil
}

// --- Identity: Groups ---

func (h *httpHandler) GetIdentityGroups(ctx context.Context, request external.GetIdentityGroupsRequestObject) (external.GetIdentityGroupsResponseObject, error) {
	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}
	if request.Params.Offset != nil {
		offset = *request.Params.Offset
	}

	groups, err := h.svc.ListGroups(ctx)
	if err != nil {
		return nil, err
	}

	// Slice the groups for pagination
	start := min(offset, len(groups))
	end := min(start+limit, len(groups))
	paginated := groups[start:end]

	data := make([]external.UserGroup, 0, len(paginated))
	for _, g := range paginated {
		data = append(data, external.UserGroup{
			Id:          openapiUUID(g.ID),
			Name:        g.Name,
			Description: &g.Description,
		})
	}
	total := len(groups)
	return external.GetIdentityGroups200JSONResponse(external.UserGroupList{
		Data: &data,
		Meta: &external.PaginationMeta{
			Limit:  &limit,
			Offset: &offset,
			Total:  &total,
		},
	}), nil
}

func (h *httpHandler) PostIdentityGroups(ctx context.Context, request external.PostIdentityGroupsRequestObject) (external.PostIdentityGroupsResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	group := &domain.UserGroup{
		Name: request.Body.Name,
	}
	if err := h.svc.CreateGroup(ctx, group); err != nil {
		return nil, err
	}
	return external.PostIdentityGroups201Response{}, nil
}

func (h *httpHandler) DeleteIdentityGroupsGroupIdMemberships(ctx context.Context, request external.DeleteIdentityGroupsGroupIdMembershipsRequestObject) (external.DeleteIdentityGroupsGroupIdMembershipsResponseObject, error) {
	if err := h.svc.RemoveGroupMember(ctx, uuid.UUID(request.GroupId), uuid.UUID(request.Params.UserId)); err != nil {
		return nil, err
	}
	return external.DeleteIdentityGroupsGroupIdMemberships204Response{}, nil
}

func (h *httpHandler) GetIdentityGroupsGroupIdMemberships(ctx context.Context, request external.GetIdentityGroupsGroupIdMembershipsRequestObject) (external.GetIdentityGroupsGroupIdMembershipsResponseObject, error) {
	groupID := uuid.UUID(request.GroupId)
	limit := 20
	offset := 0
	if request.Params.Limit != nil {
		limit = *request.Params.Limit
	}
	if request.Params.Offset != nil {
		offset = *request.Params.Offset
	}

	members, err := h.svc.ListGroupMemberships(ctx, groupID)
	if err != nil {
		return nil, err
	}

	start := min(offset, len(members))
	end := min(start+limit, len(members))
	paginated := members[start:end]

	data := make([]external.GroupMembership, 0, len(paginated))
	for _, m := range paginated {
		groupID := openapiUUID(m.GroupID)
		userID := openapiUUID(m.MemberID)
		data = append(data, external.GroupMembership{
			GroupId:   &groupID,
			UserId:    &userID,
			CreatedAt: &m.CreatedAt,
		})
	}
	total := len(members)
	return external.GetIdentityGroupsGroupIdMemberships200JSONResponse(external.GroupMembershipList{
		Data: &data,
		Meta: &external.PaginationMeta{
			Limit:  &limit,
			Offset: &offset,
			Total:  &total,
		},
	}), nil
}

func (h *httpHandler) PostIdentityGroupsGroupIdMemberships(ctx context.Context, request external.PostIdentityGroupsGroupIdMembershipsRequestObject) (external.PostIdentityGroupsGroupIdMembershipsResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	if err := h.svc.AddGroupMember(ctx, uuid.UUID(request.GroupId), uuid.UUID(request.Body.UserId)); err != nil {
		return nil, err
	}
	return external.PostIdentityGroupsGroupIdMemberships201Response{}, nil
}

// --- Identity: Users ---

func (h *httpHandler) GetIdentityUser(ctx context.Context, request external.GetIdentityUserRequestObject) (external.GetIdentityUserResponseObject, error) {
	var id uuid.UUID
	if request.Params.UserId != nil {
		id = uuid.UUID(*request.Params.UserId)
	} else {
		id = userIDFromCtx(ctx)
	}

	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		if request.Params.UserId != nil {
			return external.GetIdentityUser404Response{}, nil
		}
		user = &domain.User{ID: id, Username: id.String(), Email: id.String() + "@pdaccess.io"}
		if err := h.svc.CreateUser(ctx, user); err != nil {
			return nil, err
		}
	}

	status := external.Active
	if user.Status != "" {
		status = external.UserStatus(user.Status)
	}

	extUser := external.User{
		Id:       openapiUUID(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Status:   &status,
	}

	if request.Params.View != nil && *request.Params.View == external.Full {
		displayName := user.DisplayName
		firstName := user.FirstName
		lastName := user.LastName
		extUser.DisplayName = &displayName
		extUser.FirstName = &firstName
		extUser.LastName = &lastName
		if len(user.NotificationSettings) > 0 {
			var notif map[string]any
			if json.Unmarshal(user.NotificationSettings, &notif) == nil {
				extUser.NotificationSettings = &notif
			}
		}
	}

	return external.GetIdentityUser200JSONResponse(extUser), nil
}

func (h *httpHandler) PostIdentityUser(ctx context.Context, request external.PostIdentityUserRequestObject) (external.PostIdentityUserResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}

	id := userIDFromCtx(ctx)

	user, err := h.svc.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		user = &domain.User{ID: id, Username: id.String(), Email: id.String() + "@pdaccess.io"}
		if err := h.svc.CreateUser(ctx, user); err != nil {
			return nil, err
		}
	}

	updates := map[string]any{}

	if request.Body.DisplayName != nil {
		updates["display_name"] = *request.Body.DisplayName
	}
	if request.Body.FirstName != nil {
		updates["first_name"] = *request.Body.FirstName
	}
	if request.Body.LastName != nil {
		updates["last_name"] = *request.Body.LastName
	}
	if request.Body.Email != nil {
		updates["email"] = *request.Body.Email
	}
	if request.Body.NotificationSettings != nil {
		b, err := json.Marshal(*request.Body.NotificationSettings)
		if err != nil {
			return nil, domain.ValidationError{Field: "notification_settings", Message: "invalid format", Code: domain.ErrCodeValidation}
		}
		updates["notification_settings"] = string(b)
	}

	if err := h.svc.UpdateUser(ctx, id, updates); err != nil {
		return nil, err
	}

	user, err = h.svc.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	status := external.Active
	if user.Status != "" {
		status = external.UserStatus(user.Status)
	}

	extUser := external.User{
		Id:       openapiUUID(id),
		Username: user.Username,
		Email:    user.Email,
		Status:   &status,
	}
	displayName := user.DisplayName
	firstName := user.FirstName
	lastName := user.LastName
	extUser.DisplayName = &displayName
	extUser.FirstName = &firstName
	extUser.LastName = &lastName
	if len(user.NotificationSettings) > 0 {
		var notif map[string]any
		if json.Unmarshal(user.NotificationSettings, &notif) == nil {
			extUser.NotificationSettings = &notif
		}
	}

	return external.PostIdentityUser200JSONResponse(extUser), nil
}

func (h *httpHandler) GetIdentityUsers(ctx context.Context, request external.GetIdentityUsersRequestObject) (external.GetIdentityUsersResponseObject, error) {
	userID := userIDFromCtx(ctx)
	status := external.Active
	userIdStr := userID.String()
	return external.GetIdentityUsers200JSONResponse{
		external.User{
			Id:       openapiUUID(userID),
			Username: userIdStr,
			Email:    userIdStr + "@pdaccess.io",
			Status:   &status,
		},
	}, nil
}

func (h *httpHandler) PostIdentityUsers(ctx context.Context, request external.PostIdentityUsersRequestObject) (external.PostIdentityUsersResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	user := &domain.User{
		Username: request.Body.Username,
		Email:    request.Body.Email,
	}
	if err := h.svc.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return external.PostIdentityUsers201Response{}, nil
}

// --- Vault Memberships ---

func (h *httpHandler) GetVaultsVaultIdMemberships(ctx context.Context, request external.GetVaultsVaultIdMembershipsRequestObject) (external.GetVaultsVaultIdMembershipsResponseObject, error) {
	members, err := h.svc.ListVaultMembers(ctx, uuid.UUID(request.VaultId))
	if err != nil {
		return nil, err
	}
	_ = members
	return external.GetVaultsVaultIdMemberships200Response{}, nil
}

func (h *httpHandler) PostVaultsVaultIdMemberships(ctx context.Context, request external.PostVaultsVaultIdMembershipsRequestObject) (external.PostVaultsVaultIdMembershipsResponseObject, error) {
	if request.Body == nil {
		return nil, domain.ValidationError{Field: "body", Message: "missing request body", Code: domain.ErrCodeValidation}
	}
	vm := &domain.VaultMembership{
		VaultID:    uuid.UUID(request.VaultId),
		MemberID:   uuid.UUID(request.Body.MemberId),
		MemberType: string(request.Body.MemberType),
		Role:       string(request.Body.Role),
	}
	if err := h.svc.AddVaultMember(ctx, vm); err != nil {
		return nil, err
	}
	return external.PostVaultsVaultIdMemberships201Response{}, nil
}

func userIDFromCtx(ctx context.Context) uuid.UUID {
	client := middleware.ClientFromCtx(ctx)
	if client == nil {
		return uuid.Nil
	}
	id, err := uuid.Parse(client.User.UserId)
	if err != nil {
		return uuid.NewSHA1(uuid.NameSpaceOID, []byte(client.User.UserId))
	}
	return id
}

// --- Helpers ---

func openapiUUID(id uuid.UUID) openapi_types.UUID {
	return openapi_types.UUID(id)
}

func domainAssetToExternal(a domain.Asset) external.Asset {
	id := openapiUUID(a.ID)
	ownerID := openapiUUID(a.OwnerID)
	t := string(a.Type)
	var spec map[string]any
	if len(a.Spec) > 0 {
		json.Unmarshal(a.Spec, &spec)
	}
	var parentID *openapi_types.UUID
	if a.ParentID != nil {
		id := openapiUUID(*a.ParentID)
		parentID = &id
	}
	return external.Asset{
		Id:        &id,
		Name:      &a.Name,
		Type:      &t,
		OwnerId:   &ownerID,
		ParentId:  parentID,
		Spec:      &spec,
		CreatedAt: &a.CreatedAt,
	}
}

func domainAuditEntryToExternal(e domain.AuditEntry) external.AuditEntry {
	actorID := openapiUUID(e.ActorID)
	resourceID := openapiUUID(e.ResourceID)
	success := e.Success
	return external.AuditEntry{
		Action:     &e.Action,
		ActorId:    &actorID,
		ResourceId: &resourceID,
		Success:    &success,
		Timestamp:  &e.Timestamp,
	}
}
