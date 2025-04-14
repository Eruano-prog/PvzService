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

func TestProductRepository_InsertProductIfReceptionNotClosed(t *testing.T) {
	repo, mock, cleanup := newTestProductRepo(t)
	defer cleanup()

	ctx := context.Background()
	product := &models.Product{
		ID:          uuid.New(),
		ReceptionID: uuid.New(),
		Type:        models.ProductTypeElectronics,
		DateTime:    time.Now(),
	}

	t.Run("successful insertion", func(t *testing.T) {
		// Формируем параметры, как это делает метод
		rows := sqlmock.NewRows([]string{"id"}).AddRow(product.ID)
		mock.ExpectQuery("INSERT INTO products").WithArgs( // Уточняем начало запроса для точности
			product.ReceptionID,                      // reception_id (для подзапроса)
			string(models.ReceptionStatusInProgress), // reception_status (для подзапроса)
			product.ID,                               // product_id (для INSERT)
			product.ReceptionID,                      // reception_id (для INSERT, снова)
			string(product.Type),                     // type
			product.DateTime,                         // datetime
		).WillReturnRows(rows)

		err := repo.InsertProductIfReceptionNotClosed(ctx, product)
		require.NoError(t, err, "insert should succeed")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("reception already closed", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO products").WithArgs(
			product.ReceptionID,                      // reception_id
			string(models.ReceptionStatusInProgress), // reception_status
			product.ID,                               // product_id
			product.ReceptionID,                      // reception_id
			string(product.Type),                     // type
			product.DateTime,                         // datetime
		).WillReturnError(sql.ErrNoRows)

		err := repo.InsertProductIfReceptionNotClosed(ctx, product)
		require.ErrorIs(t, err, domain.ErrAlreadyClosed, "should return ErrAlreadyClosed")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("INSERT INTO products").WithArgs(
			product.ReceptionID,                      // reception_id
			string(models.ReceptionStatusInProgress), // reception_status
			product.ID,                               // product_id
			product.ReceptionID,                      // reception_id
			string(product.Type),                     // type
			product.DateTime,                         // datetime
		).WillReturnError(errors.New("database error"))

		err := repo.InsertProductIfReceptionNotClosed(ctx, product)
		require.Error(t, err, "should return error")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}

func TestProductRepository_GetProductsByReceiptID(t *testing.T) {
	repo, mock, cleanup := newTestProductRepo(t)
	defer cleanup()

	ctx := context.Background()
	receptionID := uuid.New()

	t.Run("successful query", func(t *testing.T) {
		expectedProduct := models.Product{
			ID:          uuid.New(),
			ReceptionID: receptionID,
			Type:        models.ProductTypeElectronics,
			DateTime:    time.Now(),
		}

		dbProduct := dto.ProductDTO{
			ID:          expectedProduct.ID,
			ReceptionID: expectedProduct.ReceptionID,
			Type:        string(expectedProduct.Type),
			DateTime:    expectedProduct.DateTime,
		}

		rows := sqlmock.NewRows([]string{"id", "reception_id", "type", "datetime"}).
			AddRow(dbProduct.ID, dbProduct.ReceptionID, dbProduct.Type, dbProduct.DateTime)

		mock.ExpectQuery("SELECT .* FROM products").WithArgs(receptionID).WillReturnRows(rows)

		products, err := repo.GetProductsByReceiptID(ctx, receptionID)
		require.NoError(t, err, "query should succeed")
		require.Len(t, products, 1, "should return one product")
		require.Equal(t, expectedProduct.ID, products[0].ID, "product ID should match")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("no products found", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM products").WithArgs(receptionID).WillReturnRows(sqlmock.NewRows([]string{"id"}))

		products, err := repo.GetProductsByReceiptID(ctx, receptionID)
		require.NoError(t, err, "no error should be returned")
		require.Len(t, products, 0, "should return empty slice")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT .* FROM products").WithArgs(receptionID).WillReturnError(errors.New("database error"))

		products, err := repo.GetProductsByReceiptID(ctx, receptionID)
		require.Error(t, err, "should return error")
		require.Nil(t, products, "products should be nil")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}

func TestProductRepository_DeleteLastProductByPVZID(t *testing.T) {
	repo, mock, cleanup := newTestProductRepo(t)
	defer cleanup()

	ctx := context.Background()
	pvzID := uuid.New()

	t.Run("successful deletion", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(uuid.New())
		mock.ExpectQuery("DELETE FROM products").WithArgs(
			pvzID,                                    // pvz_id
			string(models.ReceptionStatusInProgress), // status
		).WillReturnRows(rows)

		err := repo.DeleteLastProductByPVZID(ctx, pvzID)
		require.NoError(t, err, "delete should succeed")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("no product to delete", func(t *testing.T) {
		mock.ExpectQuery("DELETE FROM products").WithArgs(
			pvzID,                                    // pvz_id
			string(models.ReceptionStatusInProgress), // status
		).WillReturnError(sql.ErrNoRows)

		err := repo.DeleteLastProductByPVZID(ctx, pvzID)
		require.ErrorIs(t, err, domain.ErrEntityNotFound, "should return ErrEntityNotFound")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})

	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("DELETE FROM products").WithArgs(
			pvzID,                                    // pvz_id
			string(models.ReceptionStatusInProgress), // status
		).WillReturnError(errors.New("database error"))

		err := repo.DeleteLastProductByPVZID(ctx, pvzID)
		require.Error(t, err, "should return error")
		require.NoError(t, mock.ExpectationsWereMet(), "expectations should be met")
	})
}
