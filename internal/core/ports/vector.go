package ports

import (
	"context"

	"github.com/pdaccess/ws/internal/core/domain"
)

type VectorGenerator interface {
	Generate(ctx context.Context, queryTerm string) (domain.Vector, error)
	Close()
}
