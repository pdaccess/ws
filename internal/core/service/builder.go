package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/pdaccess/ws/internal/core/ports"
	"github.com/pdaccess/ws/internal/database"
)

type Impl struct {
	inventoryRepo       *database.InventoryRepository
	userRepo            *database.UserRepository
	activityRepo        *database.ActivityRepository
	pasteRepo           *database.PasteRepository
	userGroupRepo       *database.UserGroupRepository
	serviceSettingsRepo *database.ServiceSettingsRepository
	vector              ports.VectorGenerator
}

func New(inventoryRepo *database.InventoryRepository, userRepo *database.UserRepository, activityRepo *database.ActivityRepository, pasteRepo *database.PasteRepository, userGroupRepo *database.UserGroupRepository, serviceSettingsRepo *database.ServiceSettingsRepository, vector ports.VectorGenerator) ports.Service {
	return &Impl{
		inventoryRepo:       inventoryRepo,
		userRepo:            userRepo,
		activityRepo:        activityRepo,
		pasteRepo:           pasteRepo,
		userGroupRepo:       userGroupRepo,
		serviceSettingsRepo: serviceSettingsRepo,
		vector:              vector,
	}
}

func (s *Impl) SearchUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {
	return nil, nil
}

func (s *Impl) CreateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.Create(ctx, user)
}

func (s *Impl) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *Impl) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *Impl) CreateActivity(ctx context.Context, activity *domain.Activity) error {
	if s.activityRepo == nil {
		return nil
	}
	return s.activityRepo.Create(ctx, activity)
}

func (s *Impl) SearchActivities(ctx context.Context, opts ...domain.ActivitySearchOption) ([]domain.Activity, error) {
	if s.activityRepo == nil {
		return nil, nil
	}
	return s.activityRepo.Search(ctx, opts...)
}

func (s *Impl) GetActivitiesByResourceID(ctx context.Context, resourceID uuid.UUID, limit int) ([]domain.Activity, error) {
	if s.activityRepo == nil {
		return nil, nil
	}
	return s.activityRepo.GetByResourceID(ctx, resourceID, limit)
}

func (s *Impl) CreatePaste(ctx context.Context, paste *domain.Paste) error {
	if s.pasteRepo == nil {
		return nil
	}
	return s.pasteRepo.Create(ctx, paste)
}

func (s *Impl) GetPaste(ctx context.Context, id uuid.UUID) (*domain.Paste, error) {
	if s.pasteRepo == nil {
		return nil, nil
	}
	return s.pasteRepo.GetByID(ctx, id)
}

func (s *Impl) DeletePaste(ctx context.Context, id uuid.UUID) error {
	if s.pasteRepo == nil {
		return nil
	}
	return s.pasteRepo.Delete(ctx, id)
}

func (s *Impl) SearchPastes(ctx context.Context, opts ...domain.PasteSearchOption) ([]domain.Paste, error) {
	if s.pasteRepo == nil {
		return nil, nil
	}
	return s.pasteRepo.Search(ctx, opts...)
}

func (s *Impl) CreateUserGroup(ctx context.Context, ug *domain.UserGroup, userID, realmID uuid.UUID) error {
	if s.userGroupRepo == nil {
		return nil
	}
	return s.userGroupRepo.Create(ctx, ug)
}

func (s *Impl) GetUserGroup(ctx context.Context, id uuid.UUID) (*domain.UserGroup, error) {
	if s.userGroupRepo == nil {
		return nil, nil
	}
	return s.userGroupRepo.GetByID(ctx, id)
}

func (s *Impl) UpdateUserGroup(ctx context.Context, ug *domain.UserGroup, userID, realmID uuid.UUID) error {
	if s.userGroupRepo == nil {
		return nil
	}
	return s.userGroupRepo.Update(ctx, ug)
}

func (s *Impl) DeleteUserGroup(ctx context.Context, id uuid.UUID, userID, realmID uuid.UUID) error {
	if s.userGroupRepo == nil {
		return nil
	}
	return s.userGroupRepo.Delete(ctx, id)
}

func (s *Impl) SearchUserGroups(ctx context.Context, limit, offset int) ([]domain.UserGroup, error) {
	if s.userGroupRepo == nil {
		return nil, nil
	}
	return s.userGroupRepo.Search(ctx, limit, offset)
}

func (s *Impl) AddUserGroupMember(ctx context.Context, member *domain.UserGroupMember, userID, realmID uuid.UUID) error {
	if s.userGroupRepo == nil {
		return nil
	}
	return s.userGroupRepo.AddMember(ctx, member)
}

func (s *Impl) RemoveUserGroupMembers(ctx context.Context, userGroupID uuid.UUID, userIDs []uuid.UUID, userID, realmID uuid.UUID) error {
	if s.userGroupRepo == nil {
		return nil
	}
	return s.userGroupRepo.RemoveMembers(ctx, userGroupID, userIDs)
}

func (s *Impl) GetUserGroupMembers(ctx context.Context, userGroupID uuid.UUID) ([]domain.UserGroupMember, error) {
	if s.userGroupRepo == nil {
		return nil, nil
	}
	return s.userGroupRepo.GetMembers(ctx, userGroupID)
}
