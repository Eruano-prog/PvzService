package dto

import (
	"AvitoPvz/internal/domain/models"
	"github.com/oapi-codegen/runtime/types"
)

func UserToDTO(user models.User) (*User, error) {
	role, err := RoleToDTO(user.Role)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:    &user.ID,
		Email: types.Email(user.Email),
		Role:  role,
	}, nil
}
