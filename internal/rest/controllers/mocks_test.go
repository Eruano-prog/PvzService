package controllers

import (
	"AvitoPvz/internal/domain/models"
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, email, password string, role models.UserRole) (*models.User, error) {
	args := m.Called(ctx, email, password, role)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) DummyLogin(ctx context.Context, role models.UserRole) (string, error) {
	args := m.Called(ctx, role)
	return args.String(0), args.Error(1)
}

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) AddProduct(ctx context.Context, productType models.ProductType, pvzID uuid.UUID) (*models.Product, error) {
	args := m.Called(ctx, productType, pvzID)
	return args.Get(0).(*models.Product), args.Error(1)
}

type MockVerifier struct {
	mock.Mock
}

func (m *MockVerifier) VerifyToken(token string) (*models.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*models.Token), args.Error(1)
}

type MockReceptionService struct {
	mock.Mock
}

func (m *MockReceptionService) CreateReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(*models.Reception), args.Error(1)
}

type MockPVZService struct {
	mock.Mock
}

func (m *MockPVZService) CreatePVZ(ctx context.Context, pvz *models.PVZ) (*models.PVZ, error) {
	args := m.Called(ctx, pvz)
	return args.Get(0).(*models.PVZ), args.Error(1)
}

func (m *MockPVZService) GetPVZsWithReceptions(ctx context.Context, startDate, endDate *time.Time, page, limit int) ([]models.PVZWithReceptions, error) {
	args := m.Called(ctx, startDate, endDate, page, limit)
	return args.Get(0).([]models.PVZWithReceptions), args.Error(1)
}

func (m *MockPVZService) CloseLastReception(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	args := m.Called(ctx, pvzID)
	return args.Get(0).(*models.Reception), args.Error(1)
}

func (m *MockPVZService) DeleteLastProduct(ctx context.Context, pvzID uuid.UUID) error {
	args := m.Called(ctx, pvzID)
	return args.Error(0)
}
