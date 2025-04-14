package controllers

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest/dto"
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestReceptionController_CreateReceptionHandler(t *testing.T) {
	t.Run("successful reception creation", func(t *testing.T) {
		receptionService := new(MockReceptionService)
		log := slog.Default()
		controller := NewReceptionController(log, receptionService)
		mockVerifier := new(MockVerifier)

		pvzID := uuid.New()
		now := time.Now()
		expectedReception := &models.Reception{
			ID:       uuid.New(),
			PVZID:    pvzID,
			DateTime: now,
			Status:   models.ReceptionStatusInProgress,
		}

		receptionService.On("CreateReception", mock.Anything, pvzID).
			Return(expectedReception, nil)
		mockVerifier.On("VerifyToken", "valid-token").
			Return(&models.Token{UserRole: models.RoleEmployee}, nil)

		requestBody := dto.PostReceptionsJSONRequestBody{
			PvzId: pvzID,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/receptions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")

		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		controller.Register(mux, mockVerifier)
		mux.ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var response dto.Reception
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, expectedReception.ID, *response.Id)
		assert.Equal(t, expectedReception.PVZID, response.PvzId)
		assert.Equal(t, dto.ReceptionStatus(expectedReception.Status), response.Status)
		assert.WithinDuration(t, expectedReception.DateTime, response.DateTime, time.Second)

		receptionService.AssertExpectations(t)
		mockVerifier.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		receptionService := new(MockReceptionService)
		log := slog.Default()
		controller := NewReceptionController(log, receptionService)
		mockVerifier := new(MockVerifier)

		mockVerifier.On("VerifyToken", "valid-token").
			Return(&models.Token{UserRole: models.RoleEmployee}, nil)

		req := httptest.NewRequest("POST", "/receptions", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")

		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		controller.Register(mux, mockVerifier)
		mux.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		receptionService.AssertNotCalled(t, "CreateReception")
	})

	t.Run("service error", func(t *testing.T) {
		receptionService := new(MockReceptionService)
		log := slog.Default()
		controller := NewReceptionController(log, receptionService)
		mockVerifier := new(MockVerifier)

		pvzID := uuid.New()
		expectedErr := "service error"

		receptionService.On("CreateReception", mock.Anything, pvzID).
			Return((*models.Reception)(nil), errors.New(expectedErr))
		mockVerifier.On("VerifyToken", "valid-token").
			Return(&models.Token{UserRole: models.RoleEmployee}, nil)
		requestBody := dto.PostReceptionsJSONRequestBody{
			PvzId: pvzID,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/receptions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")

		w := httptest.NewRecorder()

		mux := http.NewServeMux()
		controller.Register(mux, mockVerifier)
		mux.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse dto.Error
		err := json.NewDecoder(w.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Message, expectedErr)

		receptionService.AssertExpectations(t)
		mockVerifier.AssertExpectations(t)
	})
}

func TestReceptionController_RegisterRoutes(t *testing.T) {
	receptionService := new(MockReceptionService)
	log := slog.Default()
	controller := NewReceptionController(log, receptionService)
	mockVerifier := new(MockVerifier)

	mux := http.NewServeMux()
	controller.Register(mux, mockVerifier)

	req := httptest.NewRequest("POST", "/receptions", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	// Проверяем, что роут зарегистрирован и возвращает не 404
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}
