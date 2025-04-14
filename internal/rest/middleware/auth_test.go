package middleware_test

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"AvitoPvz/internal/rest/middleware"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testKey = rest.RequestContextKey("testKey")

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
			verifier := new(MockVerifier)
			log := slog.Default()
			mockHandler := new(MockHandler)

			userID := uuid.New()
			verifier.On("VerifyToken", tc.token).
				Return(&models.Token{
					UserID:   userID,
					UserRole: tc.userRole,
				}, nil)

			mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
				Run(func(args mock.Arguments) {
					r := args.Get(1).(*http.Request)
					ctx := r.Context()
					assert.Equal(t, userID, ctx.Value(rest.IdKey))
				})

			authMiddleware := middleware.Auth(mockHandler, verifier, log, tc.allowedRoles)

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer "+tc.token)
			w := httptest.NewRecorder()

			authMiddleware.ServeHTTP(w, req)

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
			verifier := new(MockVerifier)
			log := slog.Default()
			mockHandler := new(MockHandler)

			if tc.mockSetup != nil {
				tc.mockSetup(verifier)
			}
			authMiddleware := middleware.Auth(mockHandler, verifier, log, tc.allowedRoles)

			req := httptest.NewRequest("GET", "/", nil)
			if tc.token != "" {
				if strings.HasPrefix(tc.token, "Bearer ") {
					req.Header.Set("Authorization", tc.token)
				} else {
					req.Header.Set("Authorization", "Bearer "+tc.token)
				}
			}
			w := httptest.NewRecorder()

			authMiddleware.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			verifier.AssertExpectations(t)
			mockHandler.AssertNotCalled(t, "ServeHTTP")
		})
	}
}

func TestAuthMiddleware_ContextValues(t *testing.T) {
	verifier := new(MockVerifier)
	log := slog.Default()
	mockHandler := new(MockHandler)

	userID := uuid.New()
	token := "valid-token"
	allowedRoles := middleware.EmployeeOnly

	verifier.On("VerifyToken", token).
		Return(&models.Token{
			UserID:   userID,
			UserRole: models.RoleEmployee,
		}, nil)

	mockHandler.On("ServeHTTP", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			r := args.Get(1).(*http.Request)
			ctx := r.Context()

			assert.Equal(t, userID, ctx.Value("userID"))

			originalValue := "test-value"
			ctx = context.WithValue(ctx, testKey, originalValue)
			assert.Equal(t, originalValue, ctx.Value(testKey))
		})

	authMiddleware := middleware.Auth(mockHandler, verifier, log, allowedRoles)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	ctx := context.WithValue(req.Context(), testKey, "original-value")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	authMiddleware.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	verifier.AssertExpectations(t)
	mockHandler.AssertExpectations(t)
}
