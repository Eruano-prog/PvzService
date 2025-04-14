package repository

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/postgres"
	"AvitoPvz/internal/postgres/dto"
	"AvitoPvz/internal/service"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	psqlTimeLayout = "2006-01-02T15:04:05.999999"
)

type PVZ struct {
	log *slog.Logger
	db  *sqlx.DB
}
type dbResponse struct {
	dto.PvzDTO
	ReceptionsJSON []byte `db:"receptions_with_products"`
}

type receptionJSON struct {
	ID       uuid.UUID `json:"id"`
	Status   string    `json:"status"`
	DateTime string    `json:"datetime"`
	PVZID    uuid.UUID `json:"pvz_id"`
}

type productJSON struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	DateTime    string    `json:"datetime"`
	ReceptionID uuid.UUID `json:"reception_id"`
}

type receptionWithProductsJSON struct {
	Reception receptionJSON `json:"reception"`
	Products  []productJSON `json:"products"`
}

func (p PVZ) InsertPVZ(ctx context.Context, pvz *models.PVZ) error {
	dbPVZ := dto.PvzDTO{
		ID:               pvz.ID,
		City:             string(pvz.City),
		RegistrationTime: pvz.RegistrationDate,
	}

	if _, err := p.db.NamedExecContext(ctx, postgres.PvzInsertQuery, dbPVZ); err != nil {
		p.log.Error("Error inserting PVZ", "error", err)
		return err
	}

	return nil
}

func (p PVZ) GetAllPvzs(ctx context.Context) ([]models.PVZ, error) {
	var resultDTO []dto.PvzDTO
	err := p.db.SelectContext(ctx, &resultDTO, postgres.PvzGetAll)
	if err != nil {
		p.log.Error("Error getting all PVZs", "error", err)
		return nil, err
	}

	result := make([]models.PVZ, 0, len(resultDTO))
	for _, v := range resultDTO {
		model, err := v.ToModel()
		if err != nil {
			p.log.Warn("Error converting PVZ to model", "error", err)
			continue
		}

		result = append(result, model)
	}

	return result, nil
}

func (p PVZ) GetPagedPVZsFilteredByReceptionTime(
	ctx context.Context,
	fromTime, toTime *time.Time,
	fromNumber, limit int,
) ([]models.PVZWithReceptions, error) {
	params := map[string]interface{}{
		"startTime": fromTime,
		"endTime":   toTime,
		"offset":    fromNumber,
		"limit":     limit,
	}

	q, args, err := p.db.BindNamed(postgres.PvzGetFilteredByReceptionTime, params)
	if err != nil {
		p.log.Error("Error binding named params", "error", err)
		return nil, err
	}

	var dbResults []dbResponse
	if err := p.db.SelectContext(ctx, &dbResults, q, args...); err != nil {
		p.log.Error("Failed to execute query", "error", err)
		return nil, err
	}

	return p.mapDBResultsToDomain(dbResults)
}

func (p PVZ) mapDBResultsToDomain(dbResults []dbResponse) ([]models.PVZWithReceptions, error) {
	result := make([]models.PVZWithReceptions, 0, len(dbResults))

	for _, dbItem := range dbResults {
		pvzWithReceptions, err := p.mapDBItemToDomain(dbItem)
		if err != nil {
			p.log.Warn("Skipping PVZ due to mapping error", "error", err)
			continue
		}
		result = append(result, *pvzWithReceptions)
	}

	return result, nil
}

func (p PVZ) mapDBItemToDomain(dbItem dbResponse) (*models.PVZWithReceptions, error) {
	city, err := models.GetCityFromString(dbItem.City)
	if err != nil {
		return nil, err
	}

	var receptionsJSON []receptionWithProductsJSON
	if err := json.Unmarshal(dbItem.ReceptionsJSON, &receptionsJSON); err != nil {
		return nil, err
	}

	receptions, err := p.mapReceptionsJSONToDomain(receptionsJSON)
	if err != nil {
		return nil, err
	}

	return &models.PVZWithReceptions{
		PVZ: models.PVZ{
			ID:               dbItem.ID,
			City:             city,
			RegistrationDate: dbItem.RegistrationTime,
		},
		Receptions: receptions,
	}, nil
}

func (p PVZ) mapReceptionsJSONToDomain(receptionsJSON []receptionWithProductsJSON) ([]models.ReceptionWithProducts, error) {
	receptions := make([]models.ReceptionWithProducts, 0, len(receptionsJSON))

	for _, rwp := range receptionsJSON {
		reception, err := p.mapReceptionJSONToDomain(rwp.Reception)
		if err != nil {
			p.log.Warn("Skipping reception due to mapping error", "error", err)
			continue
		}

		products, err := p.mapProductsJSONToDomain(rwp.Products)
		if err != nil {
			p.log.Warn("Skipping reception products due to mapping error", "error", err)
			continue
		}

		receptions = append(receptions, models.ReceptionWithProducts{
			Reception: *reception,
			Products:  products,
		})
	}

	return receptions, nil
}

func (p PVZ) mapReceptionJSONToDomain(rj receptionJSON) (*models.Reception, error) {
	status, err := models.GetStatusFromString(rj.Status)
	if err != nil {
		return nil, err
	}

	dateTime, err := time.Parse(psqlTimeLayout, rj.DateTime)
	if err != nil {
		return nil, err
	}

	return &models.Reception{
		ID:       rj.ID,
		Status:   status,
		DateTime: dateTime,
		PVZID:    rj.PVZID,
	}, nil
}

func (p PVZ) mapProductsJSONToDomain(productsJSON []productJSON) ([]models.Product, error) {
	products := make([]models.Product, 0, len(productsJSON))

	for _, pj := range productsJSON {
		product, err := p.mapProductJSONToDomain(pj)
		if err != nil {
			p.log.Warn("Skipping product due to mapping error", "error", err)
			continue
		}
		products = append(products, *product)
	}

	return products, nil
}

func (p PVZ) mapProductJSONToDomain(pj productJSON) (*models.Product, error) {
	productType, err := models.GetProductTypeFromString(pj.Type)
	if err != nil {
		return nil, err
	}

	dateTime, err := time.Parse(psqlTimeLayout, pj.DateTime)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:          pj.ID,
		Type:        productType,
		DateTime:    dateTime,
		ReceptionID: pj.ReceptionID,
	}, nil
}

func NewPVZRepo(log *slog.Logger, db *sqlx.DB) service.PVZRepository {
	return &PVZ{
		log: log,
		db:  db,
	}
}
