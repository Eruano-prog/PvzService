package middleware

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"context"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"slices"
	"strings"
)

var (
	EmployeeOnly         = []models.UserRole{models.RoleEmployee}
	ModeratorOnly        = []models.UserRole{models.RoleModerator}
	EmployeeAndModerator = []models.UserRole{models.RoleEmployee, models.RoleModerator}
)

type Verifier interface {
	VerifyToken(tokenString string) (userID uuid.UUID, role models.UserRole, err error)
}

func AuthMiddleware(next http.Handler, verifier Verifier, log *slog.Logger, allowedRoles []models.UserRole) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		id, role, err := verifier.VerifyToken(token)
		if err != nil {
			log.Debug("Failed to verify token", "err", err)
			rest.WriteError(w, log, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		if !slices.Contains(allowedRoles, role) {
			log.Debug("role not in allowed", "role:", role, "allowed:", allowedRoles)
			rest.WriteError(w, log, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "userID", id)))
	})
}
