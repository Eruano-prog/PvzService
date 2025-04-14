package grpc

import (
	"AvitoPvz/internal/domain/models"
	"context"
)

type PvzService interface {
	GetAllPvz(ctx context.Context) ([]models.PVZ, error)
}
