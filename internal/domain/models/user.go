package models

import (
	"fmt"
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Email    string
	Password string
	Role     UserRole
}

type UserRole string

const (
	RoleEmployee  UserRole = "employee"
	RoleModerator UserRole = "moderator"
)

func GetRoleFromString(role string) (UserRole, error) {
	switch role {
	case "employee":
		return RoleEmployee, nil
	case "moderator":
		return RoleModerator, nil
	default:
		return "", fmt.Errorf("invalid role: %s", role)
	}
}
