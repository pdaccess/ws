package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
)

type AuditRepository struct {
	db *DB
}

func NewAuditRepository(db *DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(ctx context.Context, entry *domain.AuditEntry) error {
	query := `INSERT INTO ws_audit_logs (id, timestamp, actor_id, action, resource_id, success) VALUES ($1, $2, $3, $4, $5, $6)`
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	_, err := r.db.ExecContext(ctx, query, entry.ID, entry.Timestamp, entry.ActorID, entry.Action, entry.ResourceID, entry.Success)
	if err != nil {
		return fmt.Errorf("failed to create audit: %w", err)
	}
	return nil
}

func (r *AuditRepository) Search(ctx context.Context, actorID, resourceID *uuid.UUID, from *string) ([]domain.AuditEntry, error) {
	query := `SELECT id, timestamp, actor_id, action, resource_id, success FROM ws_audit_logs WHERE 1=1`
	args := []any{}
	argNum := 1

	if actorID != nil {
		query += fmt.Sprintf(" AND actor_id = $%d", argNum)
		args = append(args, *actorID)
		argNum++
	}
	if resourceID != nil {
		query += fmt.Sprintf(" AND resource_id = $%d", argNum)
		args = append(args, *resourceID)
		argNum++
	}
	if from != nil && *from != "" {
		query += fmt.Sprintf(" AND timestamp >= $%d", argNum)
		args = append(args, *from)
		argNum++
	}

	query += " ORDER BY timestamp DESC LIMIT 100"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search audits: %w", err)
	}
	defer rows.Close()

	var entries []domain.AuditEntry
	for rows.Next() {
		var e domain.AuditEntry
		if err := rows.Scan(&e.ID, &e.Timestamp, &e.ActorID, &e.Action, &e.ResourceID, &e.Success); err != nil {
			return nil, fmt.Errorf("failed to scan audit: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}
