package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/database"
)

type UserGroupService struct {
	userGroupRepo *database.UserGroupRepository
	activityRepo  *database.ActivityRepository
}

func NewUserGroupService(userGroupRepo *database.UserGroupRepository, activityRepo *database.ActivityRepository) *UserGroupService {
	return &UserGroupService{
		userGroupRepo: userGroupRepo,
		activityRepo:  activityRepo,
	}
}

func (s *UserGroupService) CreateUserGroup(ctx context.Context, ug *domain.UserGroup, userID, realmID uuid.UUID) error {
	err := s.userGroupRepo.Create(ctx, ug)
	if err != nil {
		return err
	}

	s.logActivity(ctx, userID, realmID, "create", "user_group", ug.ID, fmt.Sprintf("Created user group: %s", ug.Name))
	return nil
}

func (s *UserGroupService) GetUserGroup(ctx context.Context, id uuid.UUID) (*domain.UserGroup, error) {
	return s.userGroupRepo.GetByID(ctx, id)
}

func (s *UserGroupService) UpdateUserGroup(ctx context.Context, ug *domain.UserGroup, userID, realmID uuid.UUID) error {
	err := s.userGroupRepo.Update(ctx, ug)
	if err != nil {
		return err
	}

	s.logActivity(ctx, userID, realmID, "update", "user_group", ug.ID, fmt.Sprintf("Updated user group: %s", ug.Name))
	return nil
}

func (s *UserGroupService) DeleteUserGroup(ctx context.Context, id uuid.UUID, userID, realmID uuid.UUID) error {
	ug, err := s.userGroupRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ug == nil {
		return fmt.Errorf("user group not found")
	}

	err = s.userGroupRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.logActivity(ctx, userID, realmID, "delete", "user_group", id, fmt.Sprintf("Deleted user group: %s", ug.Name))
	return nil
}

func (s *UserGroupService) SearchUserGroups(ctx context.Context, limit, offset int) ([]domain.UserGroup, error) {
	return s.userGroupRepo.Search(ctx, limit, offset)
}

func (s *UserGroupService) AddUserGroupMember(ctx context.Context, member *domain.UserGroupMember, userID, realmID uuid.UUID) error {
	err := s.userGroupRepo.AddMember(ctx, member)
	if err != nil {
		return err
	}

	details, _ := json.Marshal(map[string]any{
		"userGroupId": member.UserGroupID,
		"userId":      member.UserID,
		"role":        member.Role,
	})
	s.logActivity(ctx, userID, realmID, "add_member", "user_group", member.UserGroupID, string(details))
	return nil
}

func (s *UserGroupService) RemoveUserGroupMembers(ctx context.Context, userGroupID uuid.UUID, userIDs []uuid.UUID, userID, realmID uuid.UUID) error {
	err := s.userGroupRepo.RemoveMembers(ctx, userGroupID, userIDs)
	if err != nil {
		return err
	}

	details, _ := json.Marshal(map[string]any{
		"userGroupId":  userGroupID,
		"removedUsers": userIDs,
	})
	s.logActivity(ctx, userID, realmID, "remove_members", "user_group", userGroupID, string(details))
	return nil
}

func (s *UserGroupService) GetUserGroupMembers(ctx context.Context, userGroupID uuid.UUID) ([]domain.UserGroupMember, error) {
	return s.userGroupRepo.GetMembers(ctx, userGroupID)
}

func (s *UserGroupService) logActivity(ctx context.Context, userID, realmID uuid.UUID, action, resource string, resourceID uuid.UUID, details string) {
	if s.activityRepo == nil {
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

	_ = s.activityRepo.Create(ctx, activity)
}
