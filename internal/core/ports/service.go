package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

type UserOperations interface {
	SearchUsers(ctx context.Context, limit, offset int) ([]domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

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

type AlarmOperations interface {
	CreateAlarm(ctx context.Context, alarm *domain.Alarm) error
	GetAlarm(ctx context.Context, id uuid.UUID) (*domain.Alarm, error)
	DeleteAlarm(ctx context.Context, id uuid.UUID) error
	SearchAlarms(ctx context.Context, userID uuid.UUID, limit, offset int) ([]domain.Alarm, error)
	AcknowledgeAlarm(ctx context.Context, id uuid.UUID) error
}

type ActivityOperations interface {
	CreateActivity(ctx context.Context, activity *domain.Activity) error
	SearchActivities(ctx context.Context, opts ...domain.ActivitySearchOption) ([]domain.Activity, error)
	GetActivitiesByResourceID(ctx context.Context, resourceID uuid.UUID, limit int) ([]domain.Activity, error)
}

type PasteOperations interface {
	CreatePaste(ctx context.Context, paste *domain.Paste) error
	GetPaste(ctx context.Context, id uuid.UUID) (*domain.Paste, error)
	DeletePaste(ctx context.Context, id uuid.UUID) error
	SearchPastes(ctx context.Context, opts ...domain.PasteSearchOption) ([]domain.Paste, error)
}

type Service interface {
	UserOperations
	InventoryOperations
	AlarmOperations
	ActivityOperations
	PasteOperations

	ConfigOperations
	SnippetOperations
}
