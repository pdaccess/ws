package domain

import (
	"time"

	"github.com/google/uuid"
)

type AlarmSeverity string

const (
	AlarmSeverityCritical AlarmSeverity = "critical"
	AlarmSeverityInfo     AlarmSeverity = "info"
	AlarmSeverityWarning  AlarmSeverity = "warning"
)

type Alarm struct {
	ID           uuid.UUID     `json:"id"`
	UserID       uuid.UUID     `json:"userId"`
	Title        string        `json:"title"`
	Message      string        `json:"message"`
	Source       string        `json:"source,omitempty"`
	Severity     AlarmSeverity `json:"severity"`
	Acknowledged bool          `json:"acknowledged"`
	Time         time.Time     `json:"time"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

type AlarmSearch struct {
	UserID   *uuid.UUID
	Severity *AlarmSeverity
	Limit    int
	Offset   int
}

type AlarmSearchOption func(*AlarmSearch)

func WithAlarmSeverity(severity AlarmSeverity) AlarmSearchOption {
	return func(s *AlarmSearch) {
		s.Severity = &severity
	}
}

func WithAlarmLimit(limit int) AlarmSearchOption {
	return func(s *AlarmSearch) {
		s.Limit = limit
	}
}

func WithAlarmOffset(offset int) AlarmSearchOption {
	return func(s *AlarmSearch) {
		s.Offset = offset
	}
}
