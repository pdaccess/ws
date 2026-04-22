package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/rs/zerolog/log"
)

type CredentialRepository struct {
	db *DB
}

func NewCredentialRepository(db *DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

func (r *CredentialRepository) Create(ctx context.Context, cred *domain.Credential) error {
	query := `
		INSERT INTO credentials (id, group_id, name, description, type, metadata, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	if cred.ID == uuid.Nil {
		cred.ID = uuid.New()
	}
	now := time.Now()
	cred.CreatedAt = now
	cred.UpdatedAt = now

	var metadata []byte
	if cred.Metadata != nil {
		metadata, _ = json.Marshal(cred.Metadata)
	} else {
		metadata = []byte("{}")
	}

	_, err := r.db.ExecContext(ctx, query,
		cred.ID, cred.GroupID, cred.Name, cred.Description, cred.Type, metadata, cred.IsActive,
		cred.CreatedAt, cred.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create credential: %w", err)
	}

	log.Debug().Str("id", cred.ID.String()).Msg("credential created")
	return nil
}

func (r *CredentialRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Credential, error) {
	query := `
		SELECT id, group_id, name, description, type, metadata, is_active, created_at, updated_at, deleted_at
		FROM credentials
		WHERE id = $1 AND deleted_at IS NULL
	`

	cred := &domain.Credential{}
	var metadata []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&cred.ID, &cred.GroupID, &cred.Name, &cred.Description, &cred.Type, &metadata,
		&cred.IsActive, &cred.CreatedAt, &cred.UpdatedAt, &cred.DeletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credential: %w", err)
	}
	if metadata != nil {
		cred.Metadata = json.RawMessage(metadata)
	}

	return cred, nil
}

func (r *CredentialRepository) Update(ctx context.Context, cred *domain.Credential) error {
	cred.UpdatedAt = time.Now()

	var metadata []byte
	if cred.Metadata != nil {
		metadata, _ = json.Marshal(cred.Metadata)
	} else {
		metadata = []byte("{}")
	}

	query := `
		UPDATE credentials SET name = $2, description = $3, type = $4, metadata = $5, is_active = $6, updated_at = $7
		WHERE id = $1 AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		cred.ID, cred.Name, cred.Description, cred.Type, metadata, cred.IsActive, cred.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update credential: %w", err)
	}

	log.Debug().Str("id", cred.ID.String()).Msg("credential updated")
	return nil
}

func (r *CredentialRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE credentials SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete credential: %w", err)
	}

	log.Debug().Str("id", id.String()).Msg("credential deleted")
	return nil
}

func (r *CredentialRepository) Search(ctx context.Context, opts ...domain.CredentialSearchOption) ([]domain.Credential, error) {
	search := &domain.CredentialSearch{
		Limit:  20,
		Offset: 0,
	}
	for _, opt := range opts {
		opt(search)
	}

	query := `
		SELECT id, group_id, name, description, type, metadata, is_active, created_at, updated_at, deleted_at
		FROM credentials
		WHERE 1=1
	`

	args := []any{}
	argCount := 0

	if search.GroupID != nil {
		argCount++
		query += fmt.Sprintf(" AND group_id = $%d", argCount)
		args = append(args, *search.GroupID)
	}

	if !search.Deleted {
		query += " AND deleted_at IS NULL"
	} else {
		query += " AND deleted_at IS NOT NULL"
	}

	if search.Name != nil && *search.Name != "" {
		argCount++
		filterPattern := "%" + strings.ToLower(*search.Name) + "%"
		query += fmt.Sprintf(" AND (LOWER(name) LIKE $%d OR LOWER(description) LIKE $%d)", argCount, argCount)
		args = append(args, filterPattern)
	}

	if search.Type != "" {
		argCount++
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, string(search.Type))
	}

	if search.IsActive != nil {
		argCount++
		query += fmt.Sprintf(" AND is_active = $%d", argCount)
		args = append(args, *search.IsActive)
	}

	query += " ORDER BY created_at DESC"

	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, search.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, search.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search credentials: %w", err)
	}
	defer rows.Close()

	var items []domain.Credential
	for rows.Next() {
		var c domain.Credential
		var metadata []byte
		err := rows.Scan(
			&c.ID, &c.GroupID, &c.Name, &c.Description, &c.Type, &metadata,
			&c.IsActive, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan credential: %w", err)
		}
		if metadata != nil {
			c.Metadata = json.RawMessage(metadata)
		}
		items = append(items, c)
	}

	return items, rows.Err()
}

func (r *CredentialRepository) CreateSecret(ctx context.Context, secret *domain.CredentialSecret) error {
	query := `
		INSERT INTO credential_secrets (id, credential_id, username, password, private_key, public_key, api_key, api_secret, certificate, private_key_pass, expires_at, last_rotated, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	if secret.ID == uuid.Nil {
		secret.ID = uuid.New()
	}
	now := time.Now()
	secret.CreatedAt = now
	secret.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		secret.ID, secret.CredentialID, secret.Username, secret.Password, secret.PrivateKey,
		secret.PublicKey, secret.APIKey, secret.APIsecret, secret.Certificate, secret.PrivateKeyPass,
		secret.ExpiresAt, secret.LastRotated, secret.CreatedAt, secret.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create credential secret: %w", err)
	}

	log.Debug().Str("id", secret.ID.String()).Msg("credential secret created")
	return nil
}

func (r *CredentialRepository) GetSecretByCredentialID(ctx context.Context, credentialID uuid.UUID) (*domain.CredentialSecret, error) {
	query := `
		SELECT id, credential_id, username, password, private_key, public_key, api_key, api_secret, certificate, private_key_pass, expires_at, last_rotated, created_at, updated_at
		FROM credential_secrets
		WHERE credential_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`

	secret := &domain.CredentialSecret{}
	err := r.db.QueryRowContext(ctx, query, credentialID).Scan(
		&secret.ID, &secret.CredentialID, &secret.Username, &secret.Password, &secret.PrivateKey,
		&secret.PublicKey, &secret.APIKey, &secret.APIsecret, &secret.Certificate, &secret.PrivateKeyPass,
		&secret.ExpiresAt, &secret.LastRotated, &secret.CreatedAt, &secret.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get credential secret: %w", err)
	}

	return secret, nil
}

func (r *CredentialRepository) UpdateSecret(ctx context.Context, secret *domain.CredentialSecret) error {
	secret.UpdatedAt = time.Now()

	query := `
		UPDATE credential_secrets SET 
			username = $2, password = $3, private_key = $4, public_key = $5, 
			api_key = $6, api_secret = $7, certificate = $8, private_key_pass = $9,
			expires_at = $10, last_rotated = $11, updated_at = $12
		WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		secret.ID, secret.Username, secret.Password, secret.PrivateKey, secret.PublicKey,
		secret.APIKey, secret.APIsecret, secret.Certificate, secret.PrivateKeyPass,
		secret.ExpiresAt, secret.LastRotated, secret.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update credential secret: %w", err)
	}

	log.Debug().Str("id", secret.ID.String()).Msg("credential secret updated")
	return nil
}

func (r *CredentialRepository) DeleteSecret(ctx context.Context, credentialID uuid.UUID) error {
	query := `DELETE FROM credential_secrets WHERE credential_id = $1`

	_, err := r.db.ExecContext(ctx, query, credentialID)
	if err != nil {
		return fmt.Errorf("failed to delete credential secret: %w", err)
	}

	log.Debug().Str("credential_id", credentialID.String()).Msg("credential secret deleted")
	return nil
}
