package service_test

import (
	"AvitoPvz/internal/domain/models"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(t *models.Token) (string, error) {
	args := m.Called(t)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) VerifyToken(token string) (*models.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*models.Token), args.Error(1)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) InsertUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockPVZRepository struct {
	mock.Mock
}

func (m *MockPVZRepository) InsertPVZ(ctx context.Context, pvz *models.PVZ) error {
	args := m.Called(ctx, pvz)
	return args.Error(0)
}

func (m *MockPVZRepository) GetPagedPVZsFilteredByReceptionTime(ctx context.Context, fromTime, toTime *time.Time, fromNumber, limit int) ([]models.PVZWithReceptions, error) {
	args := m.Called(ctx, fromTime, toTime, fromNumber, limit)
	return args.Get(0).([]models.PVZWithReceptions), args.Error(1)
}

type MockReceptionRepository struct {
	mock.Mock
}

func (m *MockReceptionRepository) InsertReception(ctx context.Context, reception *models.Reception) error {
	args := m.Called(ctx, reception)
	return args.Error(0)
}

func (m *MockReceptionRepository) GetReceptionByPVZIDFilteredByReceptionTime(ctx context.Context, pvzID uuid.UUID, from, to *time.Time) ([]models.Reception, error) {
	args := m.Called(ctx, pvzID, from, to)
	return args.Get(0).([]models.Reception), args.Error(1)
}

func (m *MockReceptionRepository) GetActiveReceptionInPVZ(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *MockReceptionRepository) ChangeActiveReceptionStatusByPVZID(ctx context.Context, pvzID uuid.UUID, newStatus models.ReceptionStatus) (*models.Reception, error) {
	args := m.Called(ctx, pvzID, newStatus)
	return args.Get(0).(*models.Reception), args.Error(1)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) InsertProductIfReceptionNotClosed(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetProductsByReceiptID(ctx context.Context, receiptID uuid.UUID) ([]models.Product, error) {
	args := m.Called(ctx, receiptID)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) DeleteLastProductByPVZID(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}
