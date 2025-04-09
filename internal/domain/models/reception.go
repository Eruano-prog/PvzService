package models

import (
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
