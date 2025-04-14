package service

import (
	"AvitoPvz/internal/domain/models"
	"context"
	"time"

	"github.com/google/uuid"
)

type TokenService interface {
	GenerateToken(token *models.Token) (string, error)
	VerifyToken(tokenString string) (token *models.Token, err error)
}

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type PVZRepository interface {
	InsertPVZ(ctx context.Context, pvz *models.PVZ) error

	GetPagedPVZsFilteredByReceptionTime(ctx context.Context, fromTime, toTime *time.Time, fromNumber, limit int) ([]models.PVZWithReceptions, error)
	GetAllPvzs(ctx context.Context) ([]models.PVZ, error)
}

type ReceptionRepository interface {
	InsertReception(ctx context.Context, reception *models.Reception) error
	GetReceptionByPVZIDFilteredByReceptionTime(ctx context.Context, pvzID uuid.UUID, from, to *time.Time) ([]models.Reception, error)
	GetActiveReceptionInPVZ(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error)

	ChangeActiveReceptionStatusByPVZID(ctx context.Context, pvzID uuid.UUID, newStatus models.ReceptionStatus) (*models.Reception, error)
}

type ProductRepository interface {
	InsertProductIfReceptionNotClosed(ctx context.Context, product *models.Product) error
	GetProductsByReceiptID(ctx context.Context, receiptID uuid.UUID) ([]models.Product, error)

	DeleteLastProductByPVZID(ctx context.Context, pvzID uuid.UUID) error
}

type BusinessMetrics interface {
	IncPVZCreated()
	IncReceptionCreated()
	IncProductAdded()
}
