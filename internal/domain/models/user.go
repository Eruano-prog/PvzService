package models

import "github.com/google/uuid"

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
