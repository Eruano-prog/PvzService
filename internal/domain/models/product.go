package models

import (
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
