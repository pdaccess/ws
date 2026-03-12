package tests

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

type mockService struct{}

func NewMockService() *mockService {
	return &mockService{}
}

func (m *mockService) UpsertItemContext(ctx context.Context, context domain.ItemContext, items []domain.ConfigItem) error {
	return nil
}

func (m *mockService) GetItemContext(ctx context.Context, context domain.ItemContext) ([]domain.ConfigItem, error) {
	return nil, nil
}

func (m *mockService) UserSnippets(ctx context.Context, options ...domain.SnippetSearchOption) ([]domain.Snippet, error) {
	return nil, nil
}

func (m *mockService) CreateSnippet(ctx context.Context, snippet domain.Snippet) error {
	return nil
}

func (m *mockService) CreateInventory(ctx context.Context, inv *domain.Inventory) error {
	return nil
}

func (m *mockService) GetInventory(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	return nil, nil
}

func (m *mockService) UpdateInventory(ctx context.Context, inv *domain.Inventory) error {
	return nil
}

func (m *mockService) DeleteInventory(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockService) SearchInventory(ctx context.Context, opts ...domain.InventorySearchOption) ([]domain.Inventory, error) {
	return nil, nil
}

func (m *mockService) AddInventoryMember(ctx context.Context, member *domain.InventoryMember) error {
	return nil
}

func (m *mockService) RemoveInventoryMembers(ctx context.Context, inventoryID uuid.UUID, userIDs []uuid.UUID) error {
	return nil
}

func (m *mockService) GetInventoryMembers(ctx context.Context, inventoryID uuid.UUID, limit, offset int) ([]domain.InventoryMember, error) {
	return nil, nil
}
