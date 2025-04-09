package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID
	ReceptionID uuid.UUID
	Type        ProductType
	DateTime    time.Time
}

type ProductType string

const (
	ProductTypeElectronics ProductType = "electronics"
	ProductTypeClothing    ProductType = "clothing"
	ProductTypeShoes       ProductType = "shoes"
)

func GetProductTypeFromString(str string) (ProductType, error) {
	switch str {
	case "electronics":
		return ProductTypeElectronics, nil
	case "clothing":
		return ProductTypeClothing, nil
	case "shoes":
		return ProductTypeShoes, nil
	default:
		return "", errors.New("invalid product type")
	}
}
