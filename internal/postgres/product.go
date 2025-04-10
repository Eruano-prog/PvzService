package postgres

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
	"log/slog"
)

type Product struct {
	log *slog.Logger
	db  *sqlx.DB
}

func (p Product) InsertProductIfReceptionNotClosed(ctx context.Context, product *models.Product) error {
	params := map[string]interface{}{
		"product_id":       product.ID,
		"reception_id":     product.ReceptionID,
		"reception_status": string(models.ReceptionStatusInProgress),
		"type":             string(product.Type),
		"datetime":         product.DateTime,
	}

	q, args, err := p.db.BindNamed(productInsertQuery, params)
	if err != nil {
		p.log.Error("Error binding params tto query")
		return err
	}

	var insertedID uuid.UUID
	err = p.db.QueryRowxContext(ctx, q, args...).Scan(&insertedID)
	if errors.Is(err, pgx.ErrNoRows) {
		p.log.Info("Вставки не произошло: неверный статус или ID приёмки")
		return domain.ErrAlreadyClosed
	} else if err != nil {
		p.log.Error("Ошибка запроса:", err)
		return err
	} else {
		p.log.Debug("Успешно вставлено, ID:", insertedID)
	}

	return nil
}

func (p Product) GetProductsByReceiptID(ctx context.Context, receiptID uuid.UUID) ([]models.Product, error) {
	params := map[string]interface{}{
		"reception_id": receiptID,
	}

	var productDTOs []productDTO
	err := p.db.SelectContext(ctx, &productDTOs, productFindByReceptionIDQuery, params)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		p.log.Error("Error selecting products by reception id query")
		return nil, err
	}

	result := make([]models.Product, 0, len(productDTOs))
	for _, dto := range productDTOs {
		product, err := dto.toModel()
		if err != nil {
			return nil, err
		}
		result = append(result, *product)
	}

	return result, nil
}

func (p Product) DeleteLastProductByPVZID(ctx context.Context, pvzID uuid.UUID) error {
	params := map[string]interface{}{
		"pvz_id": pvzID,
		"status": string(models.ReceptionStatusInProgress),
	}

	var deletedID uuid.UUID
	err := p.db.QueryRowxContext(ctx, productDeleteByPvzIDQuery, params).Scan(&deletedID)

	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrEntityNotFound
	default:
		return err
	}
}

func NewProductRepo(log *slog.Logger, db *sqlx.DB) service.ProductRepository {
	return &Product{
		log: log,
		db:  db,
	}
}
