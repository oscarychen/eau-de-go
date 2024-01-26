package http

import (
	"context"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/transport/http/dto/request"
	"eau-de-go/internal/transport/http/dto/response"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AppUserService interface {
	GetAppUserById(ctx context.Context, ID uuid.UUID) (repository.AppUser, error)
	CreateAppUser(ctx context.Context, appUser repository.CreateAppUserParams) (repository.AppUser, error)
}

func (h *Handler) GetAppUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDao, err := h.Service.GetAppUserById(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	userDto := response.ConvertDbRow(userDao)

	jsonData, err := json.Marshal(userDto)
	if err != nil {
		log.Errorf("Error marshalling json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
		return
	}
}

func (h *Handler) CreateAppUser(w http.ResponseWriter, r *http.Request) {

	createAppUserParams, err := request.MakeCreateAppUserParamsFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userDao, err := h.Service.CreateAppUser(r.Context(), createAppUserParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDto := response.ConvertDbRow(userDao)

	jsonData, err := json.Marshal(userDto)
	if err != nil {
		log.Errorf("Error marshalling json: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
		return
	}
}
