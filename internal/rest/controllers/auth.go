package controllers

import (
	"AvitoPvz/internal/domain/models"
	"AvitoPvz/internal/rest"
	"AvitoPvz/internal/rest/dto"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type AuthController struct {
	log *slog.Logger

	userService rest.UserService
}

func NewAuthController(log *slog.Logger, userService rest.UserService) *AuthController {
	return &AuthController{
		log:         log,
		userService: userService,
	}
}

func (a AuthController) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /dummyLogin", a.dummyLoginHandler)
	mux.HandleFunc("POST /register", a.registerHandler)
	mux.HandleFunc("POST /login", a.loginHandler)
}

func (a AuthController) dummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var request dto.PostDummyLoginJSONRequestBody

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		a.log.Debug("Error during dummy login:", "error", err)
		rest.WriteError(w, a.log, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	role, err := models.GetRoleFromString(string(request.Role))
	if err != nil {
		a.log.Debug("Error during dummy login:", "error", err)
		rest.WriteError(w, a.log, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	var openApiToken dto.Token
	openApiToken, err = a.userService.DummyLogin(r.Context(), role)
	if err != nil {
		a.log.Debug("Error during dummy login:", "error", err)
		rest.WriteError(w, a.log, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(openApiToken); err != nil {
		a.log.Error("Error writing token to response")
	}
}

func (a AuthController) registerHandler(w http.ResponseWriter, r *http.Request) {
	var apiReq dto.PostRegisterJSONBody
	if err := json.NewDecoder(r.Body).Decode(&apiReq); err != nil {
		a.log.Debug("Error during registration:", "error", err)
		rest.WriteError(w, a.log, http.StatusBadRequest, fmt.Sprintf("invalid request. Error: %v", err))
		return
	}

	role, err := dto.RoleToModel(dto.UserRole(apiReq.Role))
	if err != nil {
		a.log.Debug("Error during registration:", "error", err)
		rest.WriteError(w, a.log, http.StatusBadRequest, fmt.Sprintf("invalid request. Error: %v", err))
		return
	}

	createdUser, err := a.userService.Register(r.Context(), string(apiReq.Email), apiReq.Password, role)
	if err != nil {
		a.log.Debug("Error during registration:", "error", err)
		rest.WriteError(w, a.log, http.StatusBadRequest, err.Error())
		return
	}

	apiResp, err := dto.UserToDTO(*createdUser)
	if err != nil {
		a.log.Debug("Error during registration:", "error", err)
		rest.WriteError(w, a.log, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(apiResp); err != nil {
		a.log.Error("failed to encode response", "error", err)
	}
}

func (a AuthController) loginHandler(w http.ResponseWriter, r *http.Request) {
	var apiReq dto.PostLoginJSONBody
	if err := json.NewDecoder(r.Body).Decode(&apiReq); err != nil {
		a.log.Debug("Error during dummy login:", "error", err)
		rest.WriteError(w, a.log, http.StatusBadRequest, fmt.Sprintf("invalid request. Error: %v", err))
		return
	}

	var token dto.Token
	token, err := a.userService.Login(r.Context(), string(apiReq.Email), apiReq.Password)
	if err != nil {
		a.log.Debug("Error during dummy login:", "error", err)
		rest.WriteError(w, a.log, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err = json.NewEncoder(w).Encode(token); err != nil {
		a.log.Error("failed to encode response", "error", err)
	}
}
