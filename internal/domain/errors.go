package domain

import "fmt"

var (
	ErrUserNotFound  = fmt.Errorf("user not found")
	ErrUnauthorized  = fmt.Errorf("failed to unauthorize")
	ErrAlreadyExists = fmt.Errorf("user already exists")
	ErrAlreadyClosed = fmt.Errorf("recept already closed")
)
