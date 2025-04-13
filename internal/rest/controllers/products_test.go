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

func TestProductController_AddProductHandler(t *testing.T) {
	t.Run("successful product creation", func(t *testing.T) {
		productService := new(MockProductService)
		log := slog.Default()
		controller := NewProductController(log, productService)
		mockVerifier := new(MockVerifier)

		pvzID := uuid.New()
		productType := models.ProductTypeElectronics
		now := time.Now()
		expectedProduct := &models.Product{
			ID:          uuid.New(),
			ReceptionID: uuid.New(),
			Type:        productType,
			DateTime:    now,
		}

		// Настройка моков
		productService.On("AddProduct", mock.Anything, productType, pvzID).
			Return(expectedProduct, nil)
		mockVerifier.On("VerifyToken", "valid-token").
			Return(&models.Token{UserRole: models.RoleEmployee}, nil)

		// Подготовка запроса
		requestBody := dto.PostProductsJSONRequestBody{
			PvzId: pvzID,
			Type:  dto.PostProductsJSONBodyType(dto.ProductTypeЭлектроника),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")

		w := httptest.NewRecorder()

		// Регистрация роута и вызов
		mux := http.NewServeMux()
		controller.Register(mux, mockVerifier)
		mux.ServeHTTP(w, req)

		// Проверки
		require.Equal(t, http.StatusCreated, w.Code)

		var response dto.Product
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, expectedProduct.ID, *response.Id)
		assert.Equal(t, expectedProduct.ReceptionID, response.ReceptionId)
		assert.Equal(t, dto.ProductTypeЭлектроника, response.Type)
		assert.WithinDuration(t, expectedProduct.DateTime, *response.DateTime, time.Second)

		productService.AssertExpectations(t)
		mockVerifier.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		productService := new(MockProductService)
		log := slog.Default()
		controller := NewProductController(log, productService)
		mockVerifier := new(MockVerifier)

		mockVerifier.On("VerifyToken", "valid-token").
			Return(&models.Token{UserRole: models.RoleEmployee}, nil)

		// Подготовка запроса с невалидным JSON
		req := httptest.NewRequest("POST", "/products", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")

		w := httptest.NewRecorder()

		// Регистрация роута и вызов
		mux := http.NewServeMux()
		controller.Register(mux, mockVerifier)
		mux.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		productService.AssertNotCalled(t, "AddProduct")
	})

	t.Run("invalid product type", func(t *testing.T) {
		productService := new(MockProductService)
		log := slog.Default()
		controller := NewProductController(log, productService)
		mockVerifier := new(MockVerifier)

		mockVerifier.On("VerifyToken", "valid-token").
			Return(&models.Token{UserRole: models.RoleEmployee}, nil)

		// Подготовка запроса с невалидным типом продукта
		requestBody := dto.PostProductsJSONRequestBody{
			PvzId: uuid.New(),
			Type:  "invalid-type",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")

		w := httptest.NewRecorder()

		// Регистрация роута и вызов
		mux := http.NewServeMux()
		controller.Register(mux, mockVerifier)
		mux.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		productService.AssertNotCalled(t, "AddProduct")
	})

	t.Run("service error", func(t *testing.T) {
		productService := new(MockProductService)
		log := slog.Default()
		controller := NewProductController(log, productService)
		mockVerifier := new(MockVerifier)

		mockVerifier.On("VerifyToken", "valid-token").
			Return(&models.Token{UserRole: models.RoleEmployee}, nil)

		pvzID := uuid.New()
		productType := models.ProductTypeElectronics
		expectedErr := "service error"

		// Настройка моков
		productService.On("AddProduct", mock.Anything, productType, pvzID).
			Return((*models.Product)(nil), errors.New(expectedErr))
		mockVerifier.On("VerifyToken", "valid-token").
			Return(&models.Token{UserRole: models.RoleEmployee}, nil)

		// Подготовка запроса
		requestBody := dto.PostProductsJSONRequestBody{
			PvzId: pvzID,
			Type:  dto.PostProductsJSONBodyType(dto.ProductTypeЭлектроника),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/products", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer valid-token")

		w := httptest.NewRecorder()

		// Регистрация роута и вызов
		mux := http.NewServeMux()
		controller.Register(mux, mockVerifier)
		mux.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)

		var errorResponse dto.Error
		err := json.NewDecoder(w.Body).Decode(&errorResponse)
		require.NoError(t, err)
		assert.Contains(t, errorResponse.Message, expectedErr)

		productService.AssertExpectations(t)
		mockVerifier.AssertExpectations(t)
	})
}
