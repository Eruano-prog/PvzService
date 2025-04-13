package dto

import "AvitoPvz/internal/domain/models"

func ProductToDTO(product models.Product) (*Product, error) {
	t, err := TypeToDTO(product.Type)
	if err != nil {
		return nil, err
	}

	return &Product{
		Id:          &product.ID,
		ReceptionId: product.ReceptionID,
		DateTime:    &product.DateTime,
		Type:        t,
	}, nil
}
