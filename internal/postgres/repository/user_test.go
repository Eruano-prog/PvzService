package repository

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/postgres/dto"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestUserRepository_InsertUser(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	user := &models.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: "password123",
		Role:     models.RoleEmployee,
	}

	t.Run("successful insertion", func(t *testing.T) {
		dbUser := dto.UserDTO{
			ID:       user.ID,
			Email:    user.Email,
			Password: user.Password,
			Role:     string(user.Role),
		}

		mock.ExpectExec("INSERT INTO users").WithArgs(
			dbUser.ID, dbUser.Email, dbUser.Password, dbUser.Role,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.InsertUser(ctx, user)
		require.NoError(t, err, "insert should succeed")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("failed insertion", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("database error"))

		err := repo.InsertUser(ctx, user)
		require.Error(t, err, "insert should fail")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}

func TestUserRepository_FindUserByEmail(t *testing.T) {
	repo, mock, cleanup := newTestUserRepo(t)
	defer cleanup()

	ctx := context.Background()
	email := "test@example.com"
	expectedUser := &models.User{
		ID:       uuid.New(),
		Email:    email,
		Password: "password123",
		Role:     models.RoleModerator,
	}

	t.Run("user found", func(t *testing.T) {
		dbUser := dto.UserDTO{
			ID:       expectedUser.ID,
			Email:    expectedUser.Email,
			Password: expectedUser.Password,
			Role:     string(expectedUser.Role),
		}

		rows := sqlmock.NewRows([]string{"id", "email", "password", "role"}).
			AddRow(dbUser.ID, dbUser.Email, dbUser.Password, dbUser.Role)

		mock.ExpectQuery("SELECT .* FROM users").WithArgs(email).WillReturnRows(rows)

		user, err := repo.FindUserByEmail(ctx, email)
		require.NoError(t, err, "find should succeed")
		require.NotNil(t, user, "user should not be nil")
		require.Equal(t, expectedUser.Email, user.Email, "email should match")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("user not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM users").WithArgs(email).WillReturnError(sql.ErrNoRows)

		user, err := repo.FindUserByEmail(ctx, email)
		require.ErrorIs(t, err, domain.ErrEntityNotFound, "should return ErrEntityNotFound")
		require.Nil(t, user, "user should be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM users").WithArgs(email).WillReturnError(errors.New("database error"))

		user, err := repo.FindUserByEmail(ctx, email)
		require.Error(t, err, "should return error")
		require.Nil(t, user, "user should be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}
