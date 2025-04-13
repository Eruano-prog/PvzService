package controllers

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest/dto"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPVZController_CreatePVZHandler(t *testing.T) {
	t.Run("successful PVZ creation", func(t *testing.T) {
		pvzService := new(MockPVZService)
		log := slog.Default()
		controller := NewPVZController(log, pvzService)

		pvzID := uuid.New()
		regDate := time.Now()
		expectedPVZ := &models.PVZ{
			ID:               pvzID,
			City:             models.CityMSK,
			RegistrationDate: regDate,
		}

		pvzService.On("CreatePVZ", mock.Anything, mock.AnythingOfType("*models.PVZ")).
			Return(expectedPVZ, nil)

		requestBody := dto.PVZ{
			Id:               &pvzID,
			City:             dto.Москва,
			RegistrationDate: &regDate,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/pvz", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		controller.createPVZHandler(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var response dto.PVZ
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, pvzID, *response.Id)
		assert.Equal(t, dto.Москва, response.City)
		assert.WithinDuration(t, regDate, *response.RegistrationDate, time.Second)

		pvzService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		pvzService := new(MockPVZService)
		log := slog.Default()
		controller := NewPVZController(log, pvzService)

		req := httptest.NewRequest("POST", "/pvz", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		controller.createPVZHandler(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		pvzService.AssertNotCalled(t, "CreatePVZ")
	})

	t.Run("invalid city", func(t *testing.T) {
		pvzService := new(MockPVZService)
		log := slog.Default()
		controller := NewPVZController(log, pvzService)

		requestBody := map[string]interface{}{
			"city": "invalid-city",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/pvz", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		controller.createPVZHandler(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		pvzService.AssertNotCalled(t, "CreatePVZ")
	})

	t.Run("service error", func(t *testing.T) {
		pvzService := new(MockPVZService)
		log := slog.Default()
		controller := NewPVZController(log, pvzService)

		expectedErr := "service error"
		pvzService.On("CreatePVZ", mock.Anything, mock.Anything).
			Return((*models.PVZ)(nil), errors.New(expectedErr))

		requestBody := dto.PVZ{
			City: dto.Москва,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/pvz", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		controller.createPVZHandler(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse dto.Error
		err := json.NewDecoder(w.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Message, expectedErr)

		pvzService.AssertExpectations(t)
	})
}
