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

	"github.com/google/uuid"
	"github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthController_DummyLoginHandler(t *testing.T) {
	t.Run("successful dummy login", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		role := models.RoleEmployee
		expectedToken := "dummy-token"
		userService.On("DummyLogin", mock.Anything, role).Return(expectedToken, nil)

		requestBody := dto.PostDummyLoginJSONRequestBody{
			Role: dto.PostDummyLoginJSONBodyRole(role),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.dummyLoginHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var response dto.Token
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, expectedToken, response)

		userService.AssertExpectations(t)
	})

	t.Run("invalid role", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		requestBody := dto.PostDummyLoginJSONRequestBody{
			Role: "invalid-role",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.dummyLoginHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		userService.AssertNotCalled(t, "DummyLogin")
	})

	t.Run("service error", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		role := models.RoleEmployee
		expectedErr := "service error"
		userService.On("DummyLogin", mock.Anything, role).Return("", errors.New(expectedErr))

		requestBody := dto.PostDummyLoginJSONRequestBody{
			Role: dto.PostDummyLoginJSONBodyRole(role),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/dummyLogin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.dummyLoginHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		userService.AssertExpectations(t)
	})
}

func TestAuthController_RegisterHandler(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		email := "test@example.com"
		password := "password"
		role := models.RoleEmployee
		userID := uuid.New()

		expectedUser := &models.User{
			ID:    userID,
			Email: email,
			Role:  role,
		}

		userService.On("Register", mock.Anything, email, password, role).
			Return(expectedUser, nil)

		requestBody := dto.PostRegisterJSONBody{
			Email:    types.Email(email),
			Password: password,
			Role:     dto.PostRegisterJSONBodyRole(role),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.registerHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusCreated, w.Code)

		var response dto.User
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, email, string(response.Email))
		assert.Equal(t, userID, *response.Id)
		assert.Equal(t, dto.UserRole(role), response.Role)

		userService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.registerHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		userService.AssertNotCalled(t, "Register")
	})

	t.Run("invalid role", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		requestBody := dto.PostRegisterJSONBody{
			Email:    "test@example.com",
			Password: "password",
			Role:     "invalid-role",
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.registerHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		userService.AssertNotCalled(t, "Register")
	})

	t.Run("service error", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		email := "test@example.com"
		password := "password"
		role := models.RoleEmployee
		expectedErr := "service error"

		userService.On("Register", mock.Anything, email, password, role).
			Return((*models.User)(nil), errors.New(expectedErr))

		requestBody := dto.PostRegisterJSONBody{
			Email:    types.Email(email),
			Password: password,
			Role:     dto.PostRegisterJSONBodyRole(role),
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.registerHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		userService.AssertExpectations(t)
	})
}

func TestAuthController_LoginHandler(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		email := "test@example.com"
		password := "password"
		expectedToken := "auth-token"

		userService.On("Login", mock.Anything, email, password).
			Return(expectedToken, nil)

		requestBody := dto.PostLoginJSONBody{
			Email:    types.Email(email),
			Password: password,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.loginHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		var response dto.Token
		err := json.NewDecoder(w.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, expectedToken, response)

		userService.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.loginHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		userService.AssertNotCalled(t, "Login")
	})

	t.Run("invalid credentials", func(t *testing.T) {
		userService := new(MockUserService)
		log := slog.Default()
		controller := NewAuthController(log, userService)

		email := "test@example.com"
		password := "wrong-password"

		userService.On("Login", mock.Anything, email, password).
			Return("", errors.New("invalid credentials"))

		requestBody := dto.PostLoginJSONBody{
			Email:    types.Email(email),
			Password: password,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		http.HandlerFunc(controller.loginHandler).ServeHTTP(w, req)

		require.Equal(t, http.StatusUnauthorized, w.Code)
		userService.AssertExpectations(t)
	})
}

func TestAuthController_RegisterRoutes(t *testing.T) {
	userService := new(MockUserService)
	log := slog.Default()
	controller := NewAuthController(log, userService)

	mux := http.NewServeMux()
	controller.Register(mux)

	testCases := []struct {
		method string
		path   string
	}{
		{"POST", "/dummyLogin"},
		{"POST", "/register"},
		{"POST", "/login"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			assert.NotEqual(t, http.StatusNotFound, w.Code)
		})
	}
}
