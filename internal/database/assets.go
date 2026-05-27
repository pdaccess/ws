package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pdaccess/ws/internal/core/domain"
)

type AssetRepository struct {
	db *DB
}

func NewAssetRepository(db *DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) CreateAsset(ctx context.Context, asset *domain.Asset) error {
	query := `INSERT INTO ws_assets (id, type, name, owner_id, parent_id, spec, embedding, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if asset.ID == uuid.Nil {
		asset.ID = uuid.New()
	}
	if asset.CreatedAt.IsZero() {
		asset.CreatedAt = time.Now()
	}
	spec := "{}"
	if asset.Spec != nil {
		spec = string(asset.Spec)
	}
	emb := []float64{}
	if asset.Embedding != nil {
		emb = []float64(*asset.Embedding)
	}
	_, err := r.db.ExecContext(ctx, query, asset.ID, string(asset.Type), asset.Name, asset.OwnerID, asset.ParentID, spec, pq.Array(emb), asset.CreatedAt)
	return err
}

func (r *AssetRepository) UpdateAssetEmbedding(ctx context.Context, id uuid.UUID, vec domain.Vector) error {
	_, err := r.db.ExecContext(ctx, `UPDATE ws_assets SET embedding = $1 WHERE id = $2`, pq.Array([]float64(vec)), id)
	return err
}

func (r *AssetRepository) GetAsset(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	a := &domain.Asset{}
	err := scanAsset(r.db.QueryRowContext(ctx, `SELECT id, type, name, owner_id, parent_id, spec, embedding, created_at FROM ws_assets WHERE id = $1`, id), a)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *AssetRepository) SearchAssets(ctx context.Context, assetType string, parentID *uuid.UUID, limit, offset int) ([]domain.Asset, error) {
	query := `SELECT id, type, name, owner_id, parent_id, spec, embedding, created_at FROM ws_assets WHERE 1=1`
	args := []any{}
	argNum := 1

	if assetType != "" {
		query += fmt.Sprintf(" AND type = $%d", argNum)
		args = append(args, assetType)
		argNum++
	}
	if parentID != nil {
		query += fmt.Sprintf(" AND parent_id = $%d", argNum)
		args = append(args, *parentID)
		argNum++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAssets(rows)
}

func (r *AssetRepository) HybridSearchAssets(ctx context.Context, query string, limit, offset int) ([]domain.Asset, error) {
	q := `SELECT id, type, name, owner_id, parent_id, spec, embedding, created_at FROM ws_assets WHERE name ILIKE $1 OR spec::text ILIKE $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	pattern := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, q, pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAssets(rows)
}

