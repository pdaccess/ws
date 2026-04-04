package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/rs/zerolog/log"
)

type InventoryRepository struct {
	db *DB
}

func NewInventoryRepository(db *DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func vecToString(v []float64) string {
	if v == nil {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, f := range v {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%f", float32(f)))
	}
	sb.WriteString("]")
	return sb.String()
}

func vec32ToString(v []float32) string {
	if v == nil {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteString("[")
	for i, f := range v {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf("%f", f))
	}
	sb.WriteString("]")
	return sb.String()
}

func parseVector(s string) []float64 {
	if s == "" || s == "[]" {
		return nil
	}
	s = strings.Trim(s, "[]")
	parts := strings.Split(s, ",")
	result := make([]float64, len(parts))
	for i, p := range parts {
		var v float64
		fmt.Sscanf(strings.TrimSpace(p), "%f", &v)
		result[i] = v
	}
	return result
}

func (r *InventoryRepository) CreateGroup(ctx context.Context, group *domain.Group) error {
	query := `
		INSERT INTO inventory (id, realm_id, parent_id, name, description, type, embedding, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if group.ID == uuid.Nil {
		group.ID = uuid.New()
	}
	now := time.Now()
	group.CreatedAt = now
	group.UpdatedAt = now

	if group.Type == "" {
		group.Type = domain.InventoryTypeGroup
	}

	var embedding any
	if group.Embedding != nil {
		embedding = vecToString(group.Embedding)
	} else {
		embedding = nil
	}

	_, err := r.db.ExecContext(ctx, query,
		group.ID, group.RealmID, group.ParentID, group.Name, group.Description,
		group.Type, embedding, group.CreatedAt, group.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create group: %w", err)
	}

	log.Debug().Str("id", group.ID.String()).Msg("group created")
	return nil
}

func (r *InventoryRepository) GetGroupByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	query := `
		SELECT id, realm_id, parent_id, name, description, type, embedding, created_at, updated_at, deleted_at
		FROM inventory
		WHERE id = $1 AND deleted_at IS NULL
	`

	group := &domain.Group{}
	var embedding sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID, &group.RealmID, &group.ParentID, &group.Name, &group.Description,
		&group.Type, &embedding, &group.CreatedAt, &group.UpdatedAt, &group.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}
	if embedding.Valid {
		group.Embedding = parseVector(embedding.String)
	}

	return group, nil
}

func (r *InventoryRepository) UpdateGroup(ctx context.Context, group *domain.Group) error {
	group.UpdatedAt = time.Now()

	var query string
	var args []any

	if group.Embedding != nil {
		query = `
			UPDATE inventory SET name = $2, description = $3, embedding = $4, updated_at = $5
			WHERE id = $1 AND deleted_at IS NULL
		`
		args = []any{group.ID, group.Name, group.Description, vecToString(group.Embedding), group.UpdatedAt}
	} else {
		query = `
			UPDATE inventory SET name = $2, description = $3, updated_at = $4
			WHERE id = $1 AND deleted_at IS NULL
		`
		args = []any{group.ID, group.Name, group.Description, group.UpdatedAt}
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update group: %w", err)
	}

	log.Debug().Str("id", group.ID.String()).Msg("group updated")
	return nil
}

func (r *InventoryRepository) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE inventory SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}

	log.Debug().Str("id", id.String()).Msg("group deleted")
	return nil
}

func (r *InventoryRepository) SearchGroups(ctx context.Context, opts ...domain.GroupSearchOption) ([]domain.Group, error) {
	search := &domain.GroupSearch{
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
			SELECT id, realm_id, parent_id, name, description, type, embedding, created_at, updated_at, deleted_at,
				   1 - (embedding <=> $1::vector) AS similarity
			FROM inventory
			WHERE deleted_at IS NULL
		`
	} else {
		query = `
			SELECT id, realm_id, parent_id, name, description, type, embedding, created_at, updated_at, deleted_at,
				   0 AS similarity
			FROM inventory
			WHERE 1=1
		`
	}

	args := []any{}
	argCount := 0

	if hasVectorSearch {
		argCount++
		args = append(args, vecToString(search.Vector))
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

	if !search.Deleted {
		query += " AND deleted_at IS NULL"
	} else {
		query += " AND deleted_at IS NOT NULL"
	}

	if search.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, string(search.Type))
	} else {
		argCount++
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, string(domain.InventoryTypeGroup))
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
		query += fmt.Sprintf(" ORDER BY embedding <=> $%d::vector", argCount)
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
		return nil, fmt.Errorf("failed to search groups: %w", err)
	}
	defer rows.Close()

	var items []domain.Group
	for rows.Next() {
		var g domain.Group
		var embedding sql.NullString
		var similarity float64
		err := rows.Scan(
			&g.ID, &g.RealmID, &g.ParentID, &g.Name, &g.Description,
			&g.Type, &embedding, &g.CreatedAt, &g.UpdatedAt, &g.DeletedAt, &similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		if embedding.Valid {
			g.Embedding = parseVector(embedding.String)
		}
		items = append(items, g)
	}

	return items, rows.Err()
}

