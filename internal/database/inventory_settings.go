package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"git.h2hsecure.com/core/ws/internal/core/domain"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type InventorySettingsRepository struct {
	db *DB
}

func NewInventorySettingsRepository(db *DB) *InventorySettingsRepository {
	return &InventorySettingsRepository{db: db}
}

func (r *InventorySettingsRepository) Upsert(ctx context.Context, settings *domain.InventorySettings) error {
	query := `
		INSERT INTO inventory_settings (id, inventory_id, access_protocol, auth_protocol, vendor, version, host, port, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (inventory_id) DO UPDATE SET
			access_protocol = EXCLUDED.access_protocol,
			auth_protocol = EXCLUDED.auth_protocol,
			vendor = EXCLUDED.vendor,
			version = EXCLUDED.version,
			host = EXCLUDED.host,
			port = EXCLUDED.port,
			updated_at = EXCLUDED.updated_at
	`

	if settings.ID == uuid.Nil {
		settings.ID = uuid.New()
	}
	now := time.Now()
	settings.CreatedAt = now
	settings.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		settings.ID, settings.InventoryID, settings.AccessProtocol, settings.AuthProtocol,
		settings.Vendor, settings.Version, settings.Host, settings.Port,
		settings.CreatedAt, settings.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert inventory settings: %w", err)
	}

	log.Debug().Str("inventory_id", settings.InventoryID.String()).Msg("inventory settings upserted")
	return nil
}

func (r *InventorySettingsRepository) GetByInventoryID(ctx context.Context, inventoryID uuid.UUID) (*domain.InventorySettings, error) {
	query := `
		SELECT id, inventory_id, access_protocol, auth_protocol, vendor, version, host, port, created_at, updated_at
		FROM inventory_settings
		WHERE inventory_id = $1
	`

	settings := &domain.InventorySettings{}
	err := r.db.QueryRowContext(ctx, query, inventoryID).Scan(
		&settings.ID, &settings.InventoryID, &settings.AccessProtocol, &settings.AuthProtocol,
		&settings.Vendor, &settings.Version, &settings.Host, &settings.Port,
		&settings.CreatedAt, &settings.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory settings: %w", err)
	}

	return settings, nil
}

func (r *InventorySettingsRepository) Delete(ctx context.Context, inventoryID uuid.UUID) error {
	query := `DELETE FROM inventory_settings WHERE inventory_id = $1`

	_, err := r.db.ExecContext(ctx, query, inventoryID)
	if err != nil {
		return fmt.Errorf("failed to delete inventory settings: %w", err)
	}

	log.Debug().Str("inventory_id", inventoryID.String()).Msg("inventory settings deleted")
	return nil
}
