package dto

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
)

func RoleToModel(role UserRole) (models.UserRole, error) {
	switch role {
	case UserRoleEmployee:
		return models.RoleEmployee, nil
	case UserRoleModerator:
		return models.RoleModerator, nil
	default:
		return "", domain.ErrUndefinedValue
	}
}

func RoleToDTO(role models.UserRole) (UserRole, error) {
	switch role {
	case models.RoleEmployee:
		return UserRoleEmployee, nil
	case models.RoleModerator:
		return UserRoleModerator, nil
	default:
		return "", domain.ErrUndefinedValue
	}
}
