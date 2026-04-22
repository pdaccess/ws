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
	activityRepo        *database.ActivityRepository
	pasteRepo           *database.PasteRepository
	serviceSettingsRepo *database.ServiceSettingsRepository
	credentialRepo      *database.CredentialRepository
	vector              ports.VectorGenerator
}

func New(inventoryRepo *database.InventoryRepository, activityRepo *database.ActivityRepository, pasteRepo *database.PasteRepository, serviceSettingsRepo *database.ServiceSettingsRepository, credentialRepo *database.CredentialRepository, vector ports.VectorGenerator) ports.Service {
	return &Impl{
		inventoryRepo:       inventoryRepo,
		activityRepo:        activityRepo,
		pasteRepo:           pasteRepo,
		serviceSettingsRepo: serviceSettingsRepo,
		credentialRepo:      credentialRepo,
		vector:              vector,
	}
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

func (s *Impl) CreateCredential(ctx context.Context, cred *domain.Credential) error {
	if s.credentialRepo == nil {
		return nil
	}
	return s.credentialRepo.Create(ctx, cred)
}

func (s *Impl) GetCredential(ctx context.Context, id uuid.UUID) (*domain.Credential, error) {
	if s.credentialRepo == nil {
		return nil, nil
	}
	return s.credentialRepo.GetByID(ctx, id)
}

func (s *Impl) UpdateCredential(ctx context.Context, cred *domain.Credential) error {
	if s.credentialRepo == nil {
		return nil
	}
	return s.credentialRepo.Update(ctx, cred)
}

func (s *Impl) DeleteCredential(ctx context.Context, id uuid.UUID) error {
	if s.credentialRepo == nil {
		return nil
	}
	return s.credentialRepo.Delete(ctx, id)
}

func (s *Impl) SearchCredentials(ctx context.Context, opts ...domain.CredentialSearchOption) ([]domain.Credential, error) {
	if s.credentialRepo == nil {
		return nil, nil
	}
	return s.credentialRepo.Search(ctx, opts...)
}

func (s *Impl) CreateCredentialSecret(ctx context.Context, secret *domain.CredentialSecret) error {
	if s.credentialRepo == nil {
		return nil
	}
	return s.credentialRepo.CreateSecret(ctx, secret)
}

func (s *Impl) GetCredentialSecret(ctx context.Context, credentialID uuid.UUID) (*domain.CredentialSecret, error) {
	if s.credentialRepo == nil {
		return nil, nil
	}
	return s.credentialRepo.GetSecretByCredentialID(ctx, credentialID)
}

func (s *Impl) UpdateCredentialSecret(ctx context.Context, secret *domain.CredentialSecret) error {
	if s.credentialRepo == nil {
		return nil
	}
	return s.credentialRepo.UpdateSecret(ctx, secret)
}