func (r *AssetRepository) VectorSearchAssets(ctx context.Context, embedding domain.Vector, limit, offset int) ([]domain.Asset, error) {
	q := `SELECT id, type, name, owner_id, parent_id, spec, embedding, created_at FROM ws_assets WHERE embedding != '{}' ORDER BY (SELECT SUM(v1 * v2) / GREATEST(SQRT(SUM(v1 * v1)) * SQRT(SUM(v2 * v2)), 0.000001) FROM unnest(embedding, $1::double precision[]) AS u(v1, v2)) DESC NULLS LAST LIMIT $2 OFFSET $3`
	rows, err := r.db.QueryContext(ctx, q, pq.Array([]float64(embedding)), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanAssets(rows)
}

// --- Identity: Users ---

func (r *AssetRepository) CreateUser(ctx context.Context, user *domain.User) error {
	existing, err := r.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return domain.ValidationError{Field: "email", Message: "email already exists", Code: domain.ErrCodeValidation}
	}

	query := `INSERT INTO ws_users (id, username, email, status, display_name, first_name, last_name, notification_settings) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	if user.Status == "" {
		user.Status = "active"
	}
	notif := user.NotificationSettings
	if notif == nil {
		notif = json.RawMessage(`{}`)
	}
	_, err = r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.Status, user.DisplayName, user.FirstName, user.LastName, notif)
	return err
}

func (r *AssetRepository) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	u := &domain.User{}
	var notif []byte
	err := r.db.QueryRowContext(ctx, `SELECT id, username, email, status, display_name, first_name, last_name, notification_settings FROM ws_users WHERE id = $1`, id).Scan(&u.ID, &u.Username, &u.Email, &u.Status, &u.DisplayName, &u.FirstName, &u.LastName, &notif)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(notif) > 0 {
		u.NotificationSettings = notif
	}
	return u, nil
}

func (r *AssetRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	var notif []byte
	err := r.db.QueryRowContext(ctx, `SELECT id, username, email, status, display_name, first_name, last_name, notification_settings FROM ws_users WHERE email = $1`, email).Scan(&u.ID, &u.Username, &u.Email, &u.Status, &u.DisplayName, &u.FirstName, &u.LastName, &notif)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if len(notif) > 0 {
		u.NotificationSettings = notif
	}
	return u, nil
}

func (r *AssetRepository) ListUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, username, email, status, display_name, first_name, last_name, notification_settings FROM ws_users ORDER BY username`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		var notif []byte
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Status, &u.DisplayName, &u.FirstName, &u.LastName, &notif); err != nil {
			return nil, err
		}
		if len(notif) > 0 {
			u.NotificationSettings = notif
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *AssetRepository) UpdateUser(ctx context.Context, id uuid.UUID, updates map[string]any) error {
	if len(updates) == 0 {
		return nil
	}
	var query strings.Builder
	query.WriteString(`UPDATE ws_users SET `)
	args := []any{}
	i := 1
	for col, val := range updates {
		if i > 1 {
			query.WriteString(`, `)
		}
		query.WriteString(col + ` = $` + fmt.Sprintf("%d", i))
		args = append(args, val)
		i++
	}
	query.WriteString(` WHERE id = $` + fmt.Sprintf("%d", i))
	args = append(args, id)
	_, err := r.db.ExecContext(ctx, query.String(), args...)
	return err
}

// --- Identity: Groups ---

func (r *AssetRepository) CreateGroup(ctx context.Context, group *domain.UserGroup) error {
	query := `INSERT INTO ws_groups (id, name, description) VALUES ($1, $2, $3)`
	if group.ID == uuid.Nil {
		group.ID = uuid.New()
	}
	_, err := r.db.ExecContext(ctx, query, group.ID, group.Name, group.Description)
	return err
}

func (r *AssetRepository) GetGroup(ctx context.Context, id uuid.UUID) (*domain.UserGroup, error) {
	group := &domain.UserGroup{}
	err := r.db.QueryRowContext(ctx, `SELECT id, name, description FROM ws_groups WHERE id = $1`, id).Scan(&group.ID, &group.Name, &group.Description)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *AssetRepository) ListGroups(ctx context.Context) ([]domain.UserGroup, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, description FROM ws_groups ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []domain.UserGroup
	for rows.Next() {
		var g domain.UserGroup
		if err := rows.Scan(&g.ID, &g.Name, &g.Description); err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}
	return groups, rows.Err()
}

// --- Identity: Group Memberships ---

func (r *AssetRepository) AddGroupMember(ctx context.Context, groupID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO ws_group_members (group_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, groupID, userID)
	return err
}

