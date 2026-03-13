package domain

import (
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	RealmID    uuid.UUID `json:"realmId"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID uuid.UUID `json:"resourceId"`
	Details    string    `json:"details,omitempty"`
	IPAddress  string    `json:"ipAddress,omitempty"`
	Time       time.Time `json:"time"`
	CreatedAt  time.Time `json:"createdAt"`
}

type ActivitySearch struct {
	UserID    *uuid.UUID
	RealmID   *uuid.UUID
	Action    *string
	Limit     int
	Offset    int
	StartDate *time.Time
	EndDate   *time.Time
}

type ActivitySearchOption func(*ActivitySearch)

func WithActivityUserID(userID uuid.UUID) ActivitySearchOption {
	return func(s *ActivitySearch) {
		s.UserID = &userID
	}
}

func WithActivityRealmID(realmID uuid.UUID) ActivitySearchOption {
	return func(s *ActivitySearch) {
		s.RealmID = &realmID
	}
}

func WithActivityAction(action string) ActivitySearchOption {
	return func(s *ActivitySearch) {
		s.Action = &action
	}
}

func WithActivityLimit(limit int) ActivitySearchOption {
	return func(s *ActivitySearch) {
		s.Limit = limit
	}
}

func WithActivityOffset(offset int) ActivitySearchOption {
	return func(s *ActivitySearch) {
		s.Offset = offset
	}
}
