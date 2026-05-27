package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

// --- Asset Operations ---

func (s *Impl) CreateAsset(ctx context.Context, asset *domain.Asset) error {
	if asset.Type == domain.AssetTypePolicy {
		if err := s.validatePolicyRefs(ctx, asset); err != nil {
			return err
		}
	}

	if err := s.assetRepo.CreateAsset(ctx, asset); err != nil {
		return err
	}

	if asset.Type == domain.AssetTypePolicy {
		return s.reconcilePolicyForAllAssets(ctx, asset)
	}

	return s.reconcileEffectivePolicies(ctx, asset)
}

func (s *Impl) validatePolicyRefs(ctx context.Context, asset *domain.Asset) error {
	var ps struct {
		Subjects struct {
			Users  []string `json:"users"`
			Groups []string `json:"groups"`
		} `json:"subjects"`
		Objects struct {
			AssetIDs []string `json:"asset_ids"`
		} `json:"objects"`
	}
	if err := json.Unmarshal(asset.Spec, &ps); err != nil {
		return domain.ValidationError{Field: "spec", Message: "invalid policy spec format", Code: domain.ErrCodeValidation}
	}

	var missing []string

	for _, uidStr := range ps.Subjects.Users {
		uid, err := uuid.Parse(uidStr)
		if err != nil {
			missing = append(missing, fmt.Sprintf("user(%s)", uidStr))
			continue
		}
		u, err := s.assetRepo.GetUser(ctx, uid)
		if err != nil {
			return err
		}
		if u == nil {
			missing = append(missing, fmt.Sprintf("user(%s)", uidStr))
		}
	}

	for _, gidStr := range ps.Subjects.Groups {
		gid, err := uuid.Parse(gidStr)
		if err != nil {
			missing = append(missing, fmt.Sprintf("group(%s)", gidStr))
			continue
		}
		g, err := s.assetRepo.GetGroup(ctx, gid)
		if err != nil {
			return err
		}
		if g == nil {
			missing = append(missing, fmt.Sprintf("group(%s)", gidStr))
		}
	}

	for _, aidStr := range ps.Objects.AssetIDs {
		aid, err := uuid.Parse(aidStr)
		if err != nil {
			missing = append(missing, fmt.Sprintf("asset(%s)", aidStr))
			continue
		}
		a, err := s.assetRepo.GetAsset(ctx, aid)
		if err != nil {
			return err
		}
		if a == nil {
			missing = append(missing, fmt.Sprintf("asset(%s)", aidStr))
		}
	}

	if len(missing) > 0 {
		return domain.ValidationError{
			Field:   "spec",
			Message: fmt.Sprintf("referenced resources not found: %s", strings.Join(missing, ", ")),
			Code:    domain.ErrCodeValidation,
		}
	}

	return nil
}

func (s *Impl) processPolicyForAsset(ctx context.Context, policy *domain.Asset, asset *domain.Asset) []domain.EffectivePolicy {
	var effective []domain.EffectivePolicy

	var ps struct {
		Subjects struct {
			Users  []string `json:"users"`
			Groups []string `json:"groups"`
		} `json:"subjects"`
		Actions []string `json:"actions"`
		Objects struct {
			AssetIDs []string       `json:"asset_ids"`
			Tags     map[string]any `json:"tags"`
		} `json:"objects"`
	}
	if err := json.Unmarshal(policy.Spec, &ps); err != nil {
		return nil
	}

	if !matchesPolicyObjects(asset, ps.Objects) {
		return nil
	}

	for _, uidStr := range ps.Subjects.Users {
		uid, err := uuid.Parse(uidStr)
		if err != nil {
			continue
		}
		effective = append(effective, domain.EffectivePolicy{
			UserID:   uid,
			AssetID:  asset.ID,
			Actions:  ps.Actions,
			PolicyID: policy.ID,
		})
	}

	for _, gidStr := range ps.Subjects.Groups {
		gid, err := uuid.Parse(gidStr)
		if err != nil {
			continue
		}
		members, err := s.assetRepo.ListGroupMemberships(ctx, gid)
		if err != nil {
			continue
		}
		for _, m := range members {
			groupID := gid
			effective = append(effective, domain.EffectivePolicy{
				UserID:   m.MemberID,
				GroupID:  &groupID,
				AssetID:  asset.ID,
				Actions:  ps.Actions,
				PolicyID: policy.ID,
			})
		}
	}

	return effective
}

func (s *Impl) reconcileEffectivePolicies(ctx context.Context, asset *domain.Asset) error {
	policies, err := s.assetRepo.ListPolicies(ctx)
	if err != nil {
		return nil
	}

	var effective []domain.EffectivePolicy
	for i := range policies {
		effective = append(effective, s.processPolicyForAsset(ctx, &policies[i], asset)...)
	}

	if len(effective) > 0 {
		return s.assetRepo.SetAssetEffectivePolicies(ctx, asset.ID, effective)
	}
	return nil
}

func (s *Impl) reconcilePolicyForAllAssets(ctx context.Context, policy *domain.Asset) error {
	assets, err := s.assetRepo.ListNonPolicyAssets(ctx)
	if err != nil {
		return nil
	}

	for i := range assets {
		effective := s.processPolicyForAsset(ctx, policy, &assets[i])
		if len(effective) > 0 {
			if err := s.assetRepo.AddAssetEffectivePolicies(ctx, assets[i].ID, effective); err != nil {
				return err
			}
		}
	}

	return nil
}

func matchesPolicyObjects(asset *domain.Asset, objs struct {
	AssetIDs []string       `json:"asset_ids"`
	Tags     map[string]any `json:"tags"`
}) bool {
	for _, idStr := range objs.AssetIDs {
		id, err := uuid.Parse(idStr)
		if err == nil && id == asset.ID {
			return true
		}
	}

	if len(objs.Tags) > 0 && len(asset.Spec) > 0 {
		var specMap map[string]any
		if err := json.Unmarshal(asset.Spec, &specMap); err != nil {
			return false
		}
		assetTags, ok := specMap["tags"].(string)
		if !ok {
			return false
		}
		for _, v := range objs.Tags {
			tagVal, ok := v.(string)
			if !ok {
				continue
			}
			if containsTag(assetTags, tagVal) {
				return true
			}
		}
	}

	return false
}

func containsTag(tagList, tag string) bool {
	for i := 0; i < len(tagList); {
		end := i + 1
		for end < len(tagList) && tagList[end] != ',' {
			end++
		}
		part := tagList[i:end]
		if end < len(tagList) {
			end++
		}
		i = end
		for len(part) > 0 && part[0] == ' ' {
			part = part[1:]
		}
		for len(part) > 0 && part[len(part)-1] == ' ' {
			part = part[:len(part)-1]
		}
		if part == tag {
			return true
		}
	}
	return false
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