func (r *AssetRepository) RemoveGroupMember(ctx context.Context, groupID, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ws_group_members WHERE group_id = $1 AND user_id = $2`, groupID, userID)
	return err
}

func (r *AssetRepository) ListGroupMemberships(ctx context.Context, groupID uuid.UUID) ([]domain.GroupMember, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT group_id, user_id, created_at FROM ws_group_members WHERE group_id = $1 ORDER BY created_at`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []domain.GroupMember
	for rows.Next() {
		var m domain.GroupMember
		if err := rows.Scan(&m.GroupID, &m.MemberID, &m.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

// --- Vault Memberships ---

func (r *AssetRepository) CreateVaultMembership(ctx context.Context, vm *domain.VaultMembership) error {
	query := `INSERT INTO ws_vault_memberships (vault_id, member_id, member_type, role) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, vm.VaultID, vm.MemberID, vm.MemberType, vm.Role)
	return err
}

func (r *AssetRepository) ListVaultMemberships(ctx context.Context, vaultID uuid.UUID) ([]domain.VaultMembership, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT vault_id, member_id, member_type, role, created_at FROM ws_vault_memberships WHERE vault_id = $1`, vaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vms []domain.VaultMembership
	for rows.Next() {
		var vm domain.VaultMembership
		if err := rows.Scan(&vm.VaultID, &vm.MemberID, &vm.MemberType, &vm.Role, &vm.CreatedAt); err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}
	return vms, rows.Err()
}

// --- Admin Config ---

func (r *AssetRepository) GetAdminConfig(ctx context.Context) ([]domain.AdminConfig, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT key, value FROM ws_admin_config ORDER BY key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []domain.AdminConfig
	for rows.Next() {
		var c domain.AdminConfig
		if err := rows.Scan(&c.Key, &c.Value); err != nil {
			return nil, err
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}

func (r *AssetRepository) UpsertAdminConfig(ctx context.Context, key, value string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO ws_admin_config (key, value, updated_at) VALUES ($1, $2, NOW()) ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = NOW()`, key, value)
	return err
}

// --- Effective Policies ---

func (r *AssetRepository) AddAssetEffectivePolicies(ctx context.Context, assetID uuid.UUID, policies []domain.EffectivePolicy) error {
	for _, p := range policies {
		var groupID *uuid.UUID
		if p.GroupID != nil {
			gid := *p.GroupID
			groupID = &gid
		}
		actions := pq.Array(p.Actions)
		if _, err := r.db.ExecContext(ctx,
			`INSERT INTO ws_effective_policies (user_id, group_id, asset_id, actions, policy_id) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
			p.UserID, groupID, assetID, actions, p.PolicyID,
		); err != nil {
			return err
		}
	}
	return nil
}

func (r *AssetRepository) SetAssetEffectivePolicies(ctx context.Context, assetID uuid.UUID, policies []domain.EffectivePolicy) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM ws_effective_policies WHERE asset_id = $1`, assetID); err != nil {
		return err
	}

	for _, p := range policies {
		var groupID *uuid.UUID
		if p.GroupID != nil {
			gid := *p.GroupID
			groupID = &gid
		}
		actions := pq.Array(p.Actions)
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO ws_effective_policies (user_id, group_id, asset_id, actions, policy_id) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
			p.UserID, groupID, assetID, actions, p.PolicyID,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *AssetRepository) ListNonPolicyAssets(ctx context.Context) ([]domain.Asset, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, type, name, owner_id, parent_id, spec, embedding, created_at FROM ws_assets WHERE type != 'policy' ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAssets(rows)
}

func (r *AssetRepository) ListPolicies(ctx context.Context) ([]domain.Asset, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, type, name, owner_id, parent_id, spec, embedding, created_at FROM ws_assets WHERE type = 'policy' ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAssets(rows)
}

type scannable interface {
	Scan(dest ...any) error
}

func scanAsset(row scannable, a *domain.Asset) error {
	var spec sql.NullString
	var emb []float64
	var parentID sql.NullString
	err := row.Scan(&a.ID, &a.Type, &a.Name, &a.OwnerID, &parentID, &spec, pq.Array(&emb), &a.CreatedAt)
	if err != nil {
		return err
	}
	if parentID.Valid {
		id, _ := uuid.Parse(parentID.String)
		a.ParentID = &id
	}
	if spec.Valid {
		a.Spec = json.RawMessage(spec.String)
	}
	if len(emb) > 0 {
		v := domain.Vector(emb)
		a.Embedding = &v
	}
	return nil
}

func scanAssets(rows *sql.Rows) ([]domain.Asset, error) {
	var assets []domain.Asset
	for rows.Next() {
		var a domain.Asset
		if err := scanAsset(rows, &a); err != nil {
			return nil, err
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}
