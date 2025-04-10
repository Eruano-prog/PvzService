package domain

import (
	"errors"
)

var (
	ErrEntityNotFound    = errors.New("entity not found")
	ErrUnauthorized      = errors.New("failed to authorize")
	ErrAlreadyExists     = errors.New("entity already exists")
	ErrAlreadyClosed     = errors.New("receipt already closed")
	ErrUndefinedValue    = errors.New("undefined value")
	ErrInconsistentState = errors.New("inconsistent state")
)
