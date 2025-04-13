package models

import (
	"github.com/google/uuid"
)

type Token struct {
	UserID   uuid.UUID
	UserRole UserRole
}
