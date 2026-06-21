package domain

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrInvalidInput      = errors.New("invalid input")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrEmailAlreadyInUse = errors.New("email already in use")
	ErrInvalidCreds      = errors.New("invalid credentials")
	ErrTaskNotBelong     = errors.New("task does not belong to user")
)
