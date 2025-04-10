package controllers

import (
	"AvitoPvz/internal/rest"
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type PVZController struct {
	log        *slog.Logger
	pvzService rest.PVZService
}

func NewPVZController(log *slog.Logger, pvzService rest.PVZService) *PVZController {
	return &PVZController{
		log:        log,
		pvzService: pvzService,
	}
}

func (p *PVZController) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /pvz", p.createPVZHandler)
	mux.HandleFunc("GET /pvz", p.getPVZsHandler)
	mux.HandleFunc("POST /pvz/{pvzId}/close_last_reception", p.closeLastReceptionHandler)
	mux.HandleFunc("POST /pvz/{pvzId}/delete_last_product", p.deleteLastProductHandler)
}

// TODO: city validation
func (p *PVZController) createPVZHandler(w http.ResponseWriter, r *http.Request) {
	var req rest.PVZ
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, "invalid request")
		return
	}

	createdPVZ, err := p.pvzService.CreatePVZ(r.Context(), string(req.City))
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, err.Error())
		return
	}

	resp := rest.PVZ{
		City:             rest.PVZCity(createdPVZ.City),
		Id:               &createdPVZ.ID,
		RegistrationDate: &createdPVZ.RegistrationDate,
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		p.log.Error("failed to encode response", "error", err)
	}
}

func (p *PVZController) getPVZsHandler(w http.ResponseWriter, r *http.Request) {
	params := rest.GetPvzParams{}
	var err error

	if startDateStr := r.URL.Query().Get("startDate"); startDateStr != "" {
		params.StartDate = new(time.Time)
		*params.StartDate, err = time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			rest.WriteError(w, p.log, http.StatusBadRequest, "invalid startDate format")
			return
		}
	}

	if endDateStr := r.URL.Query().Get("endDate"); endDateStr != "" {
		params.EndDate = new(time.Time)
		*params.EndDate, err = time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			rest.WriteError(w, p.log, http.StatusBadRequest, "invalid endDate format")
			return
		}
	}

	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			rest.WriteError(w, p.log, http.StatusBadRequest, "invalid page value")
			return
		}
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 30 {
			rest.WriteError(w, p.log, http.StatusBadRequest, "invalid limit value")
			return
		}
	}

	pvzs, err := p.pvzService.GetPVZsWithReceptions(r.Context(), params.StartDate, params.EndDate, page, limit)
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, err.Error())
		return
	}

	response := make([]struct {
		PVZ        rest.PVZ `json:"pvz"`
		Receptions []struct {
			Reception rest.Reception `json:"reception"`
			Products  []rest.Product `json:"products"`
		} `json:"receptions"`
	}, len(pvzs))

	for i, pvz := range pvzs {
		response[i].PVZ = rest.PVZ{
			Id:               &pvz.PVZ.ID,
			RegistrationDate: &pvz.PVZ.RegistrationDate,
			City:             rest.PVZCity(pvz.PVZ.City),
		}

		response[i].Receptions = make([]struct {
			Reception rest.Reception `json:"reception"`
			Products  []rest.Product `json:"products"`
		}, len(pvz.Receptions))

		for j, rec := range pvz.Receptions {
			response[i].Receptions[j].Reception = rest.Reception{
				Id:       &rec.Reception.ID,
				DateTime: rec.Reception.DateTime,
				PvzId:    rec.Reception.PVZID,
				Status:   rest.ReceptionStatus(rec.Reception.Status),
			}

			response[i].Receptions[j].Products = make([]rest.Product, len(rec.Products))
			for k, prod := range rec.Products {
				response[i].Receptions[j].Products[k] = rest.Product{
					Id:          &prod.ID,
					DateTime:    &prod.DateTime,
					Type:        rest.ProductType(prod.Type),
					ReceptionId: prod.ReceptionID,
				}
			}
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		p.log.Error("failed to encode response", "error", err)
	}
}

func (p *PVZController) closeLastReceptionHandler(w http.ResponseWriter, r *http.Request) {
	pvzIDStr := r.PathValue("pvzId")
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, "invalid pvzId format")
		return
	}

	reception, err := p.pvzService.CloseLastReception(r.Context(), pvzID)
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, err.Error())
		return
	}

	resp := rest.Reception{
		Id:       &reception.ID,
		DateTime: reception.DateTime,
		PvzId:    reception.PVZID,
		Status:   rest.ReceptionStatus(reception.Status),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		p.log.Error("failed to encode response", "error", err)
	}
}

func (p *PVZController) deleteLastProductHandler(w http.ResponseWriter, r *http.Request) {
	pvzIDStr := r.PathValue("pvzId")
	pvzID, err := uuid.Parse(pvzIDStr)
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, "invalid pvzId format")
		return
	}

	if err := p.pvzService.DeleteLastProduct(r.Context(), pvzID); err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
