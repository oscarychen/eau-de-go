package http

import (
	"context"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/transport/http/request_dto"
	"eau-de-go/internal/transport/http/response_dto"
	"eau-de-go/settings"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type AppUserService interface {
	Login(ctx context.Context, username string, password string) (repository.AppUser, error)
	GetAppUserById(ctx context.Context, ID uuid.UUID) (repository.AppUser, error)
	CreateAppUser(ctx context.Context, appUserParams repository.CreateAppUserParams) (repository.AppUser, error)
	UpdateAppUser(ctx context.Context, appUserParams repository.UpdateAppUserParams) (repository.AppUser, error)
	UpdateAppUserPassword(ctx context.Context, userId uuid.UUID, oldPassword string, newPassword string) (repository.AppUser, error)
	GetAppUserTokens(appUser repository.AppUser) (string, map[string]interface{}, string, map[string]interface{}, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, map[string]interface{}, repository.AppUser, error)
}

var refreshTokenCookieName = "refresh"
var refreshTokenCookiePath = "/auth/token-refresh"

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var loginDto request_dto.AppUserLoginRequestDto
	err = json.Unmarshal(bodyBytes, &loginDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDao, err := h.AppUserService.Login(r.Context(), loginDto.Username, loginDto.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	refreshToken, refreshTokenClaims, accessToken, _, err := h.AppUserService.GetAppUserTokens(userDao)

	cookie := http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refreshToken,
		HttpOnly: true,
		Path:     refreshTokenCookiePath,
		SameSite: http.SameSiteStrictMode,
		Secure:   settings.RefreshCookieSecure,
		Expires:  time.Unix(refreshTokenClaims["exp"].(int64), 0),
	}
	http.SetCookie(w, &cookie)

	userDto := response_dto.ConvertDbRow(userDao)
	responseData := response_dto.AppUserLoginResponse{
		AppUserDto:  userDto,
		AccessToken: accessToken,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		log.Errorf("Error marshalling json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
		return
	}
}

func (h *Handler) TokenRefresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie(refreshTokenCookieName)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	refreshToken := refreshTokenCookie.Value
	accessToken, _, appUser, err := h.AppUserService.RefreshToken(r.Context(), refreshToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	userDto := response_dto.ConvertDbRow(appUser)

	responseData := response_dto.AppUserLoginResponse{
		AppUserDto:  userDto,
		AccessToken: accessToken,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		log.Errorf("Error marshalling json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	userDto := response_dto.ConvertDbRow(userDao)

	jsonData, err := json.Marshal(userDto)
	if err != nil {
		log.Errorf("Error marshalling json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		log.Errorf("Error writing response: %v", err)
		return
	}
}

func (h *Handler) CreateAppUser(w http.ResponseWriter, r *http.Request) {

	appUserParams, err := request_dto.MakeCreateAppUserParamsFromRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	userDao, err := h.AppUserService.CreateAppUser(r.Context(), appUserParams)
	if err != nil { // TODO: Refactor this error handling
		var duplicateKeyError *repository.DuplicateKeyError

		if errors.As(err, &duplicateKeyError) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	userDto := response_dto.ConvertDbRow(userDao)

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

func (h *Handler) UpdateAppUser(w http.ResponseWriter, r *http.Request) {

	appUserParams, err := request_dto.MakeUpdateAppUserParamsFromRequest(r)
	if err != nil {
		log.Errorf("Error unmarshalling json: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDao, err := h.AppUserService.UpdateAppUser(r.Context(), appUserParams)
	if err != nil {
		log.Errorf("Error updating user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDto := response_dto.ConvertDbRow(userDao)

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

func (h *Handler) UpdateAppUserPassword(w http.ResponseWriter, r *http.Request) {

	jwtClaims, ok := r.Context().Value("jwt_claims").(map[string]interface{})
	if !ok {
		http.Error(w, "Error parsing access token", http.StatusUnauthorized)
		return
	}

	idStr, ok := jwtClaims["id"].(string)
	if !ok {
		http.Error(w, "Error parsing access token", http.StatusUnauthorized)
	}

	userId, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	var updatePasswordDto request_dto.UpdateAppUserPasswordRequestDto
	err = json.NewDecoder(r.Body).Decode(&updatePasswordDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDao, err := h.AppUserService.UpdateAppUserPassword(r.Context(), userId, updatePasswordDto.OldPassword, updatePasswordDto.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userDto := response_dto.ConvertDbRow(userDao)

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
