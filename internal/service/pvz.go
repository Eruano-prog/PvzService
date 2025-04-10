package service

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"context"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type PVZ struct {
	log *slog.Logger

	pvzRepository       PVZRepository
	receptionRepository ReceptionRepository
	productRepository   ProductRepository
}

func (p PVZ) CreatePVZ(ctx context.Context, city string) (*models.PVZ, error) {
	pvz := &models.PVZ{
		ID:               uuid.New(),
		City:             city,
		RegistrationDate: time.Now(),
	}

	err := p.pvzRepository.InsertPVZ(ctx, pvz)
	if err != nil {
		p.log.Error("Error inserting pvz", "err", err)
		return nil, err
	}

	return pvz, nil
}

// TODO: подумать о многопоточке здесь
// TODO: Solve N+M problem
func (p PVZ) GetPVZsWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	firstElem := (page - 1) * limit

	pvzs, err := p.pvzRepository.GetPagedPVZsFilteredByReceptionTime(ctx, startDate, endDate, firstElem, limit)
	if err != nil {
		p.log.Error("Error getting pvzs", "err", err)
		return nil, err
	}

	results := make([]models.PVZWithReceptions, 0, len(pvzs))
	for _, pvz := range pvzs {
		receptions, err := p.receptionRepository.GetReceptionByPVZIDFilteredByReceptionTime(ctx, pvz.ID, startDate, endDate)
		if err != nil {
			p.log.Error("Error getting receptions", "err", err)
			continue
		}

		receptionsWithProducts := make([]models.ReceptionWithProducts, 0, len(receptions))
		for _, reception := range receptions {
			products, err := p.productRepository.GetProductsByReceiptID(ctx, reception.ID)
			if err != nil {
				p.log.Error("Error getting products", "err", err)
				continue
			}
			receptionsWithProducts = append(receptionsWithProducts, models.ReceptionWithProducts{
				Reception: reception,
				Products:  products,
			})
		}

		results = append(results, models.PVZWithReceptions{
			PVZ:        pvz,
			Receptions: receptionsWithProducts,
		})
	}

	return results, nil
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
