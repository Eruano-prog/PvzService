package service_test

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"testing"
)

func TestProductService_AddProduct(t *testing.T) {
	t.Run("successful add", func(t *testing.T) {
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		pvzID := uuid.New()
		reception := &models.Reception{ID: uuid.New(), PVZID: pvzID, Status: models.ReceptionStatusInProgress}
		receptionRepo.On("GetActiveReceptionInPVZ", mock.Anything, pvzID).Return(reception, nil)
		productRepo.On("InsertProductIfReceptionNotClosed", mock.Anything, mock.AnythingOfType("*models.Product")).Return(nil)

		s := service.NewProductService(log, productRepo, receptionRepo)

		result, err := s.AddProduct(context.Background(), models.ProductTypeElectronics, pvzID)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, reception.ID, result.ReceptionID)
		receptionRepo.AssertExpectations(t)
		productRepo.AssertExpectations(t)
	})

	t.Run("error when no active reception", func(t *testing.T) {
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		pvzID := uuid.New()
		expectedErr := errors.New("no active reception")
		receptionRepo.On("GetActiveReceptionInPVZ", mock.Anything, pvzID).Return((*models.Reception)(nil), expectedErr)

		s := service.NewProductService(log, productRepo, receptionRepo)

		_, err := s.AddProduct(context.Background(), models.ProductTypeElectronics, pvzID)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "no active reception")
		receptionRepo.AssertExpectations(t)
		productRepo.AssertNotCalled(t, "InsertProductIfReceptionNotClosed")
	})

	t.Run("error when product insertion fails", func(t *testing.T) {
		receptionRepo := new(MockReceptionRepository)
		productRepo := new(MockProductRepository)
		log := slog.Default()

		pvzID := uuid.New()
		reception := &models.Reception{ID: uuid.New(), PVZID: pvzID, Status: models.ReceptionStatusInProgress}
		insertErr := errors.New("database error")
		receptionRepo.On("GetActiveReceptionInPVZ", mock.Anything, pvzID).Return(reception, nil)
		productRepo.On("InsertProductIfReceptionNotClosed", mock.Anything, mock.AnythingOfType("*models.Product")).Return(insertErr)

		s := service.NewProductService(log, productRepo, receptionRepo)

		_, err := s.AddProduct(context.Background(), models.ProductTypeElectronics, pvzID)

		require.Error(t, err)
		assert.Equal(t, insertErr, err)
		receptionRepo.AssertExpectations(t)
		productRepo.AssertExpectations(t)
	})
}
