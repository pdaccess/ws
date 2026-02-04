package handlers

import (
	"context"

	"git.h2hsecure.com/core/ws/internal/handlers/external"
)

type httpHandler struct{}

func NewHttpHandler() external.StrictServerInterface {
	return &httpHandler{}
}

// AlarmIndex implements external.StrictServerInterface.
func (h *httpHandler) AlarmIndex(ctx context.Context, request external.AlarmIndexRequestObject) (external.AlarmIndexResponseObject, error) {
	panic("unimplemented")
}

// CreateGroupAlarm implements external.StrictServerInterface.
func (h *httpHandler) CreateGroupAlarm(ctx context.Context, request external.CreateGroupAlarmRequestObject) (external.CreateGroupAlarmResponseObject, error) {
	panic("unimplemented")
}

// CreateGroupMember implements external.StrictServerInterface.
func (h *httpHandler) CreateGroupMember(ctx context.Context, request external.CreateGroupMemberRequestObject) (external.CreateGroupMemberResponseObject, error) {
	panic("unimplemented")
}

// CreateGroupMessage implements external.StrictServerInterface.
func (h *httpHandler) CreateGroupMessage(ctx context.Context, request external.CreateGroupMessageRequestObject) (external.CreateGroupMessageResponseObject, error) {
	panic("unimplemented")
}

// DeleteGroup implements external.StrictServerInterface.
func (h *httpHandler) DeleteGroup(ctx context.Context, request external.DeleteGroupRequestObject) (external.DeleteGroupResponseObject, error) {
	panic("unimplemented")
}

// DeleteGroupAlarm implements external.StrictServerInterface.
func (h *httpHandler) DeleteGroupAlarm(ctx context.Context, request external.DeleteGroupAlarmRequestObject) (external.DeleteGroupAlarmResponseObject, error) {
	panic("unimplemented")
}

// DeleteGroupMembers implements external.StrictServerInterface.
func (h *httpHandler) DeleteGroupMembers(ctx context.Context, request external.DeleteGroupMembersRequestObject) (external.DeleteGroupMembersResponseObject, error) {
	panic("unimplemented")
}

// DeleteGroupMessage implements external.StrictServerInterface.
func (h *httpHandler) DeleteGroupMessage(ctx context.Context, request external.DeleteGroupMessageRequestObject) (external.DeleteGroupMessageResponseObject, error) {
	panic("unimplemented")
}

// DeleteService implements external.StrictServerInterface.
func (h *httpHandler) DeleteService(ctx context.Context, request external.DeleteServiceRequestObject) (external.DeleteServiceResponseObject, error) {
	panic("unimplemented")
}

// DeleteSnippet implements external.StrictServerInterface.
func (h *httpHandler) DeleteSnippet(ctx context.Context, request external.DeleteSnippetRequestObject) (external.DeleteSnippetResponseObject, error) {
	panic("unimplemented")
}

// FetchContext implements external.StrictServerInterface.
func (h *httpHandler) FetchContext(ctx context.Context, request external.FetchContextRequestObject) (external.FetchContextResponseObject, error) {
	panic("unimplemented")
}

// FetchGroupAlarms implements external.StrictServerInterface.
func (h *httpHandler) FetchGroupAlarms(ctx context.Context, request external.FetchGroupAlarmsRequestObject) (external.FetchGroupAlarmsResponseObject, error) {
	panic("unimplemented")
}

// GetConfiguration implements external.StrictServerInterface.
func (h *httpHandler) GetConfiguration(ctx context.Context, request external.GetConfigurationRequestObject) (external.GetConfigurationResponseObject, error) {
	panic("unimplemented")
}

// GetSnippet implements external.StrictServerInterface.
func (h *httpHandler) GetSnippet(ctx context.Context, request external.GetSnippetRequestObject) (external.GetSnippetResponseObject, error) {
	panic("unimplemented")
}

// GroupMembers implements external.StrictServerInterface.
func (h *httpHandler) GroupMembers(ctx context.Context, request external.GroupMembersRequestObject) (external.GroupMembersResponseObject, error) {
	panic("unimplemented")
}

// NewGroup implements external.StrictServerInterface.
func (h *httpHandler) NewGroup(ctx context.Context, request external.NewGroupRequestObject) (external.NewGroupResponseObject, error) {
	panic("unimplemented")
}

// NewService implements external.StrictServerInterface.
func (h *httpHandler) NewService(ctx context.Context, request external.NewServiceRequestObject) (external.NewServiceResponseObject, error) {
	panic("unimplemented")
}

// NewSnippet implements external.StrictServerInterface.
func (h *httpHandler) NewSnippet(ctx context.Context, request external.NewSnippetRequestObject) (external.NewSnippetResponseObject, error) {
	panic("unimplemented")
}

// SearchGroup implements external.StrictServerInterface.
func (h *httpHandler) SearchGroup(ctx context.Context, request external.SearchGroupRequestObject) (external.SearchGroupResponseObject, error) {
	panic("unimplemented")
}

// SearchMessage implements external.StrictServerInterface.
func (h *httpHandler) SearchMessage(ctx context.Context, request external.SearchMessageRequestObject) (external.SearchMessageResponseObject, error) {
	panic("unimplemented")
}

// SearchService implements external.StrictServerInterface.
func (h *httpHandler) SearchService(ctx context.Context, request external.SearchServiceRequestObject) (external.SearchServiceResponseObject, error) {
	panic("unimplemented")
}

// ServiceActiveMessage implements external.StrictServerInterface.
func (h *httpHandler) ServiceActiveMessage(ctx context.Context, request external.ServiceActiveMessageRequestObject) (external.ServiceActiveMessageResponseObject, error) {
	panic("unimplemented")
}

// ServiceById implements external.StrictServerInterface.
func (h *httpHandler) ServiceById(ctx context.Context, request external.ServiceByIdRequestObject) (external.ServiceByIdResponseObject, error) {
	panic("unimplemented")
}

// ServiceMessage implements external.StrictServerInterface.
func (h *httpHandler) ServiceMessage(ctx context.Context, request external.ServiceMessageRequestObject) (external.ServiceMessageResponseObject, error) {
	panic("unimplemented")
}

// ServiceUpdate implements external.StrictServerInterface.
func (h *httpHandler) ServiceUpdate(ctx context.Context, request external.ServiceUpdateRequestObject) (external.ServiceUpdateResponseObject, error) {
	panic("unimplemented")
}

// UpsertContext implements external.StrictServerInterface.
func (h *httpHandler) UpsertContext(ctx context.Context, request external.UpsertContextRequestObject) (external.UpsertContextResponseObject, error) {
	panic("unimplemented")
}

// UserSnippets implements external.StrictServerInterface.
func (h *httpHandler) UserSnippets(ctx context.Context, request external.UserSnippetsRequestObject) (external.UserSnippetsResponseObject, error) {
	panic("unimplemented")
}
