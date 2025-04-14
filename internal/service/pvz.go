package service

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	log *slog.Logger

	pvzRepository       PVZRepository
	receptionRepository ReceptionRepository
	productRepository   ProductRepository
}

func (p PVZ) CreatePVZ(ctx context.Context, pvz *models.PVZ) (*models.PVZ, error) {
	if pvz.ID == uuid.Nil {
		pvz.ID = uuid.New()
	}

	if pvz.RegistrationDate.IsZero() {
		pvz.RegistrationDate = time.Now()
	}

	err := p.pvzRepository.InsertPVZ(ctx, pvz)
	if err != nil {
		p.log.Error("Error inserting pvz", "err", err)
		return nil, err
	}

	return pvz, nil
}

func (p PVZ) GetPVZsWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	firstElem := min((page-1)*limit, 0)

	pvzs, err := p.pvzRepository.GetPagedPVZsFilteredByReceptionTime(ctx, startDate, endDate, firstElem, limit)
	if err != nil {
		p.log.Error("Error getting pvzs", "err", err)
		return nil, err
	}

	return pvzs, nil
}

func (p PVZ) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	reception, err := p.receptionRepository.ChangeActiveReceptionStatusByPVZID(ctx, pvzID, models.ReceptionStatusClosed)
	if err != nil {
		p.log.Debug("Error closing last reception", "err", err)
		return nil, err
	}

	return reception, nil
}

func (p PVZ) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	err := p.productRepository.DeleteLastProductByPVZID(ctx, pvzID)
	if err != nil {
		p.log.Error("Error deleting last product", "err", err)
		return err
	}

	return nil
}

func NewPVZService(log *slog.Logger, pvzRepository PVZRepository, receptionRepository ReceptionRepository, productRepository ProductRepository) rest.PVZService {
	return &PVZ{
		log:                 log,
		pvzRepository:       pvzRepository,
		receptionRepository: receptionRepository,
		productRepository:   productRepository,
	}
}
