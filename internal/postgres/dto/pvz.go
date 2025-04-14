package dto

import (
	"time"

	"github.com/google/uuid"
)

type PvzDTO struct {
	ID               uuid.UUID `db:"id"`
	City             string    `db:"city"`
	RegistrationTime time.Time `db:"registration_time"`
}
