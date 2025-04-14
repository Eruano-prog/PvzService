package rest

import (
	"AvitoPvz/internal/domain/models"
	"context"
	"time"

	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, email, password string, role models.UserRole) (user *models.User, err error)
	Login(ctx context.Context, email, password string) (token string, err error)
	DummyLogin(ctx context.Context, role models.UserRole) (token string, err error)
}

type PVZService interface {
	CreatePVZ(ctx context.Context, pvz *models.PVZ) (*models.PVZ, error)
	GetPVZsWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error)
	CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)
	DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error
}

type ReceptionService interface {
	CreateReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)
}

type ProductService interface {
	AddProduct(ctx context.Context, productType models.ProductType, pvzID uuid.UUID) (*models.Product, error)
}
