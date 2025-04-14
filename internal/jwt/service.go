package jwt

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service struct {
	log        *slog.Logger
	secret     []byte
	expiration time.Duration
}

func (s Service) GenerateToken(t *models.Token) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": t.UserID.String(),
		"role":   t.UserRole,
		"exp":    time.Now().Add(s.expiration).Unix(),
	})

	return token.SignedString(s.secret)
}

func (s Service) VerifyToken(tokenString string) (tokenInfo *models.Token, err error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		s.log.Debug("Error parsing token", "error", err)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		s.log.Debug("Invalid token claims or token expired")
		return nil, errors.New("invalid token or expired")
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok {
		s.log.Debug("Invalid userID type in token")
		return nil, errors.New("invalid userID format")
	}

	parsedUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.log.Debug("Failed to parse userID as UUID", "userID", userIDStr)
		return nil, fmt.Errorf("invalid userID: %w", err)
	}

	r, ok := claims["role"].(string)
	if !ok {
		s.log.Debug("Invalid role type in token")
		return nil, errors.New("invalid role format")
	}

	role, err := models.GetRoleFromString(r)
	if err != nil {
		s.log.Debug("Invalid role type in token")
		return nil, err
	}

	return &models.Token{
		UserID:   parsedUUID,
		UserRole: role,
	}, nil
}

func NewJWTService(l *slog.Logger, secret string, expiration time.Duration) service.TokenService {
	return &Service{
		log:        l,
		secret:     []byte(secret),
		expiration: expiration,
	}
}
