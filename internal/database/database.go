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

func (d *DB) InventoryRepo() *InventoryRepository {
	return NewInventoryRepository(d)
}

func (d *DB) UserRepo() *UserRepository {
	return NewUserRepository(d)
}

func (d *DB) UserGroupRepo() *UserGroupRepository {
	return NewUserGroupRepository(d)
}

func (d *DB) ActivityRepo() *ActivityRepository {
	return NewActivityRepository(d)
}

func (d *DB) PasteRepo() *PasteRepository {
	return NewPasteRepository(d)
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
CREATE EXTENSION IF NOT EXISTS vector;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- User groups table
CREATE TABLE IF NOT EXISTS user_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_user_groups_deleted_at ON user_groups(deleted_at);

-- User group members table
CREATE TABLE IF NOT EXISTS user_group_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_group_id UUID NOT NULL REFERENCES user_groups(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('user', 'admin')),
    membership_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_group_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_group_members_user_group_id ON user_group_members(user_group_id);
CREATE INDEX IF NOT EXISTS idx_user_group_members_user_id ON user_group_members(user_id);

-- Inventory table (unified service and group)
CREATE TABLE IF NOT EXISTS inventory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    realm_id UUID NOT NULL,
    parent_id UUID REFERENCES inventory(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    item_type VARCHAR(50) NOT NULL CHECK (item_type IN ('group', 'service')),
    embedding vector(384),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_inventory_realm_id ON inventory(realm_id);
CREATE INDEX IF NOT EXISTS idx_inventory_parent_id ON inventory(parent_id);
CREATE INDEX IF NOT EXISTS idx_inventory_item_type ON inventory(item_type);
CREATE INDEX IF NOT EXISTS idx_inventory_deleted_at ON inventory(deleted_at);
CREATE INDEX IF NOT EXISTS idx_inventory_embedding ON inventory USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- Inventory settings table (specific fields for services)
CREATE TABLE IF NOT EXISTS inventory_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_id UUID NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    access_protocol VARCHAR(50) CHECK (access_protocol IN ('ssh', 'sql', 'vnc', 'rdp', 'http', 'none')),
    auth_protocol VARCHAR(50) CHECK (auth_protocol IN ('radius', 'oauth2', 'ldap', 'tacacs', 'none')),
    vendor VARCHAR(255),
    version VARCHAR(100),
    host VARCHAR(255),
    port INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_inventory_settings_inventory_id ON inventory_settings(inventory_id);

-- Group members table (users within inventory items)
CREATE TABLE IF NOT EXISTS inventory_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_id UUID NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL CHECK (role IN ('user', 'admin')),
    membership_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(inventory_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_inventory_members_inventory_id ON inventory_members(inventory_id);
CREATE INDEX IF NOT EXISTS idx_inventory_members_user_id ON inventory_members(user_id);

-- Inventory messages table (MOTD)
CREATE TABLE IF NOT EXISTS inventory_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_id UUID NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    background_color VARCHAR(50),
    font_color VARCHAR(50),
    font_size INTEGER,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inventory_messages_inventory_id ON inventory_messages(inventory_id);

-- Alarms table
CREATE TABLE IF NOT EXISTS alarms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_id UUID REFERENCES inventory(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(255),
    name VARCHAR(255),
    pattern VARCHAR(500),
    create_time TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alarms_user_id ON alarms(user_id);
CREATE INDEX IF NOT EXISTS idx_alarms_inventory_id ON alarms(inventory_id);

-- Snippets table
CREATE TABLE IF NOT EXISTS snippets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    content TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    marked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_snippets_user_id ON snippets(user_id);

-- Activities table
CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    realm_id UUID NOT NULL,
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    resource_id UUID,
    details TEXT,
    ip_address VARCHAR(45),
    activity_time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_activities_user_id ON activities(user_id);
CREATE INDEX IF NOT EXISTS idx_activities_realm_id ON activities(realm_id);
CREATE INDEX IF NOT EXISTS idx_activities_action ON activities(action);
CREATE INDEX IF NOT EXISTS idx_activities_activity_time ON activities(activity_time);

-- Pastes table (create if not exists, fix FK to allow NULL user_id)
DROP TABLE IF EXISTS pastes;
CREATE TABLE pastes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID,
    content TEXT NOT NULL,
    language VARCHAR(50),
    expires_at TIMESTAMP WITH TIME ZONE,
    views INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pastes_user_id ON pastes(user_id);
CREATE INDEX IF NOT EXISTS idx_pastes_expires_at ON pastes(expires_at);

-- Config contexts table
CREATE TABLE IF NOT EXISTS config_contexts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    realm VARCHAR(255) NOT NULL,
    context VARCHAR(255) NOT NULL,
    entry VARCHAR(255) NOT NULL,
    value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(realm, context, entry)
);

CREATE INDEX IF NOT EXISTS idx_config_contexts_realm_context ON config_contexts(realm, context);
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
