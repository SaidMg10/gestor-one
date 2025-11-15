package domain

import "errors"

var (
	ErrNotFound          = errors.New("record not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrEmailExists       = errors.New("email already exist")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
	ErrPasswordRequired  = errors.New("password is required")
)
