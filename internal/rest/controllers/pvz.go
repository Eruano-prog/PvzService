package controllers

import (
	"AvitoPvz/internal/rest"
	"AvitoPvz/internal/rest/dto"
	"AvitoPvz/internal/rest/middleware"
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

func (p *PVZController) Register(mux *http.ServeMux, tokenVerifier middleware.Verifier) {
	mux.Handle("POST /pvz", middleware.AuthMiddleware(http.HandlerFunc(p.createPVZHandler), tokenVerifier, p.log, middleware.ModeratorOnly))
	mux.Handle("GET /pvz", middleware.AuthMiddleware(http.HandlerFunc(p.getPVZsHandler), tokenVerifier, p.log, middleware.EmployeeAndModerator))
	mux.Handle("POST /pvz/{pvzId}/close_last_reception", middleware.AuthMiddleware(http.HandlerFunc(p.closeLastReceptionHandler), tokenVerifier, p.log, middleware.EmployeeAndModerator))
	mux.Handle("POST /pvz/{pvzId}/delete_last_product", middleware.AuthMiddleware(http.HandlerFunc(p.deleteLastProductHandler), tokenVerifier, p.log, middleware.EmployeeOnly))
}

func (p *PVZController) createPVZHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.PVZ
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		p.log.Debug("Error decoding create pvz request", "error", err)
		rest.WriteError(w, p.log, http.StatusBadRequest, "invalid request")
		return
	}

	pvz, err := req.ToModel()
	if err != nil {
		p.log.Debug("Error decoding create pvz request", "error", err)
		rest.WriteError(w, p.log, http.StatusBadRequest, "invalid request")
		return
	}

	createdPVZ, err := p.pvzService.CreatePVZ(r.Context(), pvz)
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := dto.PVZToDTO(*createdPVZ)
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		p.log.Error("failed to encode response", "error", err)
	}
}

func (p *PVZController) getPVZsHandler(w http.ResponseWriter, r *http.Request) {
	params := dto.GetPvzParams{}
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
		PVZ        dto.PVZ `json:"pvz"`
		Receptions []struct {
			Reception dto.Reception `json:"reception"`
			Products  []dto.Product `json:"products"`
		} `json:"receptions"`
	}, len(pvzs))

	for i, pvz := range pvzs {
		pvzDTO, err := dto.PVZToDTO(pvz.PVZ)
		if err != nil {
			p.log.Warn("failed to encode pvz to dto", "error", err)
			continue
		}

		response[i].PVZ = *pvzDTO

		response[i].Receptions = make([]struct {
			Reception dto.Reception `json:"reception"`
			Products  []dto.Product `json:"products"`
		}, len(pvz.Receptions))

		for j, rec := range pvz.Receptions {
			receptionDTO, err := dto.ReceptionToDTO(rec.Reception)
			if err != nil {
				p.log.Warn("failed to encode reception to dto", "error", err)
				continue
			}
			response[i].Receptions[j].Reception = *receptionDTO

			response[i].Receptions[j].Products = make([]dto.Product, len(rec.Products))
			for k, prod := range rec.Products {
				productDTO, err := dto.ProductToDTO(prod)
				if err != nil {
					p.log.Warn("failed to encode product to dto", "error", err)
					continue
				}

				response[i].Receptions[j].Products[k] = *productDTO
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

	resp, err := dto.ReceptionToDTO(*reception)
	if err != nil {
		rest.WriteError(w, p.log, http.StatusBadRequest, err.Error())
		return
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
