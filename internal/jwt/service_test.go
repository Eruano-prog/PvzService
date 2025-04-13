package jwt_test

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
	jwtService "AvitoPvz/internal/jwt"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndVerifyToken(t *testing.T) {
	log := slog.Default()
	secret := "test-secret"
	expiration := 15 * time.Minute

	service := jwtService.NewJWTService(log, secret, expiration)

	userID := uuid.New()
	role := models.RoleModerator

	tokenInfo := &models.Token{
		UserID:   userID,
		UserRole: role,
	}

	t.Run("successful generation and verification", func(t *testing.T) {
		tokenString, err := service.GenerateToken(tokenInfo)
		require.NoError(t, err)
		require.NotEmpty(t, tokenString)

		verifiedToken, err := service.VerifyToken(tokenString)
		require.NoError(t, err)
		require.NotNil(t, verifiedToken)

		assert.Equal(t, userID, verifiedToken.UserID)
		assert.Equal(t, role, verifiedToken.UserRole)
	})

	// Test token with invalid signature
	t.Run("invalid signature", func(t *testing.T) {
		tokenString, err := service.GenerateToken(tokenInfo)
		require.NoError(t, err)

		// Create another service with different secret
		invalidService := jwtService.NewJWTService(log, "different-secret", expiration)

		_, err = invalidService.VerifyToken(tokenString)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	// Test expired token
	t.Run("expired token", func(t *testing.T) {
		// Create service with very short expiration
		shortExpService := jwtService.NewJWTService(log, secret, time.Nanosecond)
		tokenString, err := shortExpService.GenerateToken(tokenInfo)
		require.NoError(t, err)

		// Wait for token to expire
		time.Sleep(time.Nanosecond)

		_, err = service.VerifyToken(tokenString)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "token is expired")
	})

	// Test invalid token format
	t.Run("invalid token format", func(t *testing.T) {
		_, err := service.VerifyToken("invalid.token.format")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid token")
	})

	// Test missing userID
	t.Run("missing userID", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"role": string(role),
			"exp":  time.Now().Add(expiration).Unix(),
		})
		tokenString, err := token.SignedString([]byte(secret))
		require.NoError(t, err)

		_, err = service.VerifyToken(tokenString)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid userID format")
	})

	// Test invalid userID format
	t.Run("invalid userID format", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": "not-a-uuid",
			"role":   string(role),
			"exp":    time.Now().Add(expiration).Unix(),
		})
		tokenString, err := token.SignedString([]byte(secret))
		require.NoError(t, err)

		_, err = service.VerifyToken(tokenString)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid userID")
	})

	// Test missing role
	t.Run("missing role", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": userID.String(),
			"exp":    time.Now().Add(expiration).Unix(),
		})
		tokenString, err := token.SignedString([]byte(secret))
		require.NoError(t, err)

		_, err = service.VerifyToken(tokenString)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role format")
	})

	// Test invalid role
	t.Run("invalid role", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userID": userID.String(),
			"role":   "invalid-role",
			"exp":    time.Now().Add(expiration).Unix(),
		})
		tokenString, err := token.SignedString([]byte(secret))
		require.NoError(t, err)

		_, err = service.VerifyToken(tokenString)
		require.Error(t, err)
		assert.True(t, errors.Is(err, domain.ErrUndefinedValue))
	})
}

func TestNewJWTService(t *testing.T) {
	log := slog.Default()
	secret := "test-secret"
	expiration := 15 * time.Minute

	t.Run("valid initialization", func(t *testing.T) {
		service := jwtService.NewJWTService(log, secret, expiration)
		require.NotNil(t, service)
	})

	t.Run("empty secret", func(t *testing.T) {
		service := jwtService.NewJWTService(log, "", expiration)
		require.NotNil(t, service)
		// This should still work, but tokens won't be secure
	})

	t.Run("zero expiration", func(t *testing.T) {
		service := jwtService.NewJWTService(log, secret, 0)
		require.NotNil(t, service)
		// Tokens will be generated but will expire immediately
	})
}
