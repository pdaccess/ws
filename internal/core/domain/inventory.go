package domain

import (
	"time"

	"github.com/google/uuid"
)

type ItemType string

const (
	ItemTypeGroup   ItemType = "group"
	ItemTypeService ItemType = "service"
)

type Vector []float32

type Inventory struct {
	ID          uuid.UUID  `json:"id"`
	RealmID     uuid.UUID  `json:"realmId"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	ItemType    ItemType   `json:"itemType"`
	Embedding   Vector     `json:"embedding,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type InventorySettings struct {
	ID             uuid.UUID `json:"id"`
	InventoryID    uuid.UUID `json:"inventoryId"`
	AccessProtocol string    `json:"accessProtocol,omitempty"`
	AuthProtocol   string    `json:"authProtocol,omitempty"`
	Vendor         string    `json:"vendor,omitempty"`
	Version        string    `json:"version,omitempty"`
	Host           string    `json:"host,omitempty"`
	Port           *int      `json:"port,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type InventoryMember struct {
	ID             uuid.UUID `json:"id"`
	InventoryID    uuid.UUID `json:"inventoryId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
	MembershipTime time.Time `json:"membershipTime"`
}

type User struct {
	ID        uuid.UUID  `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type UserGroup struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

type UserGroupMember struct {
	ID             uuid.UUID `json:"id"`
	UserGroupID    uuid.UUID `json:"userGroupId"`
	UserID         uuid.UUID `json:"userId"`
	Role           string    `json:"role"`
	MembershipTime time.Time `json:"membershipTime"`
}

type InventorySearch struct {
	RealmID   *uuid.UUID
	ParentID  *uuid.UUID
	ItemType  *ItemType
	Name      *string
	Deleted   bool
	Limit     int
	Offset    int
	StartDate *time.Time
	EndDate   *time.Time
	Filter    *string
	Vector    Vector
}

type InventorySearchOption func(*InventorySearch)

func WithRealmID(id uuid.UUID) InventorySearchOption {
	return func(s *InventorySearch) {
		s.RealmID = &id
	}
}

func WithParentID(id uuid.UUID) InventorySearchOption {
	return func(s *InventorySearch) {
		s.ParentID = &id
	}
}

func WithItemType(t ItemType) InventorySearchOption {
	return func(s *InventorySearch) {
		s.ItemType = &t
	}
}

func WithName(name string) InventorySearchOption {
	return func(s *InventorySearch) {
		s.Name = &name
	}
}

func WithDeleted(deleted bool) InventorySearchOption {
	return func(s *InventorySearch) {
		s.Deleted = deleted
	}
}

func WithLimit(limit int) InventorySearchOption {
	return func(s *InventorySearch) {
		s.Limit = limit
	}
}

func WithOffset(offset int) InventorySearchOption {
	return func(s *InventorySearch) {
		s.Offset = offset
	}
}

func WithDateRange(start, end time.Time) InventorySearchOption {
	return func(s *InventorySearch) {
		s.StartDate = &start
		s.EndDate = &end
	}
}

func WithFilter(filter string) InventorySearchOption {
	return func(s *InventorySearch) {
		s.Filter = &filter
	}
}

func WithVector(vec Vector) InventorySearchOption {
	return func(s *InventorySearch) {
		s.Vector = vec
	}
}
