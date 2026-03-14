package domain

import "errors"

var (
	ErrInvalidID    = errors.New("PDA-1001")
	ErrNotFound     = errors.New("PDA-1002")
	ErrValidation   = errors.New("PDA-1003")
	ErrUnauthorized = errors.New("PDA-1004")
	ErrInternal     = errors.New("PDA-1005")
)

const (
	ErrCodeInvalidID    string = "PDA-1001"
	ErrCodeNotFound     string = "PDA-1002"
	ErrCodeValidation   string = "PDA-1003"
	ErrCodeUnauthorized string = "PDA-1004"
	ErrCodeInternal     string = "PDA-1005"
)

type ValidationError struct {
	Field   string
	Message string
	Code    string
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e ValidationError) Unwrap() error {
	return ErrValidation
}

type NotFoundError struct {
	Resource string
	ID       any
	Code     string
}

func (e NotFoundError) Error() string {
	return e.Resource + " not found"
}

func (e NotFoundError) Unwrap() error {
	return ErrNotFound
}

type InvalidIDError struct {
	Message string
	Code    string
}

func (e InvalidIDError) Error() string {
	return e.Message
}

func (e InvalidIDError) Unwrap() error {
	return ErrInvalidID
}

type InternalError struct {
	Message string
	Code    string
	cause   error
}

func (e InternalError) Error() string {
	return e.Message
}

func (e InternalError) Unwrap() error {
	return e.cause
}

func (e InternalError) Cause() error {
	return e.cause
}

func NewInternalError(message string, cause error) InternalError {
	return InternalError{
		Message: message,
		Code:    ErrCodeInternal,
		cause:   cause,
	}
}
