package dto

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
)

func CityToModel(city PVZCity) (models.City, error) {
	switch city {
	case Москва:
		return models.CityMSK, nil
	case СанктПетербург:
		return models.CitySPB, nil
	case Казань:
		return models.CityKZN, nil
	default:
		return "", domain.ErrUndefinedValue
	}
}

func ModelToCity(city models.City) (PVZCity, error) {
	switch city {
	case models.CityMSK:
		return Москва, nil
	case models.CitySPB:
		return СанктПетербург, nil
	case models.CityKZN:
		return Казань, nil
	default:
		return "", domain.ErrUndefinedValue
	}
}
