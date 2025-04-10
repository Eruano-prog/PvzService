package models

import (
	"AvitoPvz/internal/domain"
	"github.com/google/uuid"
	"time"
)

type Reception struct {
	ID       uuid.UUID
	PVZID    uuid.UUID
	DateTime time.Time
	Status   ReceptionStatus
}

type ReceptionStatus string

const (
	ReceptionStatusInProgress ReceptionStatus = "in_progress"
	ReceptionStatusClosed     ReceptionStatus = "closed"
)

type ReceptionWithProducts struct {
	Reception Reception
	Products  []Product
}

func GetStatusFromString(status string) (ReceptionStatus, error) {
	switch status {
	case "in_progress":
		return ReceptionStatusInProgress, nil
	case "closed":
		return ReceptionStatusClosed, nil
	default:
		return ReceptionStatus(""), domain.ErrUndefinedValue
	}
}
