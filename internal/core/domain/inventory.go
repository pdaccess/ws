package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Vector []float64

type InventoryType string

const (
	InventoryTypeGroup   InventoryType = "group"
	InventoryTypeService InventoryType = "service"
	InventoryTypeVault   InventoryType = "vault"
)

type Group struct {
	ID          uuid.UUID     `json:"id"`
	RealmID     uuid.UUID     `json:"realmId"`
	ParentID    *uuid.UUID    `json:"parentId,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Type        InventoryType `json:"type,omitempty"`
	Embedding   Vector        `json:"embedding,omitempty"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	DeletedAt   *time.Time    `json:"deletedAt,omitempty"`
}

type GroupMember struct {
	ID             uuid.UUID `json:"id"`
	GroupID        uuid.UUID `json:"groupId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
	MembershipTime time.Time `json:"membershipTime"`
}

type GroupSearch struct {
	RealmID   *uuid.UUID
	ParentID  *uuid.UUID
	Name      *string
	Deleted   bool
	Limit     int
	Offset    int
	StartDate *time.Time
	EndDate   *time.Time
	Filter    *string
	Vector    Vector
	Type      InventoryType
}

type GroupSearchOption func(*GroupSearch)

func WithGroupRealmID(id uuid.UUID) GroupSearchOption {
	return func(s *GroupSearch) {
		s.RealmID = &id
	}
}

func WithGroupParentID(id uuid.UUID) GroupSearchOption {
	return func(s *GroupSearch) {
		s.ParentID = &id
	}
}

func WithGroupName(name string) GroupSearchOption {
	return func(s *GroupSearch) {
		s.Name = &name
	}
}

func WithGroupDeleted(deleted bool) GroupSearchOption {
	return func(s *GroupSearch) {
		s.Deleted = deleted
	}
}

func WithGroupLimit(limit int) GroupSearchOption {
	return func(s *GroupSearch) {
		s.Limit = limit
	}
}

func WithGroupOffset(offset int) GroupSearchOption {
	return func(s *GroupSearch) {
		s.Offset = offset
	}
}

func WithGroupDateRange(start, end time.Time) GroupSearchOption {
	return func(s *GroupSearch) {
		s.StartDate = &start
		s.EndDate = &end
	}
}

func WithGroupFilter(filter string) GroupSearchOption {
	return func(s *GroupSearch) {
		s.Filter = &filter
	}
}

func WithGroupVector(vec Vector) GroupSearchOption {
	return func(s *GroupSearch) {
		s.Vector = vec
	}
}

type Service struct {
	ID          uuid.UUID        `json:"id"`
	RealmID     uuid.UUID        `json:"realmId"`
	ParentID    *uuid.UUID       `json:"parentId,omitempty"`
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Type        InventoryType    `json:"type,omitempty"`
	Embedding   Vector           `json:"embedding,omitempty"`
	Settings    *ServiceSettings `json:"settings,omitempty"`
	Credentials []Credential     `json:"credentials,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	DeletedAt   *time.Time       `json:"deletedAt,omitempty"`
}

type ServiceSettings struct {
	ID             uuid.UUID `json:"id"`
	ServiceID      uuid.UUID `json:"serviceId"`
	AccessProtocol string    `json:"accessProtocol,omitempty"`
	AuthProtocol   string    `json:"authProtocol,omitempty"`
	Vendor         string    `json:"vendor,omitempty"`
	Version        string    `json:"version,omitempty"`
	Host           string    `json:"host,omitempty"`
	Port           *int      `json:"port,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type ServiceMember struct {
	ID             uuid.UUID `json:"id"`
	ServiceID      uuid.UUID `json:"serviceId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
	MembershipTime time.Time `json:"membershipTime"`
}

type ServiceSearch struct {
	RealmID   *uuid.UUID
	ParentID  *uuid.UUID
	Name      *string
	Deleted   bool
	Limit     int
	Offset    int
	StartDate *time.Time
	EndDate   *time.Time
	Filter    *string
	Vector    Vector
	Type      InventoryType
}

type ServiceSearchOption func(*ServiceSearch)

func WithServiceRealmID(id uuid.UUID) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.RealmID = &id
	}
}

func WithServiceParentID(id uuid.UUID) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.ParentID = &id
	}
}

func WithServiceName(name string) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.Name = &name
	}
}

func WithServiceDeleted(deleted bool) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.Deleted = deleted
	}
}

func WithServiceLimit(limit int) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.Limit = limit
	}
}

func WithServiceOffset(offset int) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.Offset = offset
	}
}

func WithServiceDateRange(start, end time.Time) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.StartDate = &start
		s.EndDate = &end
	}
}

func WithServiceFilter(filter string) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.Filter = &filter
	}
}

func WithServiceVector(vec Vector) ServiceSearchOption {
	return func(s *ServiceSearch) {
		s.Vector = vec
	}
}

type Credential struct {
	ID          uuid.UUID       `json:"id"`
	GroupID     uuid.UUID       `json:"groupId"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Type        CredentialType  `json:"type"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	DeletedAt   *time.Time      `json:"deletedAt,omitempty"`
}

type CredentialType string

const (
	CredentialTypePassword    CredentialType = "password"
	CredentialTypeSSHKey      CredentialType = "ssh_key"
	CredentialTypeAPIKey      CredentialType = "api_key"
	CredentialTypeCertificate CredentialType = "certificate"
	CredentialTypeOAuth       CredentialType = "oauth"
)

type CredentialSecret struct {
	ID             uuid.UUID  `json:"id"`
	CredentialID   uuid.UUID  `json:"credentialId"`
	Username       string     `json:"username,omitempty"`
	Password       string     `json:"password,omitempty"`
	PrivateKey     string     `json:"privateKey,omitempty"`
	PublicKey      string     `json:"publicKey,omitempty"`
	APIKey         string     `json:"apiKey,omitempty"`
	APIsecret      string     `json:"apiSecret,omitempty"`
	Certificate    string     `json:"certificate,omitempty"`
	PrivateKeyPass string     `json:"privateKeyPass,omitempty"`
	ExpiresAt      *time.Time `json:"expiresAt,omitempty"`
	LastRotated    *time.Time `json:"lastRotated,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type CredentialSearch struct {
	GroupID  *uuid.UUID
	Name     *string
	Type     CredentialType
	IsActive *bool
	Deleted  bool
	Limit    int
	Offset   int
}

type CredentialSearchOption func(*CredentialSearch)

func WithCredentialGroupID(id uuid.UUID) CredentialSearchOption {
	return func(s *CredentialSearch) {
		s.GroupID = &id
	}
}

func WithCredentialName(name string) CredentialSearchOption {
	return func(s *CredentialSearch) {
		s.Name = &name
	}
}

func WithCredentialType(t CredentialType) CredentialSearchOption {
	return func(s *CredentialSearch) {
		s.Type = t
	}
}

func WithCredentialActive(active bool) CredentialSearchOption {
	return func(s *CredentialSearch) {
		s.IsActive = &active
	}
}

func WithCredentialLimit(limit int) CredentialSearchOption {
	return func(s *CredentialSearch) {
		s.Limit = limit
	}
}

func WithCredentialOffset(offset int) CredentialSearchOption {
	return func(s *CredentialSearch) {
		s.Offset = offset
	}
}
