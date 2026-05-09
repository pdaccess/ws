package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

type AssetOperations interface {
	CreateAsset(ctx context.Context, asset *domain.Asset) error
	GetAsset(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	SearchAssets(ctx context.Context, assetType string, parentID *uuid.UUID, limit, offset int) ([]domain.Asset, error)
	HybridSearchAssets(ctx context.Context, query string, limit, offset int) ([]domain.Asset, error)
}

type IdentityOperations interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]any) error
	ListUsers(ctx context.Context) ([]domain.User, error)
	CreateGroup(ctx context.Context, group *domain.UserGroup) error
	ListGroups(ctx context.Context) ([]domain.UserGroup, error)
	AddGroupMember(ctx context.Context, groupID, userID uuid.UUID) error
	RemoveGroupMember(ctx context.Context, groupID, userID uuid.UUID) error
	ListGroupMemberships(ctx context.Context, groupID uuid.UUID) ([]domain.GroupMember, error)
}

type VaultOperations interface {
	AddVaultMember(ctx context.Context, vm *domain.VaultMembership) error
	ListVaultMembers(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMembership, error)
}

type AdminConfigOperations interface {
	GetConfig(ctx context.Context) ([]domain.AdminConfig, error)
	PatchConfig(ctx context.Context, configs map[string]string) error
}

type AuditOperations interface {
	CreateAudit(ctx context.Context, audit *domain.AuditEntry) error
	SearchAudits(ctx context.Context, actorID, resourceID *uuid.UUID, from *string) ([]domain.AuditEntry, error)
}

type Service interface {
	AssetOperations
	IdentityOperations
	VaultOperations
	AdminConfigOperations
	AuditOperations
}
