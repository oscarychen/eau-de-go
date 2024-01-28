package http

import (
	"context"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/settings"
	"eau-de-go/internal/transport/http/dto/request"
	"eau-de-go/internal/transport/http/dto/response"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

var NowFunc = time.Now

type AppUserService interface {
	Login(ctx context.Context, username string, password string) (repository.AppUser, error)
	GetAppUserById(ctx context.Context, ID uuid.UUID) (repository.AppUser, error)
	CreateAppUser(ctx context.Context, appUser repository.CreateAppUserParams) (repository.AppUser, error)
	GetAppUserTokens(appUser repository.AppUser) (string, string, error)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var loginDto request.AppUserLoginRequestDto
	err = json.Unmarshal(bodyBytes, &loginDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDao, err := h.AppUserService.Login(r.Context(), loginDto.Username, loginDto.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	refreshToken, accessToken, err := h.AppUserService.GetAppUserTokens(userDao)

	cookie := http.Cookie{
		Name:     "refresh",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/api/user/refresh",
		SameSite: http.SameSiteStrictMode,
		Secure:   settings.RefreshCookieSecure,
		Expires:  NowFunc().Add(settings.RefreshTokenLife),
	}
	http.SetCookie(w, &cookie)

	userDto := response.ConvertDbRow(userDao)
	responseData := response.AppUserLoginResponse{
		AppUserDto:  userDto,
		AccessToken: accessToken,
	}

	jsonData, err := json.Marshal(responseData)
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

func (h *Handler) GetAppUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDao, err := h.AppUserService.GetAppUserById(r.Context(), id)
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
	userDao, err := h.AppUserService.CreateAppUser(r.Context(), createAppUserParams)
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
