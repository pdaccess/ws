package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/rs/zerolog/log"
)

type UserGroupRepository struct {
	db *DB
}

func NewUserGroupRepository(db *DB) *UserGroupRepository {
	return &UserGroupRepository{db: db}
}

func (r *UserGroupRepository) Create(ctx context.Context, ug *domain.UserGroup) error {
	query := `
		INSERT INTO user_groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	if ug.ID == uuid.Nil {
		ug.ID = uuid.New()
	}
	now := time.Now()
	ug.CreatedAt = now
	ug.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		ug.ID, ug.Name, ug.Description, ug.CreatedAt, ug.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create user group: %w", err)
	}

	log.Debug().Str("id", ug.ID.String()).Msg("user group created")
	return nil
}

func (r *UserGroupRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.UserGroup, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM user_groups
		WHERE id = $1 AND deleted_at IS NULL
	`

	ug := &domain.UserGroup{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ug.ID, &ug.Name, &ug.Description, &ug.CreatedAt, &ug.UpdatedAt, &ug.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user group: %w", err)
	}

	return ug, nil
}

func (r *UserGroupRepository) Update(ctx context.Context, ug *domain.UserGroup) error {
	query := `
		UPDATE user_groups SET name = $2, description = $3, updated_at = $4
		WHERE id = $1 AND deleted_at IS NULL
	`

	ug.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query, ug.ID, ug.Name, ug.Description, ug.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user group: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Debug().Str("id", ug.ID.String()).Msg("user group updated")
	return nil
}

func (r *UserGroupRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE user_groups SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete user group: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Debug().Str("id", id.String()).Msg("user group deleted")
	return nil
}

func (r *UserGroupRepository) Search(ctx context.Context, limit, offset int) ([]domain.UserGroup, error) {
	query := `
		SELECT id, name, description, created_at, updated_at, deleted_at
		FROM user_groups
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search user groups: %w", err)
	}
	defer rows.Close()

	var groups []domain.UserGroup
	for rows.Next() {
		var ug domain.UserGroup
		err := rows.Scan(&ug.ID, &ug.Name, &ug.Description, &ug.CreatedAt, &ug.UpdatedAt, &ug.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user group: %w", err)
		}
		groups = append(groups, ug)
	}

	return groups, rows.Err()
}

func (r *UserGroupRepository) AddMember(ctx context.Context, member *domain.UserGroupMember) error {
	query := `
		INSERT INTO user_group_members (id, user_group_id, user_id, role, membership_time)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_group_id, user_id) DO UPDATE SET role = $4, membership_time = $5
	`

	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	if member.MembershipTime.IsZero() {
		member.MembershipTime = time.Now()
	}

	_, err := r.db.ExecContext(ctx, query,
		member.ID, member.UserGroupID, member.UserID, member.Role, member.MembershipTime,
	)
	if err != nil {
		return fmt.Errorf("failed to add user group member: %w", err)
	}

	return nil
}

func (r *UserGroupRepository) RemoveMembers(ctx context.Context, userGroupID uuid.UUID, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	query := `DELETE FROM user_group_members WHERE user_group_id = $1 AND user_id = ANY($2)`
	_, err := r.db.ExecContext(ctx, query, userGroupID, pq.Array(userIDs))
	if err != nil {
		return fmt.Errorf("failed to remove user group members: %w", err)
	}

	return nil
}

func (r *UserGroupRepository) GetMembers(ctx context.Context, userGroupID uuid.UUID) ([]domain.UserGroupMember, error) {
	query := `
		SELECT id, user_group_id, user_id, role, membership_time
		FROM user_group_members
		WHERE user_group_id = $1
		ORDER BY membership_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userGroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user group members: %w", err)
	}
	defer rows.Close()

	var members []domain.UserGroupMember
	for rows.Next() {
		var m domain.UserGroupMember
		err := rows.Scan(&m.ID, &m.UserGroupID, &m.UserID, &m.Role, &m.MembershipTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user group member: %w", err)
		}
		members = append(members, m)
	}

	return members, rows.Err()
}
