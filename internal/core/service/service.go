package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

func (s *Impl) CreateService(ctx context.Context, svc *domain.Service, userID, realmID uuid.UUID) error {
	var err error

	if s.vector != nil {
		svc.Embedding, err = s.vector.Generate(ctx, fmt.Sprintf("%s %s", svc.Name, svc.Description))
		if err != nil {
			return fmt.Errorf("vector generation: %w", err)
		}
	}

	if err := s.inventoryRepo.CreateService(ctx, svc); err != nil {
		return err
	}

	if svc.Settings != nil {
		svc.Settings.ServiceID = svc.ID
		err = s.serviceSettingsRepo.Upsert(ctx, svc.Settings)
		if err != nil {
			return err
		}
	}

	s.logActivity(ctx, userID, realmID, "create", "service", svc.ID, fmt.Sprintf("Created service: %s", svc.Name))
	return nil
}

func (s *Impl) GetService(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	svc, err := s.inventoryRepo.GetServiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}

	settings, err := s.serviceSettingsRepo.GetByInventoryID(ctx, id)
	if err != nil {
		return nil, err
	}
	if settings != nil {
		svc.Settings = settings
	}

	return svc, nil
}

func (s *Impl) UpdateService(ctx context.Context, svc *domain.Service, userID, realmID uuid.UUID) error {
	err := s.inventoryRepo.UpdateService(ctx, svc)
	if err != nil {
		return err
	}

	if svc.Settings != nil {
		svc.Settings.ServiceID = svc.ID
		err = s.serviceSettingsRepo.Upsert(ctx, svc.Settings)
		if err != nil {
			return err
		}
	}

	s.logActivity(ctx, userID, realmID, "update", "service", svc.ID, fmt.Sprintf("Updated service: %s", svc.Name))
	return nil
}

func (s *Impl) DeleteService(ctx context.Context, id uuid.UUID, userID, realmID uuid.UUID) error {
	svc, err := s.inventoryRepo.GetServiceByID(ctx, id)
	if err != nil {
		return err
	}
	if svc == nil {
		return fmt.Errorf("service not found")
	}

	err = s.serviceSettingsRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	err = s.inventoryRepo.DeleteService(ctx, id)
	if err != nil {
		return err
	}

	s.logActivity(ctx, userID, realmID, "delete", "service", id, fmt.Sprintf("Deleted service: %s", svc.Name))
	return nil
}

func (s *Impl) SearchServices(ctx context.Context, opts ...domain.ServiceSearchOption) ([]domain.Service, error) {
	return s.inventoryRepo.SearchServices(ctx, opts...)
}

func (s *Impl) AddServiceMember(ctx context.Context, member *domain.ServiceMember, userID, realmID uuid.UUID) error {
	err := s.inventoryRepo.AddServiceMember(ctx, member)
	if err != nil {
		return err
	}

	details, _ := json.Marshal(map[string]any{
		"serviceId": member.ServiceID,
		"userId":    member.UserID,
		"role":      member.Role,
	})
	s.logActivity(ctx, userID, realmID, "add_member", "service", member.ServiceID, string(details))
	return nil
}

func (s *Impl) RemoveServiceMembers(ctx context.Context, serviceID uuid.UUID, userIDs []uuid.UUID, userID, realmID uuid.UUID) error {
	err := s.inventoryRepo.RemoveServiceMembers(ctx, serviceID, userIDs)
	if err != nil {
		return err
	}

	details, _ := json.Marshal(map[string]any{
		"serviceId":    serviceID,
		"removedUsers": userIDs,
	})
	s.logActivity(ctx, userID, realmID, "remove_members", "service", serviceID, string(details))
	return nil
}

func (s *Impl) GetServiceMembers(ctx context.Context, serviceID uuid.UUID, limit, offset int) ([]domain.ServiceMember, error) {
	return s.inventoryRepo.GetServiceMembers(ctx, serviceID, limit, offset)
}

func (s *Impl) UpsertServiceSettings(ctx context.Context, settings *domain.ServiceSettings, userID, realmID uuid.UUID) error {
	return s.serviceSettingsRepo.Upsert(ctx, settings)
}

func (s *Impl) GetServiceSettings(ctx context.Context, serviceID uuid.UUID) (*domain.ServiceSettings, error) {
	return s.serviceSettingsRepo.GetByInventoryID(ctx, serviceID)
}

func (s *Impl) SearchServicesVector(ctx context.Context, vector domain.Vector, limit, offset int) ([]domain.Service, error) {
	return s.inventoryRepo.SearchServices(ctx,
		domain.WithServiceVector(vector),
		domain.WithServiceLimit(limit),
		domain.WithServiceOffset(offset),
	)
}

func (s *Impl) SearchServicesWithQuery(ctx context.Context, query string, limit, offset int) ([]domain.Service, error) {
	if s.vector != nil {
		vector, err := s.vector.Generate(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("vector generation: %w", err)
		}
		return s.SearchServicesVector(ctx, vector, limit, offset)
	}

	return s.inventoryRepo.SearchServices(ctx,
		domain.WithServiceFilter(query),
		domain.WithServiceLimit(limit),
		domain.WithServiceOffset(offset),
	)
}
