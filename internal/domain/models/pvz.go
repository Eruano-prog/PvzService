package models

import (
	"AvitoPvz/internal/domain"
	"github.com/google/uuid"
	"time"
)

type PVZ struct {
	ID               uuid.UUID
	City             City
	RegistrationDate time.Time
}

type PVZWithReceptions struct {
	PVZ        PVZ
	Receptions []ReceptionWithProducts
}

type City string

var (
	CityMSK City = "Moscow"
	CitySPB City = "Saint-Petersburg"
	CityKZN City = "Kazan"
)

func GetCityFromString(str string) (City, error) {
	switch str {
	case "Moscow":
		return CityMSK, nil
	case "Saint-Petersburg":
		return CitySPB, nil
	case "Kazan":
		return CityKZN, nil
	default:
		return City(""), domain.ErrUndefinedValue
	}
}
