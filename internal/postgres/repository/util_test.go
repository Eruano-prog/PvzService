package repository

import (
	"github.com/jmoiron/sqlx"
	"testing"

	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log/slog"
)

func newTestProductRepo(t *testing.T) (*Product, sqlmock.Sqlmock, func() error) {
	t.Helper()

	db, mock, err := sqlmock.Newx()
	require.NoError(t, err, "failed to create mock database")

	sqlxDB := sqlx.NewDb(db.DB, "sqlmock")

	repo := NewProductRepo(slog.Default(), sqlxDB).(*Product)
	require.NotNil(t, repo, "repository should not be nil")

	return repo, mock, db.Close
}

func newTestReceptionRepo(t *testing.T) (*Reception, sqlmock.Sqlmock, func() error) {
	t.Helper()

	db, mock, err := sqlmock.Newx()
	require.NoError(t, err, "failed to create mock database")

	sqlxDB := sqlx.NewDb(db.DB, "sqlmock")

	repo := NewReceptionRepo(slog.Default(), sqlxDB).(*Reception)
	require.NotNil(t, repo, "repository should not be nil")

	return repo, mock, db.Close
}

func newTestUserRepo(t *testing.T) (*User, sqlmock.Sqlmock, func() error) {
	t.Helper()

	db, mock, err := sqlmock.Newx()
	require.NoError(t, err, "failed to create mock database")

	sqlxDB := sqlx.NewDb(db.DB, "sqlmock")

	repo := NewUserRepo(slog.Default(), sqlxDB).(*User)
	require.NotNil(t, repo, "repository should not be nil")

	return repo, mock, db.Close
}

func newTestPVZRepo(t *testing.T) (*PVZ, sqlmock.Sqlmock, func() error) {
	t.Helper()

	db, mock, err := sqlmock.Newx()
	require.NoError(t, err, "failed to create mock database")

	sqlxDB := sqlx.NewDb(db.DB, "sqlmock")

	repo := NewPVZRepo(slog.Default(), sqlxDB).(*PVZ)
	require.NotNil(t, repo, "repository should not be nil")

	return repo, mock, db.Close
}

func callCleanupOrLog(t *testing.T, clean func() error) {
	err := clean()
	if err != nil {
		t.Logf("failed to cleanup: %v", err)
	}
}
