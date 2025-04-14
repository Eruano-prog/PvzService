package service_test

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPVZService_CreatePVZ(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		testPVZ := &models.PVZ{ID: uuid.New(), City: models.CityMSK}
		pvzRepo.On("InsertPVZ", mock.Anything, testPVZ).Return(nil)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		result, err := s.CreatePVZ(context.Background(), testPVZ)

		require.NoError(t, err)
		assert.Equal(t, testPVZ, result)
		pvzRepo.AssertExpectations(t)
	})

	t.Run("error when PVZ repository fails", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		testPVZ := &models.PVZ{ID: uuid.New(), City: models.CityMSK}
		expectedErr := errors.New("database error")
		pvzRepo.On("InsertPVZ", mock.Anything, testPVZ).Return(expectedErr)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		_, err := s.CreatePVZ(context.Background(), testPVZ)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		pvzRepo.AssertExpectations(t)
	})
}

func TestPVZService_GetPVZsWithReceptions(t *testing.T) {
	t.Run("successful get with filters", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		startDate := time.Now().Add(-24 * time.Hour)
		endDate := time.Now()
		expectedPVZs := []models.PVZWithReceptions{
			{
				PVZ: models.PVZ{ID: uuid.New(), City: models.CityMSK},
				Receptions: []models.ReceptionWithProducts{
					{
						Reception: models.Reception{ID: uuid.New(), PVZID: uuid.New(), Status: models.ReceptionStatusInProgress},
						Products:  []models.Product{},
					},
				},
			},
		}

		pvzRepo.On("GetPagedPVZsFilteredByReceptionTime",
			mock.Anything,
			&startDate,
			&endDate,
			0,
			10).
			Return(expectedPVZs, nil)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		result, err := s.GetPVZsWithReceptions(context.Background(), &startDate, &endDate, 1, 10)

		require.NoError(t, err)
		assert.Equal(t, expectedPVZs, result)
		pvzRepo.AssertExpectations(t)
	})

	t.Run("empty result", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		pvzRepo.On("GetPagedPVZsFilteredByReceptionTime",
			mock.Anything,
			(*time.Time)(nil),
			(*time.Time)(nil),
			0,
			10).
			Return([]models.PVZWithReceptions{}, nil)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		result, err := s.GetPVZsWithReceptions(context.Background(), nil, nil, 1, 10)

		require.NoError(t, err)
		assert.Empty(t, result)
		pvzRepo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		expectedErr := errors.New("database error")
		pvzRepo.On("GetPagedPVZsFilteredByReceptionTime",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything).
			Return([]models.PVZWithReceptions(nil), expectedErr)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		_, err := s.GetPVZsWithReceptions(context.Background(), nil, nil, 1, 10)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		pvzRepo.AssertExpectations(t)
	})
}

func TestPVZService_CloseLastReception(t *testing.T) {
	t.Run("successful close", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		pvzID := uuid.New()
		expectedReception := &models.Reception{
			ID:     uuid.New(),
			PVZID:  pvzID,
			Status: models.ReceptionStatusClosed,
		}

		receptionRepo.On("ChangeActiveReceptionStatusByPVZID",
			mock.Anything,
			pvzID,
			models.ReceptionStatusClosed).
			Return(expectedReception, nil)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		result, err := s.CloseLastReception(context.Background(), pvzID)

		require.NoError(t, err)
		assert.Equal(t, expectedReception, result)
		receptionRepo.AssertExpectations(t)
	})

	t.Run("error when no active reception", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		pvzID := uuid.New()
		expectedErr := errors.New("no active reception")

		receptionRepo.On("ChangeActiveReceptionStatusByPVZID",
			mock.Anything,
			pvzID,
			models.ReceptionStatusClosed).
			Return((*models.Reception)(nil), expectedErr)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		_, err := s.CloseLastReception(context.Background(), pvzID)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		receptionRepo.AssertExpectations(t)
	})
}

func TestPVZService_DeleteLastProduct(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		pvzID := uuid.New()
		productRepo.On("DeleteLastProductByPVZID", mock.Anything, pvzID).Return(nil)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		err := s.DeleteLastProduct(context.Background(), pvzID)

		require.NoError(t, err)
		productRepo.AssertExpectations(t)
	})

	t.Run("error when product not found", func(t *testing.T) {
		pvzRepo := new(MockPVZRepository)
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		pvzID := uuid.New()
		expectedErr := errors.New("product not found")
		productRepo.On("DeleteLastProductByPVZID", mock.Anything, pvzID).Return(expectedErr)

		s := service.NewPVZService(log, pvzRepo, receptionRepo, productRepo)

		err := s.DeleteLastProduct(context.Background(), pvzID)

		require.Error(t, err)
		assert.Equal(t, expectedErr, err)
		productRepo.AssertExpectations(t)
	})
}
