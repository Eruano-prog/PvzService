package dto

import (
	"AvitoPvz/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type ProductDTO struct {
	ID          uuid.UUID `db:"id"`
	ReceptionID uuid.UUID `db:"reception_id"`
	Type        string    `db:"type"`
	DateTime    time.Time `db:"datetime"`
}

func (p ProductDTO) ToModel() (*models.Product, error) {
	t, err := models.GetProductTypeFromString(p.Type)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:          p.ID,
		ReceptionID: p.ReceptionID,
		Type:        t,
		DateTime:    p.DateTime,
	}, nil
}
