package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"git.h2hsecure.com/core/ws/internal/core/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type InventoryRepository struct {
	db *DB
}

func NewInventoryRepository(db *DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) Create(ctx context.Context, inv *domain.Inventory) error {
	query := `
		INSERT INTO inventory (id, realm_id, parent_id, name, description, item_type, embedding, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	now := time.Now()
	inv.CreatedAt = now
	inv.UpdatedAt = now

	embedding := pq.Array(inv.Embedding)
	_, err := r.db.ExecContext(ctx, query,
		inv.ID, inv.RealmID, inv.ParentID, inv.Name, inv.Description, inv.ItemType,
		embedding, inv.CreatedAt, inv.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}

	log.Debug().Str("id", inv.ID.String()).Msg("inventory created")
	return nil
}

func (r *InventoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Inventory, error) {
	query := `
		SELECT id, realm_id, parent_id, name, description, item_type, embedding, created_at, updated_at, deleted_at
		FROM inventory
		WHERE id = $1 AND deleted_at IS NULL
	`

	inv := &domain.Inventory{}
	var embedding []float32
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&inv.ID, &inv.RealmID, &inv.ParentID, &inv.Name, &inv.Description, &inv.ItemType,
		pq.Array(&embedding), &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	inv.Embedding = embedding

	return inv, nil
}

func (r *InventoryRepository) Update(ctx context.Context, inv *domain.Inventory) error {
	inv.UpdatedAt = time.Now()

	if inv.Embedding != nil {
		query := `
			UPDATE inventory SET name = $2, description = $3, embedding = $4, updated_at = $5
			WHERE id = $1 AND deleted_at IS NULL
		`
		embedding := pq.Array(inv.Embedding)
		result, err := r.db.ExecContext(ctx, query, inv.ID, inv.Name, inv.Description, embedding, inv.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to update inventory: %w", err)
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return sql.ErrNoRows
		}
	} else {
		query := `
			UPDATE inventory SET name = $2, description = $3, updated_at = $4
			WHERE id = $1 AND deleted_at IS NULL
		`
		result, err := r.db.ExecContext(ctx, query, inv.ID, inv.Name, inv.Description, inv.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to update inventory: %w", err)
		}
		rows, _ := result.RowsAffected()
		if rows == 0 {
			return sql.ErrNoRows
		}
	}

	log.Debug().Str("id", inv.ID.String()).Msg("inventory updated")
	return nil
}

func (r *InventoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE inventory SET deleted_at = $2 WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, id, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete inventory: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	log.Debug().Str("id", id.String()).Msg("inventory deleted")
	return nil
}

func (r *InventoryRepository) Search(ctx context.Context, opts ...domain.InventorySearchOption) ([]domain.Inventory, error) {
	search := &domain.InventorySearch{
		Limit:  20,
		Offset: 0,
	}
	for _, opt := range opts {
		opt(search)
	}

	hasVectorSearch := search.Vector != nil && len(search.Vector) > 0

	var query string
	if hasVectorSearch {
		query = `
			SELECT id, realm_id, parent_id, name, description, item_type, embedding, created_at, updated_at, deleted_at,
				   1 - (embedding <=> $1) AS similarity
			FROM inventory
			WHERE deleted_at IS NULL
		`
	} else {
		query = `
			SELECT id, realm_id, parent_id, name, description, item_type, embedding, created_at, updated_at, deleted_at
			FROM inventory
			WHERE 1=1
		`
	}

	args := []interface{}{}
	argCount := 0

	if hasVectorSearch {
		argCount++
		vec := pq.Array(search.Vector)
		args = append(args, vec)
	}

	if search.RealmID != nil {
		argCount++
		query += fmt.Sprintf(" AND realm_id = $%d", argCount)
		args = append(args, *search.RealmID)
	}

	if search.ParentID != nil {
		argCount++
		query += fmt.Sprintf(" AND parent_id = $%d", argCount)
		args = append(args, *search.ParentID)
	}

	if search.ItemType != nil {
		argCount++
		query += fmt.Sprintf(" AND item_type = $%d", argCount)
		args = append(args, *search.ItemType)
	}

	if !search.Deleted {
		query += " AND deleted_at IS NULL"
	} else {
		query += " AND deleted_at IS NOT NULL"
	}

	if search.Filter != nil && *search.Filter != "" {
		argCount++
		filterPattern := "%" + strings.ToLower(*search.Filter) + "%"
		query += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d)", argCount, argCount)
		args = append(args, filterPattern)
	}

	if search.StartDate != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *search.StartDate)
	}

	if search.EndDate != nil {
		argCount++
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *search.EndDate)
	}

	if hasVectorSearch {
		argCount++
		query += fmt.Sprintf(" ORDER BY embedding <=> $%d", argCount)
		args = append(args, pq.Array(search.Vector))
	} else {
		query += " ORDER BY created_at DESC"
	}

	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, search.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, search.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search inventory: %w", err)
	}
	defer rows.Close()

	var items []domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		var embedding []float32
		err := rows.Scan(
			&inv.ID, &inv.RealmID, &inv.ParentID, &inv.Name, &inv.Description, &inv.ItemType,
			pq.Array(&embedding), &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}
		inv.Embedding = embedding
		items = append(items, inv)
	}

	return items, rows.Err()
}

func (r *InventoryRepository) SearchSimilar(ctx context.Context, vector domain.Vector, limit int) ([]domain.Inventory, error) {
	query := `
		SELECT id, realm_id, parent_id, name, description, item_type, embedding, created_at, updated_at, deleted_at
		FROM inventory
		WHERE deleted_at IS NULL AND embedding IS NOT NULL
		ORDER BY embedding <=> $1
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, pq.Array(vector), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar inventory: %w", err)
	}
	defer rows.Close()

	var items []domain.Inventory
	for rows.Next() {
		var inv domain.Inventory
		var embedding []float32
		err := rows.Scan(
			&inv.ID, &inv.RealmID, &inv.ParentID, &inv.Name, &inv.Description, &inv.ItemType,
			pq.Array(&embedding), &inv.CreatedAt, &inv.UpdatedAt, &inv.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory: %w", err)
		}
		inv.Embedding = embedding
		items = append(items, inv)
	}

	return items, rows.Err()
}

