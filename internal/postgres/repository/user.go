package repository

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/postgres"
	"AvitoPvz/internal/postgres/dto"
	"AvitoPvz/internal/service"
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type User struct {
	log *slog.Logger

	db *sqlx.DB
}

func (u User) InsertUser(ctx context.Context, user *models.User) error {
	dbUser := dto.UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Role:     string(user.Role),
	}

	_, err := u.db.NamedExecContext(ctx, postgres.UserInsertQuery, dbUser)
	if err != nil {
		u.log.Error("Failed to insert user", "error", err)
		return err
	}

	return nil
}

func (u User) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var dbUser dto.UserDTO

	err := u.db.GetContext(ctx, &dbUser, postgres.UserFindByEmailQuery, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrEntityNotFound
		}
		u.log.Error("Failed to find user by email", "error", err)
		return nil, err
	}

	user, err := dbUser.ToModel()
	if err != nil {
		u.log.Error("Failed to transform DTO to model", "error", err)
		return nil, err
	}

	return user, nil
}

func NewUserRepo(log *slog.Logger, db *sqlx.DB) service.UserRepository {
	return &User{
		log: log,
		db:  db,
	}
}
