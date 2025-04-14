package dto

import (
	"AvitoPvz/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type PvzDTO struct {
	ID               uuid.UUID `db:"id"`
	City             string    `db:"city"`
	RegistrationTime time.Time `db:"registration_time"`
}

func (p PvzDTO) ToModel() (models.PVZ, error) {
	city, err := models.GetCityFromString(p.City)
	if err != nil {
		return models.PVZ{}, err
	}

	return models.PVZ{
		ID:               p.ID,
		City:             city,
		RegistrationDate: p.RegistrationTime,
	}, nil
}
