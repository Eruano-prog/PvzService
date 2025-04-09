package controllers

import (
	"AvitoPvz/internal/rest"
	"encoding/json"
	"log/slog"
	"net/http"
)

type ReceptionController struct {
	log              *slog.Logger
	receptionService rest.ReceptionService
}

func NewReceptionController(log *slog.Logger, receptionService rest.ReceptionService) *ReceptionController {
	return &ReceptionController{
		log:              log,
		receptionService: receptionService,
	}
}

func (r *ReceptionController) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /receptions", r.createReceptionHandler)
}

func (r *ReceptionController) createReceptionHandler(w http.ResponseWriter, req *http.Request) {
	var request rest.PostReceptionsJSONRequestBody
	if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
		rest.WriteError(w, r.log, http.StatusBadRequest, "invalid request")
		return
	}

	reception, err := r.receptionService.CreateReception(req.Context(), request.PvzId)
	if err != nil {
		rest.WriteError(w, r.log, http.StatusBadRequest, err.Error())
		return
	}

	resp := rest.Reception{
		Id:       &reception.ID,
		DateTime: reception.DateTime,
		PvzId:    reception.PVZID,
		Status:   rest.ReceptionStatus(reception.Status),
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		r.log.Error("failed to encode response", "error", err)
	}
}
