package domain

import (
	"time"

	"github.com/google/uuid"
)

type Paste struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	Content   string     `json:"content"`
	Language  string     `json:"language,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	Views     int        `json:"views"`
	CreatedAt time.Time  `json:"createdAt"`
}

type PasteSearch struct {
	UserID *uuid.UUID
	Limit  int
	Offset int
}

type PasteSearchOption func(*PasteSearch)

func WithPasteUserID(userID uuid.UUID) PasteSearchOption {
	return func(s *PasteSearch) {
		s.UserID = &userID
	}
}

func WithPasteLimit(limit int) PasteSearchOption {
	return func(s *PasteSearch) {
		s.Limit = limit
	}
}

func WithPasteOffset(offset int) PasteSearchOption {
	return func(s *PasteSearch) {
		s.Offset = offset
	}
}
