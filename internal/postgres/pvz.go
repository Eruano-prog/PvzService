package postgres

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/service"
	"context"
	"encoding/json"
	"github.com/google/uuid"
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
		ID:               pvz.ID,
		City:             string(pvz.City),
		RegistrationTime: pvz.RegistrationDate,
	}

	_, err := p.db.NamedExecContext(ctx, pvzInsertQuery, dbPVZ)
	if err != nil {
		p.log.Error("Error inserting PVZ", "error", err)
		return err
	}

	return nil
}

func (p PVZ) GetPagedPVZsFilteredByReceptionTime(ctx context.Context, fromTime, toTime *time.Time, fromNumber, limit int) ([]models.PVZWithReceptions, error) {
	params := map[string]interface{}{
		"startTime": fromTime,
		"endTime":   toTime,
		"offset":    fromNumber,
		"limit":     limit,
	}

	type dbResponse struct {
		pvzDTO
		ReceptionsJSON []byte `db:"receptions_with_products"`
	}

	q, args, err := p.db.BindNamed(pvzGetFilteredByReceptionTime, params)
	if err != nil {
		p.log.Error("Error binding named params", "error", err)
		return nil, err
	}

	var dbResults []dbResponse
	err = p.db.SelectContext(ctx, &dbResults, q, args...)
	if err != nil {
		p.log.Error("failed to execute query", "error", err)
		return nil, err
	}

	result := make([]models.PVZWithReceptions, 0, len(dbResults))
	for _, dbItem := range dbResults {
		city, err := models.GetCityFromString(dbItem.City)
		if err != nil {
			p.log.Warn("invalid city in DB", "city", dbItem.City, "error", err)
			continue
		}

		var receptionsWithProducts []struct {
			Reception struct {
				ID       uuid.UUID `json:"id"`
				Status   string    `json:"status"`
				DateTime string    `json:"datetime"`
				PVZID    uuid.UUID `json:"pvz_id"`
			} `json:"reception"`
			Products []struct {
				ID          uuid.UUID `json:"id"`
				Type        string    `json:"type"`
				DateTime    string    `json:"datetime"`
				ReceptionID uuid.UUID `json:"reception_id"`
			} `json:"products"`
		}

		if err := json.Unmarshal(dbItem.ReceptionsJSON, &receptionsWithProducts); err != nil {
			p.log.Error("failed to unmarshal receptions JSON", "error", err)
			continue
		}

		pvzReceptions := make([]models.ReceptionWithProducts, 0, len(receptionsWithProducts))
		for _, rwp := range receptionsWithProducts {
			receptionTime, err := parsePostgreTime(rwp.Reception.DateTime)
			if err != nil {
				p.log.Error("failed to parse reception time",
					"time", rwp.Reception.DateTime, "error", err)
				continue
			}

			status, err := models.GetStatusFromString(rwp.Reception.Status)
			if err != nil {
				p.log.Warn("invalid reception status", "status", rwp.Reception.Status, "error", err)
				continue
			}

			products := make([]models.Product, 0, len(rwp.Products))
			for _, pr := range rwp.Products {
				t, err := models.GetProductTypeFromString(pr.Type)
				if err != nil {
					p.log.Warn("invalid product type", "type", pr.Type, "error", err)
					continue
				}
				productTime, err := parsePostgreTime(pr.DateTime)
				if err != nil {
					p.log.Error("failed to parse product time",
						"time", pr.DateTime, "error", err)
					continue
				}

				products = append(products, models.Product{
					ID:          pr.ID,
					Type:        t,
					DateTime:    productTime,
					ReceptionID: pr.ReceptionID,
				})
			}
			pvzReceptions = append(pvzReceptions, models.ReceptionWithProducts{
				Reception: models.Reception{
					ID:       rwp.Reception.ID,
					Status:   status,
					DateTime: receptionTime,
					PVZID:    rwp.Reception.PVZID,
				},
				Products: products,
			})
		}

		result = append(result, models.PVZWithReceptions{
			PVZ: models.PVZ{
				ID:               dbItem.ID,
				City:             city,
				RegistrationDate: dbItem.RegistrationTime,
			},
			Receptions: pvzReceptions,
		})
	}

	return result, nil
}

func NewPVZRepo(log *slog.Logger, db *sqlx.DB) service.PVZRepository {
	return &PVZ{
		log: log,
		db:  db,
	}
}

func parsePostgreTime(t string) (time.Time, error) {
	const psqlTimeLayout = "2006-01-02T15:04:05.999999"

	return time.Parse(psqlTimeLayout, t)
}
