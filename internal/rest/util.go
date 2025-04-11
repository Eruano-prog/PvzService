package rest

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
