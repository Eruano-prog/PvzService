package dto

import (
	"AvitoPvz/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

func (p PVZ) ToModel() (*models.PVZ, error) {
	var id uuid.UUID
	if p.Id == nil {
		id = uuid.Nil
	} else {
		id = *p.Id
	}

	city, err := CityToModel(p.City)
	if err != nil {
		return nil, err
	}

	var t time.Time
	if p.RegistrationDate == nil {
		t = time.Time{}
	} else {
		t = *p.RegistrationDate
	}

	return &models.PVZ{
		ID:               id,
		City:             city,
		RegistrationDate: t,
	}, nil
}

func PVZToDTO(pvz models.PVZ) (*PVZ, error) {
	city, err := ModelToCity(pvz.City)
	if err != nil {
		return nil, err
	}

	return &PVZ{
		Id:               &pvz.ID,
		City:             city,
		RegistrationDate: &pvz.RegistrationDate,
	}, nil
}
