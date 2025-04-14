package repository

import (
	"AvitoPvz/internal/domain"
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/postgres"
	"AvitoPvz/internal/postgres/dto"
	"AvitoPvz/internal/service"
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Reception struct {
	log *slog.Logger

	db *sqlx.DB
}

func (r Reception) InsertReception(ctx context.Context, reception *models.Reception) error {
	dbReception := dto.ReceptionDTO{
		ID:       reception.ID,
		PVZID:    reception.PVZID,
		DateTime: reception.DateTime,
		Status:   string(reception.Status),
	}

	_, err := r.db.NamedExecContext(ctx, postgres.ReceptionInsertQuery, dbReception)
	if err != nil {
		r.log.Error("Error inserting reception", "error", err)
		return err
	}

	return nil
}

func (r Reception) GetReceptionByPVZIDFilteredByReceptionTime(ctx context.Context, pvzID uuid.UUID, from, to *time.Time) ([]models.Reception, error) {
	params := map[string]interface{}{
		"pvz_id":    pvzID,
		"from_time": from,
		"to_time":   to,
	}

	q, args, err := r.db.BindNamed(postgres.ReceptionFindByPvzAndTimeQuery, params)
	if err != nil {
		r.log.Error("Error binding reception get query", "error", err)
		return nil, err
	}

	var resultsDTO []dto.ReceptionDTO
	err = r.db.SelectContext(ctx, &resultsDTO, q, args...)
	if err != nil {
		r.log.Error("Error binding reception get query", "error", err)
		return nil, err
	}

	results := make([]models.Reception, 0, len(resultsDTO))
	for _, resultDTO := range resultsDTO {
		status, err := models.GetStatusFromString(resultDTO.Status)
		if err != nil {
			r.log.Error("Error getting reception status", "error", err)
			continue
		}

		results = append(results, models.Reception{
			ID:       resultDTO.ID,
			PVZID:    resultDTO.PVZID,
			DateTime: resultDTO.DateTime,
			Status:   status,
		})
	}

	return results, nil
}

func (r Reception) GetActiveReceptionInPVZ(ctx context.Context, pvzID uuid.UUID) (*models.Reception, error) {
	params := map[string]interface{}{
		"pvz_id": pvzID,
		"status": string(models.ReceptionStatusInProgress),
	}
	q, args, err := r.db.BindNamed(postgres.ReceptionFindByStatusPvzQuery, params)
	if err != nil {
		r.log.Error("Error binding reception get active query", "error", err)
		return nil, err
	}

	var resultDTO dto.ReceptionDTO
	err = r.db.GetContext(ctx, &resultDTO, q, args...)
	if err != nil {
		r.log.Error("Error binding reception get active query", "error", err)
		return nil, err
	}

	status, err := models.GetStatusFromString(resultDTO.Status)
	if err != nil {
		r.log.Error("Error getting reception status", "error", err)
		return nil, err
	}

	return &models.Reception{
		ID:       resultDTO.ID,
		PVZID:    resultDTO.PVZID,
		DateTime: resultDTO.DateTime,
		Status:   status,
	}, nil

}

func (r Reception) ChangeActiveReceptionStatusByPVZID(ctx context.Context, pvzID uuid.UUID, newStatus models.ReceptionStatus) (*models.Reception, error) {
	params := map[string]interface{}{
		"pvz_id":      pvzID,
		"status":      string(newStatus),
		"last_status": string(models.ReceptionStatusInProgress),
	}

	q, args, err := r.db.BindNamed(postgres.ReceptionUpdateByStatusAndPvzQuery, params)
	if err != nil {
		r.log.Error("Error binding reception change active query", "error", err)
		return nil, err
	}

	var resultDTO []dto.ReceptionDTO
	err = r.db.SelectContext(ctx, &resultDTO, q, args...)
	if err != nil {
		r.log.Error("Error binding reception change active query", "error", err)
		return nil, err
	}

	if len(resultDTO) == 0 {
		return nil, domain.ErrEntityNotFound
	} else if len(resultDTO) > 1 {
		// Expected that only one non-closed receipt exists per pvz
		return nil, domain.ErrInconsistentState
	}

	rec, err := resultDTO[0].ToModel()
	if err != nil {
		r.log.Error("Error getting reception status", "error", err)
		return nil, err
	}

	return rec, nil
}

func NewReceptionRepo(log *slog.Logger, db *sqlx.DB) service.ReceptionRepository {
	return &Reception{
		log: log,
		db:  db,
	}
}
