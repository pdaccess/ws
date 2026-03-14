package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/rs/zerolog/log"
)

type PasteRepository struct {
	db *DB
}

func NewPasteRepository(db *DB) *PasteRepository {
	return &PasteRepository{db: db}
}

func (r *PasteRepository) Create(ctx context.Context, paste *domain.Paste) error {
	query := `
		INSERT INTO pastes (id, user_id, content, language, expires_at, views, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	if paste.ID == uuid.Nil {
		paste.ID = uuid.New()
	}
	if paste.CreatedAt.IsZero() {
		paste.CreatedAt = time.Now()
	}
	if paste.Views == 0 {
		paste.Views = 0
	}

	_, err := r.db.ExecContext(ctx, query,
		paste.ID, paste.UserID, paste.Content, paste.Language,
		paste.ExpiresAt, paste.Views, paste.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create paste: %w", err)
	}

	log.Debug().Str("id", paste.ID.String()).Msg("paste created")
	return nil
}

func (r *PasteRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Paste, error) {
	query := `
		SELECT id, user_id, content, language, expires_at, views, created_at
		FROM pastes
		WHERE id = $1
	`

	paste := &domain.Paste{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&paste.ID, &paste.UserID, &paste.Content, &paste.Language,
		&paste.ExpiresAt, &paste.Views, &paste.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get paste: %w", err)
	}

	return paste, nil
}

func (r *PasteRepository) Search(ctx context.Context, opts ...domain.PasteSearchOption) ([]domain.Paste, error) {
	search := &domain.PasteSearch{
		Limit:  20,
		Offset: 0,
	}
	for _, opt := range opts {
		opt(search)
	}

	query := `
		SELECT id, user_id, content, language, expires_at, views, created_at
		FROM pastes
		WHERE (expires_at IS NULL OR expires_at > NOW())
	`
	args := []any{}
	argNum := 1

	if search.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, *search.UserID)
		argNum++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, search.Limit, search.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search pastes: %w", err)
	}
	defer rows.Close()

	var pastes []domain.Paste
	for rows.Next() {
		var p domain.Paste
		err := rows.Scan(
			&p.ID, &p.UserID, &p.Content, &p.Language,
			&p.ExpiresAt, &p.Views, &p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan paste: %w", err)
		}
		pastes = append(pastes, p)
	}

	return pastes, nil
}

func (r *PasteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM pastes WHERE id = $1`

	rs, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete paste: %w", err)
	}

	count, err := rs.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if count == 0 {
		return domain.ErrNotFound
	}

	log.Debug().Str("id", id.String()).Msg("paste deleted")
	return nil
}

func (r *PasteRepository) IncrementViews(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE pastes SET views = views + 1 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment paste views: %w", err)
	}

	return nil
}
