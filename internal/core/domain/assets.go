package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AssetType string

const (
	AssetTypeServiceGroup AssetType = "service_group"
	AssetTypeService      AssetType = "service"
	AssetTypeVault        AssetType = "vault"
	AssetTypeSecret       AssetType = "secret"
	AssetTypePolicy       AssetType = "policy"
	AssetTypePaste        AssetType = "paste"
	AssetTypeJIT          AssetType = "jit"
)

type Asset struct {
	ID        uuid.UUID       `json:"id"`
	Type      AssetType       `json:"type"`
	Name      string          `json:"name"`
	OwnerID   uuid.UUID       `json:"owner_id"`
	ParentID  *uuid.UUID      `json:"parent_id,omitempty"`
	Spec      json.RawMessage `json:"spec,omitempty"`
	Embedding *Vector         `json:"embedding,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

type User struct {
	ID                   uuid.UUID       `json:"id"`
	Username             string          `json:"username"`
	Email                string          `json:"email"`
	Status               string          `json:"status,omitempty"`
	DisplayName          string          `json:"display_name,omitempty"`
	FirstName            string          `json:"first_name,omitempty"`
	LastName             string          `json:"last_name,omitempty"`
	NotificationSettings json.RawMessage `json:"notification_settings,omitempty"`
}

type UserGroup struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
}

type GroupMember struct {
	GroupID   uuid.UUID `json:"groupId"`
	MemberID  uuid.UUID `json:"memberId"`
	CreatedAt time.Time `json:"createdAt"`
}

type VaultMembership struct {
	VaultID    uuid.UUID `json:"vaultId"`
	MemberID   uuid.UUID `json:"memberId"`
	MemberType string    `json:"memberType"`
	Role       string    `json:"role"`
	CreatedAt  time.Time `json:"createdAt"`
}

type AuditEntry struct {
	ID         uuid.UUID `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	ActorID    uuid.UUID `json:"actor_id"`
	Action     string    `json:"action"`
	ResourceID uuid.UUID `json:"resource_id"`
	Success    bool      `json:"success"`
}

type AdminConfig struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EffectivePolicy struct {
	UserID   uuid.UUID  `json:"user_id"`
	GroupID  *uuid.UUID `json:"group_id,omitempty"`
	AssetID  uuid.UUID  `json:"asset_id"`
	Actions  []string   `json:"actions"`
	PolicyID uuid.UUID  `json:"policy_id"`
}
