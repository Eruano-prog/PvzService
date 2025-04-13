package dto

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
)

func TypeToModel(t ProductType) (models.ProductType, error) {
	switch t {
	case ProductTypeЭлектроника:
		return models.ProductTypeElectronics, nil
	case ProductTypeОдежда:
		return models.ProductTypeClothing, nil
	case ProductTypeОбувь:
		return models.ProductTypeShoes, nil
	default:
		return "", domain.ErrUndefinedValue
	}
}

func TypeToDTO(t models.ProductType) (ProductType, error) {
	switch t {
	case models.ProductTypeElectronics:
		return ProductTypeЭлектроника, nil
	case models.ProductTypeClothing:
		return ProductTypeОдежда, nil
	case models.ProductTypeShoes:
		return ProductTypeОбувь, nil
	default:
		return "", domain.ErrUndefinedValue
	}
}
