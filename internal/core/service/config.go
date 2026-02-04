package service

import (
	"context"

	"git.h2hsecure.com/core/ws/internal/core/domain"
)

// GetItemContext implements ports.Service.
func (i *Impl) GetItemContext(ctx context.Context, context domain.ItemContext) ([]domain.ConfigItem, error) {
	panic("unimplemented")
}

// UpsertItemContext implements ports.Service.
func (i *Impl) UpsertItemContext(ctx context.Context, context domain.ItemContext, items []domain.ConfigItem) error {
	panic("unimplemented")
}
