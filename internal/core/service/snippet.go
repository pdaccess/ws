package service

import (
	"context"

	"git.h2hsecure.com/core/ws/internal/core/domain"
)

// CreateSnippet implements ports.Service.
func (i *Impl) CreateSnippet(ctx context.Context, snippet domain.Snippet) error {
	panic("unimplemented")
}

// UserSnippets implements ports.Service.
func (i *Impl) UserSnippets(ctx context.Context, options ...domain.SnippetSearchOption) ([]domain.Snippet, error) {
	panic("unimplemented")
}
