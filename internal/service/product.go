package service

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	log *slog.Logger

	receptionRepository ReceptionRepository
	productRepository   ProductRepository

	metrics BusinessMetrics
}

func (p Product) AddProduct(ctx context.Context, productType models.ProductType, pvzID uuid.UUID) (*models.Product, error) {
	reception, err := p.receptionRepository.GetActiveReceptionInPVZ(ctx, pvzID)
	if err != nil {
		p.log.Info("No active reception in pvz")
		return nil, err
	}

	product := &models.Product{
		ID:          uuid.New(),
		ReceptionID: reception.ID,
		Type:        productType,
		DateTime:    time.Now(),
	}

	err = p.productRepository.InsertProductIfReceptionNotClosed(ctx, product)
	if err != nil {
		p.log.Info("Failed to insert product")
		return nil, err
	}

	p.metrics.IncProductAdded()

	return product, nil
}

func NewProductService(log *slog.Logger, productRepository ProductRepository, receptionRepository ReceptionRepository, metrics BusinessMetrics) rest.ProductService {
	return &Product{
		log:                 log,
		productRepository:   productRepository,
		receptionRepository: receptionRepository,
		metrics:             metrics,
	}
}