func (r *InventoryRepository) AddMember(ctx context.Context, member *domain.InventoryMember) error {
	query := `
		INSERT INTO inventory_members (id, inventory_id, user_id, role, membership_time)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (inventory_id, user_id) DO UPDATE SET role = $4, membership_time = $5
	`

	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	if member.MembershipTime.IsZero() {
		member.MembershipTime = time.Now()
	}

	_, err := r.db.ExecContext(ctx, query,
		member.ID, member.InventoryID, member.UserID, member.Role, member.MembershipTime,
	)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}

func (r *InventoryRepository) RemoveMembers(ctx context.Context, inventoryID uuid.UUID, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return nil
	}

	query := `DELETE FROM inventory_members WHERE inventory_id = $1 AND user_id = ANY($2)`
	_, err := r.db.ExecContext(ctx, query, inventoryID, pq.Array(userIDs))
	if err != nil {
		return fmt.Errorf("failed to remove members: %w", err)
	}

	return nil
}

func (r *InventoryRepository) GetMembers(ctx context.Context, inventoryID uuid.UUID, limit, offset int) ([]domain.InventoryMember, error) {
	query := `
		SELECT id, inventory_id, user_id, role, membership_time
		FROM inventory_members
		WHERE inventory_id = $1
		ORDER BY membership_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, inventoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	defer rows.Close()

	var members []domain.InventoryMember
	for rows.Next() {
		var m domain.InventoryMember
		err := rows.Scan(&m.ID, &m.InventoryID, &m.UserID, &m.Role, &m.MembershipTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, m)
	}

	return members, rows.Err()
}
