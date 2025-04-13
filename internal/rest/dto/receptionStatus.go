package dto

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
)

func StatusToDTO(status models.ReceptionStatus) (ReceptionStatus, error) {
	switch status {
	case models.ReceptionStatusInProgress:
		return InProgress, nil
	case models.ReceptionStatusClosed:
		return Close, nil
	default:
		return "", domain.ErrUndefinedValue
	}
}
