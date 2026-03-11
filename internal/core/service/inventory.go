package service

import (
	"context"

	"git.h2hsecure.com/core/ws/internal/core/domain"
	"github.com/google/uuid"
)

func (s *Impl) CreateInventory(ctx context.Context, inv *domain.Inventory) error {
	return s.inventoryRepo.Create(ctx, inv)
}

func (s *Impl) GetInventory(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	return s.inventoryRepo.GetByID(ctx, id)
}

func (s *Impl) UpdateInventory(ctx context.Context, inv *domain.Inventory) error {
	return s.inventoryRepo.Update(ctx, inv)
}

func (s *Impl) DeleteInventory(ctx context.Context, id uuid.UUID) error {
	return s.inventoryRepo.Delete(ctx, id)
}

func (s *Impl) SearchInventory(ctx context.Context, opts ...domain.InventorySearchOption) ([]domain.Inventory, error) {
	return s.inventoryRepo.Search(ctx, opts...)
}

func (s *Impl) AddInventoryMember(ctx context.Context, member *domain.InventoryMember) error {
	return s.inventoryRepo.AddMember(ctx, member)
}

func (s *Impl) RemoveInventoryMembers(ctx context.Context, inventoryID uuid.UUID, userIDs []uuid.UUID) error {
	return s.inventoryRepo.RemoveMembers(ctx, inventoryID, userIDs)
}

func (s *Impl) GetInventoryMembers(ctx context.Context, inventoryID uuid.UUID, limit, offset int) ([]domain.InventoryMember, error) {
	return s.inventoryRepo.GetMembers(ctx, inventoryID, limit, offset)
}
