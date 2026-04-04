package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/rs/zerolog/log"
)

func (s *Impl) CreateGroup(ctx context.Context, group *domain.Group, userID, realmID uuid.UUID) error {
	err := s.inventoryRepo.CreateGroup(ctx, group)
	if err != nil {
		return err
	}

	s.logActivity(ctx, userID, realmID, "create", "group", group.ID, fmt.Sprintf("Created group: %s", group.Name))
	return nil
}

func (s *Impl) GetGroup(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	return s.inventoryRepo.GetGroupByID(ctx, id)
}

func (s *Impl) UpdateGroup(ctx context.Context, group *domain.Group, userID, realmID uuid.UUID) error {
	err := s.inventoryRepo.UpdateGroup(ctx, group)
	if err != nil {
		return err
	}

	s.logActivity(ctx, userID, realmID, "update", "group", group.ID, fmt.Sprintf("Updated group: %s", group.Name))
	return nil
}

func (s *Impl) DeleteGroup(ctx context.Context, id uuid.UUID, userID, realmID uuid.UUID) error {
	group, err := s.inventoryRepo.GetGroupByID(ctx, id)
	if err != nil {
		return err
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	err = s.inventoryRepo.DeleteGroup(ctx, id)
	if err != nil {
		return err
	}

	s.logActivity(ctx, userID, realmID, "delete", "group", id, fmt.Sprintf("Deleted group: %s", group.Name))
	return nil
}

func (s *Impl) SearchGroups(ctx context.Context, opts ...domain.GroupSearchOption) ([]domain.Group, error) {
	return s.inventoryRepo.SearchGroups(ctx, opts...)
}

func (s *Impl) AddGroupMember(ctx context.Context, member *domain.GroupMember, userID, realmID uuid.UUID) error {
	err := s.inventoryRepo.AddGroupMember(ctx, member)
	if err != nil {
		return err
	}

	details, _ := json.Marshal(map[string]any{
		"groupId": member.GroupID,
		"userId":  member.UserID,
		"role":    member.Role,
	})
	s.logActivity(ctx, userID, realmID, "add_member", "group", member.GroupID, string(details))
	return nil
}

func (s *Impl) RemoveGroupMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID, userID, realmID uuid.UUID) error {
	err := s.inventoryRepo.RemoveGroupMembers(ctx, groupID, userIDs)
	if err != nil {
		return err
	}

	details, _ := json.Marshal(map[string]any{
		"groupId":      groupID,
		"removedUsers": userIDs,
	})
	s.logActivity(ctx, userID, realmID, "remove_members", "group", groupID, string(details))
	return nil
}

func (s *Impl) GetGroupMembers(ctx context.Context, groupID uuid.UUID, limit, offset int) ([]domain.GroupMember, error) {
	return s.inventoryRepo.GetGroupMembers(ctx, groupID, limit, offset)
}

func (s *Impl) SearchGroupsVector(ctx context.Context, vector domain.Vector, limit, offset int) ([]domain.Group, error) {
	return s.inventoryRepo.SearchGroups(ctx,
		domain.WithGroupVector(vector),
		domain.WithGroupLimit(limit),
		domain.WithGroupOffset(offset),
	)
}

func (s *Impl) SearchGroupsWithQuery(ctx context.Context, query string, limit, offset int) ([]domain.Group, error) {
	if s.vector != nil {
		vector, err := s.vector.Generate(ctx, query)
		if err != nil {
			log.Warn().Err(err).Msg("vector generation failed, falling back to text search")
			return s.inventoryRepo.SearchGroups(ctx,
				domain.WithGroupFilter(query),
				domain.WithGroupLimit(limit),
				domain.WithGroupOffset(offset),
			)
		}
		if vector == nil || len(vector) == 0 {
			return s.inventoryRepo.SearchGroups(ctx,
				domain.WithGroupFilter(query),
				domain.WithGroupLimit(limit),
				domain.WithGroupOffset(offset),
			)
		}
		return s.SearchGroupsVector(ctx, vector, limit, offset)
	}

	return s.inventoryRepo.SearchGroups(ctx,
		domain.WithGroupFilter(query),
		domain.WithGroupLimit(limit),
		domain.WithGroupOffset(offset),
	)
}

func (s *Impl) logActivity(ctx context.Context, userID, realmID uuid.UUID, action, resource string, resourceID uuid.UUID, details string) {
	if s.activityRepo == nil {
		log.Warn().Msg("activity repository is nil, skipping activity log")
		return
	}

	activity := &domain.Activity{
		ID:         uuid.New(),
		UserID:     userID,
		RealmID:    realmID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		Time:       time.Now(),
		CreatedAt:  time.Now(),
	}

	if err := s.activityRepo.Create(ctx, activity); err != nil {
		log.Error().Err(err).Str("action", action).Str("resource", resource).Str("resourceID", resourceID.String()).Msg("failed to create activity log")
		return
	}

	log.Debug().Str("id", activity.ID.String()).Str("action", action).Str("resource", resource).Msg("activity logged")
}
