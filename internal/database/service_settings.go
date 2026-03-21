package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pdaccess/ws/internal/core/domain"
	"github.com/rs/zerolog/log"
)

type ServiceSettingsRepository struct {
	db *DB
}

func NewServiceSettingsRepository(db *DB) *ServiceSettingsRepository {
	return &ServiceSettingsRepository{db: db}
}

func (r *ServiceSettingsRepository) Upsert(ctx context.Context, settings *domain.ServiceSettings) error {
	query := `
		INSERT INTO services (id, inventory_id, access_protocol, auth_protocol, vendor, version, host, port, created_at, updated_at)
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
		settings.ID, settings.ServiceID, settings.AccessProtocol, settings.AuthProtocol,
		settings.Vendor, settings.Version, settings.Host, settings.Port,
		settings.CreatedAt, settings.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to upsert service settings: %w", err)
	}

	log.Debug().Str("id", settings.ID.String()).Msg("service settings upserted")
	return nil
}

func (r *ServiceSettingsRepository) GetByInventoryID(ctx context.Context, inventoryID uuid.UUID) (*domain.ServiceSettings, error) {
	query := `
		SELECT id, inventory_id, access_protocol, auth_protocol, vendor, version, host, port, created_at, updated_at
		FROM services
		WHERE inventory_id = $1
	`

	settings := &domain.ServiceSettings{}
	var accessProtocol, authProtocol, vendor, version, host sql.NullString
	var port sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, inventoryID).Scan(
		&settings.ID, &settings.ServiceID, &accessProtocol, &authProtocol,
		&vendor, &version, &host, &port, &settings.CreatedAt, &settings.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get service settings: %w", err)
	}

	if accessProtocol.Valid {
		settings.AccessProtocol = accessProtocol.String
	}
	if authProtocol.Valid {
		settings.AuthProtocol = authProtocol.String
	}
	if vendor.Valid {
		settings.Vendor = vendor.String
	}
	if version.Valid {
		settings.Version = version.String
	}
	if host.Valid {
		settings.Host = host.String
	}
	if port.Valid {
		p := int(port.Int64)
		settings.Port = &p
	}

	return settings, nil
}

func (r *ServiceSettingsRepository) Delete(ctx context.Context, inventoryID uuid.UUID) error {
	query := `DELETE FROM services WHERE inventory_id = $1`

	_, err := r.db.ExecContext(ctx, query, inventoryID)
	if err != nil {
		return fmt.Errorf("failed to delete service settings: %w", err)
	}

	log.Debug().Str("inventory_id", inventoryID.String()).Msg("service settings deleted")
	return nil
}
