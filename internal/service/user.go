package service

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type User struct {
	log *slog.Logger

	token          TokenService
	userRepository UserRepository
}

func (u User) Register(ctx context.Context, email, password string, role models.UserRole) (user *models.User, err error) {
	_, err = u.userRepository.FindUserByEmail(ctx, email)
	if err == nil || !errors.Is(err, domain.ErrEntityNotFound) {
		u.log.Debug("user already exists or error happened")
		return nil, domain.ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error("failed to hash password", "Error", err)
		return nil, err
	}

	userToInsert := &models.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	err = u.userRepository.InsertUser(ctx, userToInsert)
	if err != nil {
		u.log.Error("failed to insert user", "Error", err)
		return nil, err
	}

	return userToInsert, nil
}

func (u User) Login(ctx context.Context, email, password string) (token string, err error) {
	user, err := u.userRepository.FindUserByEmail(ctx, email)

	if err != nil {
		u.log.Debug("failed to find user by email", "Error", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		u.log.Debug("failed to compare password", "Error", err)
		return "", domain.ErrUnauthorized
	}

	token, err = u.token.GenerateToken(&models.Token{UserID: user.ID, UserRole: user.Role})
	if err != nil {
		u.log.Error("failed to generate token", "Error", err)
		return "", err
	}
	return token, nil
}

func (u User) DummyLogin(ctx context.Context, role models.UserRole) (token string, err error) {
	token, err = u.token.GenerateToken(&models.Token{UserID: uuid.Nil, UserRole: role})
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewUserService(log *slog.Logger, tokenService TokenService, userRepository UserRepository) rest.UserService {
	return &User{
		log:            log,
		token:          tokenService,
		userRepository: userRepository,
	}
}
