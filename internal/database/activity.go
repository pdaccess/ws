package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/rs/zerolog/log"
)

type ActivityRepository struct {
	db *DB
}

func NewActivityRepository(db *DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) Create(ctx context.Context, activity *domain.Activity) error {
	query := `
		INSERT INTO activities (id, user_id, realm_id, action, resource, resource_id, details, ip_address, activity_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	if activity.ID == uuid.Nil {
		activity.ID = uuid.New()
	}
	if activity.CreatedAt.IsZero() {
		activity.CreatedAt = time.Now()
	}
	if activity.Time.IsZero() {
		activity.Time = activity.CreatedAt
	}

	_, err := r.db.ExecContext(ctx, query,
		activity.ID, activity.UserID, activity.RealmID, activity.Action,
		activity.Resource, activity.ResourceID, activity.Details,
		activity.IPAddress, activity.Time, activity.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	log.Debug().Str("id", activity.ID.String()).Msg("activity created")
	return nil
}

func (r *ActivityRepository) Search(ctx context.Context, opts ...domain.ActivitySearchOption) ([]domain.Activity, error) {
	search := &domain.ActivitySearch{
		Limit:  20,
		Offset: 0,
	}
	for _, opt := range opts {
		opt(search)
	}

	query := `
		SELECT id, user_id, realm_id, action, resource, resource_id, details, ip_address, activity_time, created_at
		FROM activities
		WHERE 1=1
	`
	args := []any{}
	argNum := 1

	if search.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argNum)
		args = append(args, *search.UserID)
		argNum++
	}
	if search.RealmID != nil {
		query += fmt.Sprintf(" AND realm_id = $%d", argNum)
		args = append(args, *search.RealmID)
		argNum++
	}
	if search.Action != nil {
		query += fmt.Sprintf(" AND action = $%d", argNum)
		args = append(args, *search.Action)
		argNum++
	}
	if search.StartDate != nil {
		query += fmt.Sprintf(" AND activity_time >= $%d", argNum)
		args = append(args, *search.StartDate)
		argNum++
	}
	if search.EndDate != nil {
		query += fmt.Sprintf(" AND activity_time <= $%d", argNum)
		args = append(args, *search.EndDate)
		argNum++
	}

	query += fmt.Sprintf(" ORDER BY activity_time DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, search.Limit, search.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search activities: %w", err)
	}
	defer rows.Close()

	var activities []domain.Activity
	for rows.Next() {
		var a domain.Activity
		err := rows.Scan(
			&a.ID, &a.UserID, &a.RealmID, &a.Action, &a.Resource,
			&a.ResourceID, &a.Details, &a.IPAddress, &a.Time, &a.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, a)
	}

	return activities, nil
}

func (r *ActivityRepository) GetByResourceID(ctx context.Context, resourceID uuid.UUID, limit int) ([]domain.Activity, error) {
	query := `
		SELECT id, user_id, realm_id, action, resource, resource_id, details, ip_address, activity_time, created_at
		FROM activities
		WHERE resource_id = $1
		ORDER BY activity_time DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, resourceID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get activities by resource: %w", err)
	}
	defer rows.Close()

	var activities []domain.Activity
	for rows.Next() {
		var a domain.Activity
		err := rows.Scan(
			&a.ID, &a.UserID, &a.RealmID, &a.Action, &a.Resource,
			&a.ResourceID, &a.Details, &a.IPAddress, &a.Time, &a.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}
		activities = append(activities, a)
	}

	return activities, nil
}
