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
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jmoiron/sqlx"
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

	q, args, err := p.db.BindNamed(postgres.ProductInsertQuery, params)
	if err != nil {
		p.log.Error("Error binding params to query")
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

func (p Product) GetProductsByReceiptID(ctx context.Context, receptionID uuid.UUID) ([]models.Product, error) {
	var productDTOs []dto.ProductDTO
	err := p.db.SelectContext(ctx, &productDTOs, postgres.ProductFindByReceptionIDQuery, receptionID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		p.log.Error("Error selecting products by reception id query")
		return nil, err
	}

	result := make([]models.Product, 0, len(productDTOs))
	for _, productDTO := range productDTOs {
		product, err := productDTO.ToModel()
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

	q, args, err := p.db.BindNamed(postgres.ProductDeleteByPvzIDQuery, params)
	if err != nil {
		p.log.Error("Error binding params to query", "params", params, "err", err)
		return err
	}

	var deletedID uuid.UUID
	err = p.db.GetContext(ctx, &deletedID, q, args...)

	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrNoRows):
		p.log.Debug("No product was deleted")
		return domain.ErrEntityNotFound
	default:
		p.log.Error("Error executing delete query", "args", args, "err", err)
		return err
	}
}

func NewProductRepo(log *slog.Logger, db *sqlx.DB) service.ProductRepository {
	return &Product{
		log: log,
		db:  db,
	}
}
