package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

// --- Asset Operations ---

func (s *Impl) CreateAsset(ctx context.Context, asset *domain.Asset) error {
	if err := s.assetRepo.CreateAsset(ctx, asset); err != nil {
		return err
	}
	if s.vecGen != nil {
		vec, err := s.vecGen.Generate(ctx, asset.Name)
		if err == nil {
			s.assetRepo.UpdateAssetEmbedding(ctx, asset.ID, vec)
		}
	}
	return nil
}

func (s *Impl) GetAsset(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	return s.assetRepo.GetAsset(ctx, id)
}

func (s *Impl) SearchAssets(ctx context.Context, assetType string, parentID *uuid.UUID, limit, offset int) ([]domain.Asset, error) {
	return s.assetRepo.SearchAssets(ctx, assetType, parentID, limit, offset)
}

func (s *Impl) HybridSearchAssets(ctx context.Context, query string, limit, offset int) ([]domain.Asset, error) {
	if s.vecGen != nil {
		vec, err := s.vecGen.Generate(ctx, query)
		if err == nil {
			return s.assetRepo.VectorSearchAssets(ctx, vec, limit, offset)
		}
	}
	return s.assetRepo.HybridSearchAssets(ctx, query, limit, offset)
}

// --- Identity Operations ---

func (s *Impl) CreateUser(ctx context.Context, user *domain.User) error {
	return s.assetRepo.CreateUser(ctx, user)
}

func (s *Impl) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.assetRepo.GetUser(ctx, id)
}

func (s *Impl) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	return s.assetRepo.UpdateUser(ctx, id, updates)
}

func (s *Impl) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.assetRepo.ListUsers(ctx)
}

func (s *Impl) CreateGroup(ctx context.Context, group *domain.UserGroup) error {
	return s.assetRepo.CreateGroup(ctx, group)
}

func (s *Impl) ListGroups(ctx context.Context) ([]domain.UserGroup, error) {
	return s.assetRepo.ListGroups(ctx)
}

func (s *Impl) AddGroupMember(ctx context.Context, groupID, userID uuid.UUID) error {
	return s.assetRepo.AddGroupMember(ctx, groupID, userID)
}

func (s *Impl) RemoveGroupMember(ctx context.Context, groupID, userID uuid.UUID) error {
	return s.assetRepo.RemoveGroupMember(ctx, groupID, userID)
}

func (s *Impl) ListGroupMemberships(ctx context.Context, groupID uuid.UUID) ([]domain.GroupMember, error) {
	return s.assetRepo.ListGroupMemberships(ctx, groupID)
}

// --- Vault Operations ---

func (s *Impl) AddVaultMember(ctx context.Context, vm *domain.VaultMembership) error {
	return s.assetRepo.CreateVaultMembership(ctx, vm)
}

func (s *Impl) ListVaultMembers(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMembership, error) {
	return s.assetRepo.ListVaultMemberships(ctx, vaultID)
}

// --- Admin Config Operations ---

func (s *Impl) GetConfig(ctx context.Context) ([]domain.AdminConfig, error) {
	return s.assetRepo.GetAdminConfig(ctx)
}

func (s *Impl) PatchConfig(ctx context.Context, configs map[string]string) error {
	for k, v := range configs {
		if !domain.IsValidConfigKey(k) {
			return domain.ValidationError{
				Field:   k,
				Message: "unknown config key: " + k,
				Code:    domain.ErrCodeValidation,
			}
		}
		if err := s.assetRepo.UpsertAdminConfig(ctx, k, v); err != nil {
			return err
		}
	}
	return nil
}

// --- Audit Operations ---

func (s *Impl) CreateAudit(ctx context.Context, audit *domain.AuditEntry) error {
	return s.auditRepo.Create(ctx, audit)
}

func (s *Impl) SearchAudits(ctx context.Context, actorID, resourceID *uuid.UUID, from *string) ([]domain.AuditEntry, error) {
	return s.auditRepo.Search(ctx, actorID, resourceID, from)
}
