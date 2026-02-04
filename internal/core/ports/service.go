package ports

import (
	"context"

	"git.h2hsecure.com/core/ws/internal/core/domain"
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
}

type GroupOperations interface {
}

type AlarmOperations interface {
}

type Service interface {
	InventoryOperations
	GroupOperations
	AlarmOperations

	ConfigOperations
	SnippetOperations
}
