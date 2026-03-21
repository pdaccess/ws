package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

type UserOperations interface {
	SearchUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type ConfigOperations interface {
	UpsertItemContext(ctx context.Context, context domain.ItemContext, items []domain.ConfigItem) error
	GetItemContext(ctx context.Context, context domain.ItemContext) ([]domain.ConfigItem, error)
}

type SnippetOperations interface {
	UserSnippets(ctx context.Context, options ...domain.SnippetSearchOption) ([]domain.Snippet, error)
	CreateSnippet(ctx context.Context, snippet domain.Snippet) error
}

type GroupOperations interface {
	CreateGroup(ctx context.Context, group *domain.Group, userID, realmID uuid.UUID) error
	GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, error)
	UpdateGroup(ctx context.Context, group *domain.Group, userID, realmID uuid.UUID) error
	DeleteGroup(ctx context.Context, id uuid.UUID, userID, realmID uuid.UUID) error
	SearchGroups(ctx context.Context, opts ...domain.GroupSearchOption) ([]domain.Group, error)
	SearchGroupsWithQuery(ctx context.Context, query string, limit, offset int) ([]domain.Group, error)

	AddGroupMember(ctx context.Context, member *domain.GroupMember, userID, realmID uuid.UUID) error
	RemoveGroupMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID, userID, realmID uuid.UUID) error
	GetGroupMembers(ctx context.Context, groupID uuid.UUID, limit, offset int) ([]domain.GroupMember, error)
}

type ServiceOperations interface {
	CreateService(ctx context.Context, svc *domain.Service, userID, realmID uuid.UUID) error
	GetService(ctx context.Context, id uuid.UUID) (*domain.Service, error)
	UpdateService(ctx context.Context, svc *domain.Service, userID, realmID uuid.UUID) error
	DeleteService(ctx context.Context, id uuid.UUID, userID, realmID uuid.UUID) error
	SearchServices(ctx context.Context, opts ...domain.ServiceSearchOption) ([]domain.Service, error)
	SearchServicesWithQuery(ctx context.Context, query string, limit, offset int) ([]domain.Service, error)

	AddServiceMember(ctx context.Context, member *domain.ServiceMember, userID, realmID uuid.UUID) error
	RemoveServiceMembers(ctx context.Context, serviceID uuid.UUID, userIDs []uuid.UUID, userID, realmID uuid.UUID) error
	GetServiceMembers(ctx context.Context, serviceID uuid.UUID, limit, offset int) ([]domain.ServiceMember, error)

	UpsertServiceSettings(ctx context.Context, settings *domain.ServiceSettings, userID, realmID uuid.UUID) error
	GetServiceSettings(ctx context.Context, serviceID uuid.UUID) (*domain.ServiceSettings, error)
}

type AlarmOperations interface {
	CreateAlarm(ctx context.Context, alarm *domain.Alarm) error
	GetAlarm(ctx context.Context, id uuid.UUID) (*domain.Alarm, error)
	DeleteAlarm(ctx context.Context, id uuid.UUID) error
	SearchAlarms(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Alarm, error)
	AcknowledgeAlarm(ctx context.Context, id uuid.UUID) error
}

type ActivityOperations interface {
	CreateActivity(ctx context.Context, activity *domain.Activity) error
	SearchActivities(ctx context.Context, opts ...domain.ActivitySearchOption) ([]domain.Activity, error)
	GetActivitiesByResourceID(ctx context.Context, resourceID uuid.UUID, limit int) ([]domain.Activity, error)
}

type PasteOperations interface {
	CreatePaste(ctx context.Context, paste *domain.Paste) error
	GetPaste(ctx context.Context, id uuid.UUID) (*domain.Paste, error)
	DeletePaste(ctx context.Context, id uuid.UUID) error
	SearchPastes(ctx context.Context, opts ...domain.PasteSearchOption) ([]domain.Paste, error)
}

type UserGroupOperations interface {
	CreateUserGroup(ctx context.Context, ug *domain.UserGroup, userID, realmID uuid.UUID) error
	GetUserGroup(ctx context.Context, id uuid.UUID) (*domain.UserGroup, error)
	UpdateUserGroup(ctx context.Context, ug *domain.UserGroup, userID, realmID uuid.UUID) error
	DeleteUserGroup(ctx context.Context, id uuid.UUID, userID, realmID uuid.UUID) error
	SearchUserGroups(ctx context.Context, limit, offset int) ([]domain.UserGroup, error)

	AddUserGroupMember(ctx context.Context, member *domain.UserGroupMember, userID, realmID uuid.UUID) error
	RemoveUserGroupMembers(ctx context.Context, userGroupID uuid.UUID, userIDs []uuid.UUID, userID, realmID uuid.UUID) error
	GetUserGroupMembers(ctx context.Context, userGroupID uuid.UUID) ([]domain.UserGroupMember, error)
}

type Service interface {
	UserOperations
	GroupOperations
	ServiceOperations
	AlarmOperations
	ActivityOperations
	PasteOperations
	UserGroupOperations

	ConfigOperations
	SnippetOperations
}
