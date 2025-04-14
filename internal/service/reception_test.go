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

func TestReceptionService_CreateReception(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		repo := new(MockReceptionRepository)
		log := slog.Default()

		pvzID := uuid.New()
		repo.On("InsertReception", mock.Anything, mock.AnythingOfType("*models.Reception")).
			Return(nil).
			Run(func(args mock.Arguments) {
				reception := args.Get(1).(*models.Reception)
				assert.Equal(t, pvzID, reception.PVZID)
				assert.Equal(t, models.ReceptionStatusInProgress, reception.Status)
				assert.WithinDuration(t, time.Now(), reception.DateTime, time.Second)
			})

		s := service.NewReceptionService(log, repo)

		result, err := s.CreateReception(context.Background(), pvzID)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, pvzID, result.PVZID)
		assert.Equal(t, models.ReceptionStatusInProgress, result.Status)
		repo.AssertExpectations(t)
	})

	t.Run("error when repository fails", func(t *testing.T) {
		repo := new(MockReceptionRepository)
		log := slog.Default()

		pvzID := uuid.New()
		expectedErr := errors.New("database error")
		repo.On("InsertReception", mock.Anything, mock.AnythingOfType("*models.Reception")).
			Return(expectedErr)

		s := service.NewReceptionService(log, repo)

		result, err := s.CreateReception(context.Background(), pvzID)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
		repo.AssertExpectations(t)
	})
}
