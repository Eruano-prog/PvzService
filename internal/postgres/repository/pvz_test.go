package repository

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/postgres/dto"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
)

func TestPVZRepository_InsertPVZ(t *testing.T) {
	repo, mock, cleanup := newTestPVZRepo(t)
	defer callCleanupOrLog(t, cleanup)

	ctx := context.Background()
	pvz := &models.PVZ{
		ID:               uuid.New(),
		City:             models.CityMSK,
		RegistrationDate: time.Now(),
	}

	t.Run("successful insertion", func(t *testing.T) {
		dbPVZ := dto.PvzDTO{
			ID:               pvz.ID,
			City:             string(pvz.City),
			RegistrationTime: pvz.RegistrationDate,
		}

		mock.ExpectExec("INSERT INTO pvzs").WithArgs(
			dbPVZ.ID,
			dbPVZ.City,
			dbPVZ.RegistrationTime,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.InsertPVZ(ctx, pvz)
		require.NoError(t, err, "insert should succeed")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("failed insertion", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO pvzs").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("database error"))

		err := repo.InsertPVZ(ctx, pvz)
		require.Error(t, err, "insert should fail")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}

func TestPVZRepository_GetPagedPVZsFilteredByReceptionTime(t *testing.T) {
	repo, mock, cleanup := newTestPVZRepo(t)
	defer callCleanupOrLog(t, cleanup)

	ctx := context.Background()
	fromTime := time.Now().Add(-24 * time.Hour)
	toTime := time.Now()
	fromNumber := 0
	limit := 10

	t.Run("successful query", func(t *testing.T) {
		pvzID := uuid.New()
		dbResponse := dbResponse{
			PvzDTO: dto.PvzDTO{
				ID:               pvzID,
				City:             string(models.CityMSK),
				RegistrationTime: time.Now(),
			},
			ReceptionsJSON: prepareMockReceptionsJSON(t),
		}

		rows := sqlmock.NewRows([]string{"id", "city", "registration_time", "receptions_with_products"}).
			AddRow(dbResponse.ID, dbResponse.City, dbResponse.RegistrationTime, dbResponse.ReceptionsJSON)

		mock.ExpectQuery("SELECT .* FROM pvzs").WithArgs(
			fromTime,
			toTime,
			fromNumber,
			limit,
		).WillReturnRows(rows)

		pvzs, err := repo.GetPagedPVZsFilteredByReceptionTime(ctx, &fromTime, &toTime, fromNumber, limit)
		require.NoError(t, err, "query should succeed")
		require.Len(t, pvzs, 1, "should return one PVZ")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("no PVZs found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM pvzs").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		pvzs, err := repo.GetPagedPVZsFilteredByReceptionTime(ctx, &fromTime, &toTime, fromNumber, limit)
		require.NoError(t, err, "no error should be returned")
		require.Len(t, pvzs, 0, "should return empty slice")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM pvzs").WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(errors.New("database error"))

		pvzs, err := repo.GetPagedPVZsFilteredByReceptionTime(ctx, &fromTime, &toTime, fromNumber, limit)
		require.Error(t, err, "should return error")
		require.Nil(t, pvzs, "PVZs should be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}

func prepareMockReceptionsJSON(t *testing.T) []byte {
	receptionID := uuid.New()
	productID := uuid.New()

	reception := receptionJSON{
		ID:       receptionID,
		Status:   string(models.ReceptionStatusInProgress),
		DateTime: time.Now().Format(psqlTimeLayout),
		PVZID:    uuid.New(),
	}

	product := productJSON{
		ID:          productID,
		Type:        string(models.ProductTypeElectronics),
		DateTime:    time.Now().Format(psqlTimeLayout),
		ReceptionID: receptionID,
	}

	receptionWithProducts := receptionWithProductsJSON{
		Reception: reception,
		Products:  []productJSON{product},
	}

	jsonData, err := json.Marshal([]receptionWithProductsJSON{receptionWithProducts})
	require.NoError(t, err, "failed to marshal mock JSON")

	return jsonData
}
