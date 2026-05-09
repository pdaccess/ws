package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type DB struct {
	*sql.DB
}

func (d *DB) AssetRepo() *AssetRepository {
	return NewAssetRepository(d)
}

func (d *DB) AuditRepo() *AuditRepository {
	return NewAuditRepository(d)
}

func New(connStr string) (*DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("database connection established")
	return &DB{db}, nil
}

func (d *DB) RunMigrations() error {
	schema := `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Unified assets table (service, service_group, vault, secret, policy, paste, jit)
CREATE TABLE IF NOT EXISTS ws_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL,
    parent_id UUID,
    spec JSONB DEFAULT '{}',
    embedding DOUBLE PRECISION[] DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ws_assets_type ON ws_assets(type);
CREATE INDEX IF NOT EXISTS idx_ws_assets_owner_id ON ws_assets(owner_id);

-- Identity: users
CREATE TABLE IF NOT EXISTS ws_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    display_name VARCHAR(255) DEFAULT '',
    first_name VARCHAR(255) DEFAULT '',
    last_name VARCHAR(255) DEFAULT '',
    notification_settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Identity: groups
CREATE TABLE IF NOT EXISTS ws_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Identity: group memberships
CREATE TABLE IF NOT EXISTS ws_group_members (
    group_id UUID NOT NULL REFERENCES ws_groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES ws_users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (group_id, user_id)
);

-- Vault memberships
CREATE TABLE IF NOT EXISTS ws_vault_memberships (
    vault_id UUID NOT NULL REFERENCES ws_assets(id) ON DELETE CASCADE,
    member_id UUID NOT NULL,
    member_type VARCHAR(20) NOT NULL CHECK (member_type IN ('user', 'vault')),
    role VARCHAR(20) NOT NULL CHECK (role IN ('owner', 'manager', 'viewer')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (vault_id, member_id)
);

-- Admin config (key-value)
CREATE TABLE IF NOT EXISTS ws_admin_config (
    key VARCHAR(255) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Audit logs
CREATE TABLE IF NOT EXISTS ws_audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    actor_id UUID NOT NULL,
    action VARCHAR(255) NOT NULL,
    resource_id UUID,
    success BOOLEAN DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_ws_audit_logs_actor_id ON ws_audit_logs(actor_id);
CREATE INDEX IF NOT EXISTS idx_ws_audit_logs_resource_id ON ws_audit_logs(resource_id);
CREATE INDEX IF NOT EXISTS idx_ws_audit_logs_timestamp ON ws_audit_logs(timestamp);
`

	_, err := d.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Info().Msg("database migrations completed")
	return nil
}

func (d *DB) InitDatabase(connStr string) error {
	u, err := url.Parse(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	dbName := u.Path[1:]
	if dbName == "" {
		return fmt.Errorf("database name not found in connection string")
	}

	tempDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	var exists bool
	err = tempDB.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		tempDB.Close()
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(dbName)))
		if err != nil {
			tempDB.Close()
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Info().Str("database", dbName).Msg("database created")
	}
	tempDB.Close()

	db, err := New(connStr)
	if err != nil {
		return err
	}
	d.DB = db.DB

	return d.RunMigrations()
}

func (d *DB) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}
