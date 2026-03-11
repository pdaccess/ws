package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

type ConfigOperations interface {
	UpsertItemContext(ctx context.Context, context domain.ItemContext, items []domain.ConfigItem) error
	GetItemContext(ctx context.Context, context domain.ItemContext) ([]domain.ConfigItem, error)
}

type SnippetOperations interface {
	UserSnippets(ctx context.Context, options ...domain.SnippetSearchOption) ([]domain.Snippet, error)
	CreateSnippet(ctx context.Context, snippet domain.Snippet) error
}

type InventoryOperations interface {
	CreateInventory(ctx context.Context, inv *domain.Inventory) error
	GetInventory(ctx context.Context, id uuid.UUID) (*domain.Inventory, error)
	UpdateInventory(ctx context.Context, inv *domain.Inventory) error
	DeleteInventory(ctx context.Context, id uuid.UUID) error
	SearchInventory(ctx context.Context, opts ...domain.InventorySearchOption) ([]domain.Inventory, error)

	AddInventoryMember(ctx context.Context, member *domain.InventoryMember) error
	RemoveInventoryMembers(ctx context.Context, inventoryID uuid.UUID, userIDs []uuid.UUID) error
	GetInventoryMembers(ctx context.Context, inventoryID uuid.UUID, limit, offset int) ([]domain.InventoryMember, error)
}

type GroupOperations any

type AlarmOperations any

type Service interface {
	InventoryOperations
	GroupOperations
	AlarmOperations

	ConfigOperations
	SnippetOperations
}