func (r *InventoryRepository) AddGroupMember(ctx context.Context, member *domain.GroupMember) error {
	query := `
		INSERT INTO group_members (id, group_id, user_id, role, membership_time)
		VALUES ($1, $2, $3, $4, $5)
	`

	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	member.MembershipTime = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		member.ID, member.GroupID, member.UserID, member.Role, member.MembershipTime,
	)
	if err != nil {
		return fmt.Errorf("failed to add group member: %w", err)
	}

	log.Debug().Str("id", member.ID.String()).Msg("group member added")
	return nil
}

func (r *InventoryRepository) RemoveGroupMember(ctx context.Context, groupID, userID uuid.UUID) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove group member: %w", err)
	}

	log.Debug().Str("group_id", groupID.String()).Str("user_id", userID.String()).Msg("group member removed")
	return nil
}

func (r *InventoryRepository) GetGroupMembers(ctx context.Context, groupID uuid.UUID, limit, offset int) ([]domain.GroupMember, error) {
	query := `
		SELECT id, group_id, user_id, role, membership_time
		FROM group_members
		WHERE group_id = $1
		ORDER BY membership_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, groupID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get group members: %w", err)
	}
	defer rows.Close()

	var members []domain.GroupMember
	for rows.Next() {
		var m domain.GroupMember
		err := rows.Scan(&m.ID, &m.GroupID, &m.UserID, &m.Role, &m.MembershipTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan group member: %w", err)
		}
		members = append(members, m)
	}

	return members, rows.Err()
}

func (r *InventoryRepository) GetGroupUserIDs(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM group_members WHERE group_id = $1`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get group user ids: %w", err)
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan user id: %w", err)
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, rows.Err()
}

