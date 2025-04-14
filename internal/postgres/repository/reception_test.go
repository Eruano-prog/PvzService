package repository

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/postgres/dto"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestReceptionRepository_InsertReception(t *testing.T) {
	repo, mock, cleanup := newTestReceptionRepo(t)
	defer cleanup()

	ctx := context.Background()
	reception := &models.Reception{
		ID:       uuid.New(),
		PVZID:    uuid.New(),
		DateTime: time.Now(),
		Status:   models.ReceptionStatusInProgress,
	}

	t.Run("successful insertion", func(t *testing.T) {
		dbReception := dto.ReceptionDTO{
			ID:       reception.ID,
			PVZID:    reception.PVZID,
			DateTime: reception.DateTime,
			Status:   string(reception.Status),
		}

		mock.ExpectExec("INSERT INTO receptions").WithArgs(
			dbReception.ID, dbReception.PVZID, dbReception.DateTime, dbReception.Status,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.InsertReception(ctx, reception)
		require.NoError(t, err, "insert should succeed")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("failed insertion", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO receptions").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("database error"))

		err := repo.InsertReception(ctx, reception)
		require.Error(t, err, "insert should fail")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}

func TestReceptionRepository_GetReceptionByPVZIDFilteredByReceptionTime(t *testing.T) {
	repo, mock, cleanup := newTestReceptionRepo(t)
	defer cleanup()

	ctx := context.Background()
	pvzID := uuid.New()
	fromTime := time.Now().Add(-24 * time.Hour)
	toTime := time.Now()

	t.Run("successful query", func(t *testing.T) {
		expectedReception := models.Reception{
			ID:       uuid.New(),
			PVZID:    pvzID,
			DateTime: time.Now(),
			Status:   models.ReceptionStatusInProgress,
		}

		dbReception := dto.ReceptionDTO{
			ID:       expectedReception.ID,
			PVZID:    expectedReception.PVZID,
			DateTime: expectedReception.DateTime,
			Status:   string(expectedReception.Status),
		}

		rows := sqlmock.NewRows([]string{"id", "pvz_id", "datetime", "status"}).
			AddRow(dbReception.ID, dbReception.PVZID, dbReception.DateTime, dbReception.Status)

		mock.ExpectQuery("SELECT .* FROM receptions").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnRows(rows)

		receptions, err := repo.GetReceptionByPVZIDFilteredByReceptionTime(ctx, pvzID, &fromTime, &toTime)
		require.NoError(t, err, "query should succeed")
		require.Len(t, receptions, 1, "should return one reception")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("no receptions found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM receptions").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		receptions, err := repo.GetReceptionByPVZIDFilteredByReceptionTime(ctx, pvzID, &fromTime, &toTime)
		require.NoError(t, err, "no error should be returned")
		require.Len(t, receptions, 0, "should return empty slice")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM receptions").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("database error"))

		receptions, err := repo.GetReceptionByPVZIDFilteredByReceptionTime(ctx, pvzID, &fromTime, &toTime)
		require.Error(t, err, "should return error")
		require.Nil(t, receptions, "receptions should be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}

func TestReceptionRepository_GetActiveReceptionInPVZ(t *testing.T) {
	repo, mock, cleanup := newTestReceptionRepo(t)
	defer cleanup()

	ctx := context.Background()
	pvzID := uuid.New()

	t.Run("active reception found", func(t *testing.T) {
		expectedReception := models.Reception{
			ID:       uuid.New(),
			PVZID:    pvzID,
			DateTime: time.Now(),
			Status:   models.ReceptionStatusInProgress,
		}

		dbReception := dto.ReceptionDTO{
			ID:       expectedReception.ID,
			PVZID:    expectedReception.PVZID,
			DateTime: expectedReception.DateTime,
			Status:   string(expectedReception.Status),
		}

		rows := sqlmock.NewRows([]string{"id", "pvz_id", "datetime", "status"}).
			AddRow(dbReception.ID, dbReception.PVZID, dbReception.DateTime, dbReception.Status)

		mock.ExpectQuery("SELECT .* FROM receptions").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnRows(rows)

		reception, err := repo.GetActiveReceptionInPVZ(ctx, pvzID)
		require.NoError(t, err, "query should succeed")
		require.NotNil(t, reception, "reception should not be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("no active reception", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM receptions").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(sql.ErrNoRows)

		reception, err := repo.GetActiveReceptionInPVZ(ctx, pvzID)
		require.ErrorIs(t, err, sql.ErrNoRows, "should return no rows error")
		require.Nil(t, reception, "reception should be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}

func TestReceptionRepository_ChangeActiveReceptionStatusByPVZID(t *testing.T) {
	repo, mock, cleanup := newTestReceptionRepo(t)
	defer cleanup()

	ctx := context.Background()
	pvzID := uuid.New()
	newStatus := models.ReceptionStatusClosed

	t.Run("status changed successfully", func(t *testing.T) {
		expectedReception := models.Reception{
			ID:       uuid.New(),
			PVZID:    pvzID,
			DateTime: time.Now(),
			Status:   newStatus,
		}

		dbReception := dto.ReceptionDTO{
			ID:       expectedReception.ID,
			PVZID:    expectedReception.PVZID,
			DateTime: expectedReception.DateTime,
			Status:   string(expectedReception.Status),
		}

		rows := sqlmock.NewRows([]string{"id", "pvz_id", "datetime", "status"}).
			AddRow(dbReception.ID, dbReception.PVZID, dbReception.DateTime, dbReception.Status)

		mock.ExpectQuery("UPDATE receptions").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnRows(rows)

		reception, err := repo.ChangeActiveReceptionStatusByPVZID(ctx, pvzID, newStatus)
		require.NoError(t, err, "change should succeed")
		require.NotNil(t, reception, "reception should not be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("no reception to change", func(t *testing.T) {
		mock.ExpectQuery("UPDATE receptions").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		reception, err := repo.ChangeActiveReceptionStatusByPVZID(ctx, pvzID, newStatus)
		require.ErrorIs(t, err, domain.ErrEntityNotFound, "should return ErrEntityNotFound")
		require.Nil(t, reception, "reception should be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}
