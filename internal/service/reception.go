package service

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	log                 *slog.Logger
	receptionRepository ReceptionRepository
	metrics             BusinessMetrics
}

func (r Reception) CreateReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	reception := &models.Reception{
		ID:       uuid.New(),
		PVZID:    pvzID,
		DateTime: time.Now(),
		Status:   models.ReceptionStatusInProgress,
	}

	err := r.receptionRepository.InsertReception(ctx, reception)
	if err != nil {
		r.log.Error("Failed to insert reception", "pvzID", pvzID, "err", err)
		return nil, err
	}

	r.metrics.IncReceptionCreated()

	return reception, nil
}

func NewReceptionService(log *slog.Logger, receptionRepo ReceptionRepository, metrics BusinessMetrics) rest.ReceptionService {
	return &Reception{
		log:                 log,
		receptionRepository: receptionRepo,
		metrics:             metrics,
	}
}
