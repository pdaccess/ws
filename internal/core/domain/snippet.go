package domain

import "time"

type Snippet struct {
	Content   string
	UserId    string
	CreatedAt time.Time
	Expire    time.Duration
}

type SnippetSearch struct {
	Id *string
}

type SnippetSearchOption func(*SnippetSearch)

func WithSnippetById(Id string) SnippetSearchOption {
	return func(s *SnippetSearch) {
		s.Id = &Id
	}
}
