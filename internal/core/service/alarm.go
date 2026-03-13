package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

func (s *Impl) CreateAlarm(ctx context.Context, alarm *domain.Alarm) error {
	return nil
}

func (s *Impl) GetAlarm(ctx context.Context, id uuid.UUID) (*domain.Alarm, error) {
	return nil, nil
}

func (s *Impl) DeleteAlarm(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (s *Impl) SearchAlarms(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Alarm, error) {
	return nil, nil
}

func (s *Impl) AcknowledgeAlarm(ctx context.Context, id uuid.UUID) error {
	return nil
}
