package dto

import (
	"AvitoPvz/internal/domain/models"

	"github.com/google/uuid"
)

type UserDTO struct {
	ID       uuid.UUID `db:"id"`
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Role     string    `db:"role"`
}

func (u UserDTO) ToModel() (*models.User, error) {
	role, err := models.GetRoleFromString(u.Role)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:       u.ID,
		Email:    u.Email,
		Password: u.Password,
		Role:     role,
	}, nil
}
