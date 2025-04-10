package postgres

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"time"
)

type PVZ struct {
	log *slog.Logger

	db *sqlx.DB
}

func (p PVZ) InsertPVZ(ctx context.Context, pvz *models.PVZ) error {
	dbPVZ := pvzDTO{
		ID:              pvz.ID,
		City:            pvz.City,
		RegitrationTime: pvz.RegistrationDate,
	}

	_, err := p.db.NamedExecContext(ctx, pvzInsertQuery, dbPVZ)
	if err != nil {
		p.log.Error("Error inserting PVZ", "error", err)
		return err
	}

	return nil
}

func (p PVZ) GetPagedPVZsFilteredByReceptionTime(ctx context.Context, fromTime, toTime *time.Time, fromNumber, limit int) ([]models.PVZ, error) {
	params := map[string]interface{}{
		"startTime": fromTime,
		"endTime":   toTime,
		"offset":    fromNumber,
		"limit":     limit,
	}

	q, args, err := p.db.BindNamed(pvzGetFilteredByReceptionTime, params)
	if err != nil {
		p.log.Error("Error binding named params", "error", err)
		return nil, err
	}

	resultDTO := make([]pvzDTO, 0, limit)
	err = p.db.SelectContext(ctx, &resultDTO, q, args...)
	if err != nil {
		p.log.Error("Error getting pvzs", "error", err)
		return nil, err
	}

	result := make([]models.PVZ, len(resultDTO))
	for i, dto := range resultDTO {
		result[i] = models.PVZ{
			ID:               dto.ID,
			City:             dto.City,
			RegistrationDate: dto.RegitrationTime,
		}
	}

	return result, nil
}

func NewPVZ(log *slog.Logger, db *sqlx.DB) service.PVZRepository {
	return &PVZ{
		log: log,
		db:  db,
	}
}
