package jwt

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
)

type Service struct {
	log          *slog.Logger
	secret       []byte
	allowedRoles map[models.UserRole]struct{} // Для валидации ролей
}

func (s Service) GenerateToken(userID uuid.UUID, role models.UserRole) (string, error) {
	// Валидация роли
	if _, ok := s.allowedRoles[role]; !ok {
		return "", fmt.Errorf("invalid role: %s", role)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID.String(),
		"role":   role,
	})

	return token.SignedString(s.secret)
}

func (s Service) VerifyToken(tokenString string) (userID uuid.UUID, role models.UserRole, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		s.log.Debug("Error parsing token", "error", err)
		return uuid.Nil, "", fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		s.log.Error("Invalid token claims")
		return uuid.Nil, "", errors.New("invalid token claims")
	}

	// Извлекаем и парсим userID (ожидаем строку)
	userIDStr, ok := (*claims)["userID"].(string)
	if !ok {
		s.log.Debug("Invalid userID type in token")
		return uuid.Nil, "", errors.New("invalid userID format")
	}

	parsedUUID, err := uuid.Parse(userIDStr)
	if err != nil {
		s.log.Debug("Failed to parse userID as UUID", "userID", userIDStr)
		return uuid.Nil, "", fmt.Errorf("invalid userID: %w", err)
	}

	r, ok := (*claims)["role"].(models.UserRole)
	if !ok {
		s.log.Debug("Invalid role type in token")
		return uuid.Nil, "", errors.New("invalid role format")
	}

	if _, ok := s.allowedRoles[r]; !ok {
		s.log.Debug("Invalid role value in token", "role", r)
		return uuid.Nil, "", fmt.Errorf("unauthorized role: %s", r)
	}

	return parsedUUID, r, nil
}

func NewJWTService(l *slog.Logger, secret string, allowedRoles []models.UserRole) service.TokenService {
	rolesMap := make(map[models.UserRole]struct{})
	for _, role := range allowedRoles {
		rolesMap[role] = struct{}{}
	}

	return &Service{
		log:          l,
		secret:       []byte(secret),
		allowedRoles: rolesMap,
	}
}
