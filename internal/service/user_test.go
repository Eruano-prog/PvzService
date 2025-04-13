package service_test

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"context"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		repo := new(MockUserRepository)
		tokenService := new(MockTokenService)
		log := slog.Default()

		// Настройка моков
		repo.On("FindUserByEmail", mock.Anything, "test@example.com").
			Return((*models.User)(nil), domain.ErrEntityNotFound)

		repo.On("InsertUser", mock.Anything, mock.AnythingOfType("*models.User")).
			Return(nil)

		s := service.NewUserService(log, tokenService, repo)

		user, err := s.Register(context.Background(), "test@example.com", "password", models.RoleModerator)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, user.ID)
		assert.Equal(t, "test@example.com", user.Email)

		repo.AssertExpectations(t)
	})
}

func TestUserService_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		repo := new(MockUserRepository)
		tokenService := new(MockTokenService)
		log := slog.Default()

		userID := uuid.New()
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		user := &models.User{
			ID:       userID,
			Email:    "test@example.com",
			Password: string(hashedPassword),
			Role:     models.RoleModerator,
		}

		// Настройка моков
		repo.On("FindUserByEmail", mock.Anything, "test@example.com").
			Return(user, nil)

		tokenService.On("GenerateToken", mock.AnythingOfType("*models.Token")).
			Return("generated-token", nil)

		s := service.NewUserService(log, tokenService, repo)

		token, err := s.Login(context.Background(), "test@example.com", "password")
		require.NoError(t, err)
		assert.Equal(t, "generated-token", token)

		repo.AssertExpectations(t)
		tokenService.AssertExpectations(t)
	})
}
