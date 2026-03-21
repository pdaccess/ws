package domain

import (
	"time"

	"github.com/google/uuid"
)

type Vector []float64

type Group struct {
	ID          uuid.UUID  `json:"id"`
	RealmID     uuid.UUID  `json:"realmId"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Embedding   Vector     `json:"embedding,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
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
	Embedding   Vector           `json:"embedding,omitempty"`
	Settings    *ServiceSettings `json:"settings,omitempty"`
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
