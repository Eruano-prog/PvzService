package dto

import "AvitoPvz/internal/domain/models"

func ReceptionToDTO(reception models.Reception) (*Reception, error) {
	status, err := StatusToDTO(reception.Status)
	if err != nil {
		return nil, err
	}

	return &Reception{
		Id:       &reception.ID,
		PvzId:    reception.PVZID,
		DateTime: reception.DateTime,
		Status:   status,
	}, nil
}
