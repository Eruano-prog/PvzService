package dto

import (
	"github.com/google/uuid"
	"time"
)

type PvzDTO struct {
	ID               uuid.UUID `db:"id"`
	City             string    `db:"city"`
	RegistrationTime time.Time `db:"registration_time"`
}