func (r *InventoryRepository) RemoveGroupMembers(ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`

	for _, userID := range userIDs {
		_, err := r.db.ExecContext(ctx, query, groupID, userID)
		if err != nil {
			return fmt.Errorf("failed to remove group member: %w", err)
		}
	}

	log.Debug().Str("group_id", groupID.String()).Msg("group members removed")
	return nil
}

func (r *InventoryRepository) CreateService(ctx context.Context, svc *domain.Service) error {
	query := `
		INSERT INTO inventory (id, realm_id, parent_id, name, description, type, embedding, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if svc.ID == uuid.Nil {
		svc.ID = uuid.New()
	}
	now := time.Now()
	svc.CreatedAt = now
	svc.UpdatedAt = now

	if svc.Type == "" {
		svc.Type = domain.InventoryTypeService
	}

	var embedding any
	if svc.Embedding != nil {
		embedding = vecToString(svc.Embedding)
	} else {
		embedding = nil
	}

	_, err := r.db.ExecContext(ctx, query,
		svc.ID, svc.RealmID, svc.ParentID, svc.Name, svc.Description,
		svc.Type, embedding, svc.CreatedAt, svc.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}

	log.Debug().Str("id", svc.ID.String()).Msg("service created")
	return nil
}

func (r *InventoryRepository) GetServiceByID(ctx context.Context, id uuid.UUID) (*domain.Service, error) {
	query := `
		SELECT id, realm_id, parent_id, name, description, type, embedding, created_at, updated_at, deleted_at
		FROM inventory
		WHERE id = $1 AND deleted_at IS NULL
	`

	svc := &domain.Service{}
	var embedding sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&svc.ID, &svc.RealmID, &svc.ParentID, &svc.Name, &svc.Description,
		&svc.Type, &embedding, &svc.CreatedAt, &svc.UpdatedAt, &svc.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	if embedding.Valid {
		svc.Embedding = parseVector(embedding.String)
	}

	return svc, nil
}

func (r *InventoryRepository) UpdateService(ctx context.Context, svc *domain.Service) error {
	svc.UpdatedAt = time.Now()

	var query string
	var args []any

	if svc.Embedding != nil {
		query = `
			UPDATE inventory SET name = $2, description = $3, embedding = $4, updated_at = $5
			WHERE id = $1 AND deleted_at IS NULL
		`
		args = []any{svc.ID, svc.Name, svc.Description, vecToString(svc.Embedding), svc.UpdatedAt}
	} else {
		query = `
			UPDATE inventory SET name = $2, description = $3, updated_at = $4
			WHERE id = $1 AND deleted_at IS NULL
		`
		args = []any{svc.ID, svc.Name, svc.Description, svc.UpdatedAt}
	}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	log.Debug().Str("id", svc.ID.String()).Msg("service updated")
	return nil
}

func (r *InventoryRepository) DeleteService(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE inventory SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	log.Debug().Str("id", id.String()).Msg("service deleted")
	return nil
}

func (r *InventoryRepository) SearchServices(ctx context.Context, opts ...domain.ServiceSearchOption) ([]domain.Service, error) {
	search := &domain.ServiceSearch{
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
			SELECT id, realm_id, parent_id, name, description, type, embedding, created_at, updated_at, deleted_at,
				   1 - (embedding <=> $1::vector) AS similarity
			FROM inventory
			WHERE deleted_at IS NULL
		`
	} else {
		query = `
			SELECT id, realm_id, parent_id, name, description, type, embedding, created_at, updated_at, deleted_at,
				   0 AS similarity
			FROM inventory
			WHERE 1=1
		`
	}

	args := []any{}
	argCount := 0

	if hasVectorSearch {
		argCount++
		args = append(args, vecToString(search.Vector))
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

	if !search.Deleted {
		query += " AND deleted_at IS NULL"
	} else {
		query += " AND deleted_at IS NOT NULL"
	}

	if search.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, string(search.Type))
	} else {
		argCount++
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, string(domain.InventoryTypeService))
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
		query += fmt.Sprintf(" ORDER BY embedding <=> $%d::vector", argCount)
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
		return nil, fmt.Errorf("failed to search services: %w", err)
	}
	defer rows.Close()

	var items []domain.Service
	for rows.Next() {
		var s domain.Service
		var embedding sql.NullString
		var similarity float64
		err := rows.Scan(
			&s.ID, &s.RealmID, &s.ParentID, &s.Name, &s.Description,
			&s.Type, &embedding, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt, &similarity,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service: %w", err)
		}
		if embedding.Valid {
			s.Embedding = parseVector(embedding.String)
		}
		items = append(items, s)
	}

	return items, rows.Err()
}

func (r *InventoryRepository) AddServiceMember(ctx context.Context, member *domain.ServiceMember) error {
	query := `
		INSERT INTO service_members (id, service_id, user_id, role, membership_time)
		VALUES ($1, $2, $3, $4, $5)
	`

	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	member.MembershipTime = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		member.ID, member.ServiceID, member.UserID, member.Role, member.MembershipTime,
	)
	if err != nil {
		return fmt.Errorf("failed to add service member: %w", err)
	}

	log.Debug().Str("id", member.ID.String()).Msg("service member added")
	return nil
}

func (r *InventoryRepository) RemoveServiceMember(ctx context.Context, serviceID, userID uuid.UUID) error {
	query := `DELETE FROM service_members WHERE service_id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, serviceID, userID)
	if err != nil {
		return fmt.Errorf("failed to remove service member: %w", err)
	}

	log.Debug().Str("service_id", serviceID.String()).Str("user_id", userID.String()).Msg("service member removed")
	return nil
}

func (r *InventoryRepository) GetServiceMembers(ctx context.Context, serviceID uuid.UUID, limit, offset int) ([]domain.ServiceMember, error) {
	query := `
		SELECT id, service_id, user_id, role, membership_time
		FROM service_members
		WHERE service_id = $1
		ORDER BY membership_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, serviceID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get service members: %w", err)
	}
	defer rows.Close()

	var members []domain.ServiceMember
	for rows.Next() {
		var m domain.ServiceMember
		err := rows.Scan(&m.ID, &m.ServiceID, &m.UserID, &m.Role, &m.MembershipTime)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service member: %w", err)
		}
		members = append(members, m)
	}

	return members, rows.Err()
}

func (r *InventoryRepository) GetServiceUserIDs(ctx context.Context, serviceID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT user_id FROM service_members WHERE service_id = $1`

	rows, err := r.db.QueryContext(ctx, query, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service user ids: %w", err)
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan user id: %w", err)
		}
		userIDs = append(userIDs, id)
	}

	return userIDs, rows.Err()
}

func (r *InventoryRepository) RemoveServiceMembers(ctx context.Context, serviceID uuid.UUID, userIDs []uuid.UUID) error {
	query := `DELETE FROM service_members WHERE service_id = $1 AND user_id = $2`

	for _, userID := range userIDs {
		_, err := r.db.ExecContext(ctx, query, serviceID, userID)
		if err != nil {
			return fmt.Errorf("failed to remove service member: %w", err)
		}
	}

	log.Debug().Str("service_id", serviceID.String()).Msg("service members removed")
	return nil
}
