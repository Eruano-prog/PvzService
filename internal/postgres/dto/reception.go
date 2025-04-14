package dto

import (
	"AvitoPvz/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type ReceptionDTO struct {
	ID       uuid.UUID `db:"id"`
	PVZID    uuid.UUID `db:"pvz_id"`
	DateTime time.Time `db:"datetime"`
	Status   string    `db:"status"`
}

func (r *ReceptionDTO) ToModel() (*models.Reception, error) {
	status, err := models.GetStatusFromString(r.Status)
	if err != nil {
		return nil, err
	}

	return &models.Reception{
		ID:       r.ID,
		PVZID:    r.PVZID,
		DateTime: r.DateTime,
		Status:   status,
	}, nil
}
