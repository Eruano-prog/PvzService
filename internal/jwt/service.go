package jwt

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Service struct {
	log        *slog.Logger
	secret     []byte
	expiration time.Duration
}

func (s Service) GenerateToken(userID uuid.UUID, role models.UserRole) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID.String(),
		"role":   role,
		"exp":    time.Now().Add(s.expiration).Unix(),
	})

	return token.SignedString(s.secret)
}

func (s Service) VerifyToken(tokenString string) (userID uuid.UUID, role models.UserRole, err error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		s.log.Debug("Error parsing token", "error", err)
		return uuid.Nil, "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		s.log.Debug("Invalid token claims or token expired")
		return uuid.Nil, "", errors.New("invalid token or expired")
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok {
		s.log.Debug("Invalid userID type in token")
		return uuid.Nil, "", errors.New("invalid userID format")
	}

	parsedUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.log.Debug("Failed to parse userID as UUID", "userID", userIDStr)
		return uuid.Nil, "", fmt.Errorf("invalid userID: %w", err)
	}

	r, ok := claims["role"].(string)
	if !ok {
		s.log.Debug("Invalid role type in token")
		return uuid.Nil, "", errors.New("invalid role format")
	}

	role, err = models.GetRoleFromString(r)
	if err != nil {
		s.log.Debug("Invalid role type in token")
		return uuid.Nil, "", err
	}

	return parsedUUID, role, nil
}

func NewJWTService(l *slog.Logger, secret string, expiration time.Duration) service.TokenService {
	l.Debug("NewJWTService", expiration)
	return &Service{
		log:        l,
		secret:     []byte(secret),
		expiration: expiration,
	}
}
