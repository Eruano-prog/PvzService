package middleware_test

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest/middleware"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockVerifier struct {
	mock.Mock
}

func (m *MockVerifier) VerifyToken(token string) (*models.Token, error) {
	args := m.Called(token)
	return args.Get(0).(*models.Token), args.Error(1)
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}

func TestAuthMiddleware_Success(t *testing.T) {
	testCases := []struct {
		name         string
		token        string
		userRole     models.UserRole
		allowedRoles []models.UserRole
	}{
		{
			name:         "employee with employee only",
			token:        "valid-token",
			userRole:     models.RoleEmployee,
			allowedRoles: middleware.EmployeeOnly,
		},
		{
			name:         "moderator with moderator only",
			token:        "valid-token",
			userRole:     models.RoleModerator,
			allowedRoles: middleware.ModeratorOnly,
		},
		{
			name:         "employee with employee and moderator",
			token:        "valid-token",
			userRole:     models.RoleEmployee,
			allowedRoles: middleware.EmployeeAndModerator,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			verifier := new(MockVerifier)
			log := slog.Default()
			mockHandler := new(MockHandler)

			// Mock verification
			userID := uuid.New()
			verifier.On("VerifyToken", tc.token).
				Return(&models.Token{
					UserID:   userID,
					UserRole: tc.userRole,
				}, nil)

			// Expect handler to be called
			mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					r := args.Get(1).(*http.Request)
					ctx := r.Context()
					assert.Equal(t, userID, ctx.Value("userID"))
				})

			// Create middleware
			authMiddleware := middleware.AuthMiddleware(mockHandler, verifier, log, tc.allowedRoles)

			// Request
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			w := httptest.NewRecorder()

			// Execute
			authMiddleware.ServeHTTP(w, req)

			// Verify
			assert.Equal(t, http.StatusOK, w.Code)
			verifier.AssertExpectations(t)
			mockHandler.AssertExpectations(t)
		})
	}
}

func TestAuthMiddleware_ErrorCases(t *testing.T) {
	testCases := []struct {
		name           string
		token          string
		mockSetup      func(*MockVerifier)
		allowedRoles   []models.UserRole
		expectedStatus int
	}{
		{
			name:  "no token",
			token: "",
			mockSetup: func(mv *MockVerifier) {
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "token verification failed",
			token: "invalid-token",
			mockSetup: func(mv *MockVerifier) {
				mv.On("VerifyToken", "invalid-token").
					Return((*models.Token)(nil), errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:  "insufficient permissions",
			token: "user-token",
			mockSetup: func(mv *MockVerifier) {
				mv.On("VerifyToken", "user-token").
					Return(&models.Token{
						UserRole: models.RoleModerator,
					}, nil)
			},
			allowedRoles:   middleware.EmployeeOnly,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			verifier := new(MockVerifier)
			log := slog.Default()
			mockHandler := new(MockHandler)

			// Setup mocks
			if tc.mockSetup != nil {
				tc.mockSetup(verifier)
			}

			// Create middleware
			authMiddleware := middleware.AuthMiddleware(mockHandler, verifier, log, tc.allowedRoles)

			// Request
			req := httptest.NewRequest("GET", "/", nil)
			if tc.token != "" {
				if strings.HasPrefix(tc.token, "Bearer ") {
					req.Header.Set("Authorization", tc.token)
				} else {
					req.Header.Set("Authorization", "Bearer "+tc.token)
				}
			}
			w := httptest.NewRecorder()

			// Execute
			authMiddleware.ServeHTTP(w, req)

			// Verify
			assert.Equal(t, tc.expectedStatus, w.Code)
			verifier.AssertExpectations(t)
			mockHandler.AssertNotCalled(t, "ServeHTTP")
		})
	}
}

func TestAuthMiddleware_ContextValues(t *testing.T) {
	// Setup
	verifier := new(MockVerifier)
	log := slog.Default()
	mockHandler := new(MockHandler)

	// Test data
	userID := uuid.New()
	token := "valid-token"
	allowedRoles := middleware.EmployeeOnly

	// Mock verification
	verifier.On("VerifyToken", token).
		Return(&models.Token{
			UserID:   userID,
			UserRole: models.RoleEmployee,
		}, nil)

	// Expect handler to be called with userID in context
	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			ctx := r.Context()

			// Check context values
			assert.Equal(t, userID, ctx.Value("userID"))

			// Check that original context values are preserved
			originalValue := "test-value"
			ctx = context.WithValue(ctx, "testKey", originalValue)
			assert.Equal(t, originalValue, ctx.Value("testKey"))
		})

	// Create middleware
	authMiddleware := middleware.AuthMiddleware(mockHandler, verifier, log, allowedRoles)

	// Request with original context
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	// Add original context value
	ctx := context.WithValue(req.Context(), "testKey", "original-value")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute
	authMiddleware.ServeHTTP(w, req)

	// Verify
	assert.Equal(t, http.StatusOK, w.Code)
	verifier.AssertExpectations(t)
	mockHandler.AssertExpectations(t)
}
