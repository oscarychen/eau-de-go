package request_dto

import (
	"database/sql"
	"eau-de-go/internal/repository"
	"eau-de-go/pkg/jwt_util"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type AppUserLoginRequestDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateAppUserRequestDto struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	IsActive  bool   `json:"is_active"`
}

func MakeCreateAppUserParamsFromRequest(r *http.Request) (repository.CreateAppUserParams, error) {
	var dto CreateAppUserRequestDto
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return repository.CreateAppUserParams{}, err
	}
	err = json.Unmarshal(bodyBytes, &dto)
	if err != nil {
		return repository.CreateAppUserParams{}, err
	}
	createAppUserParams := repository.CreateAppUserParams{
		Username:  dto.Username,
		Email:     dto.Email,
		Password:  dto.Password,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
	}
	return createAppUserParams, nil
}

type UpdateAppUserRequestDto struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}

func MakeUpdateAppUserParamsFromRequest(r *http.Request) (repository.UpdateAppUserParams, error) {
	var dto UpdateAppUserRequestDto
	var firstName, lastName sql.NullString

	jwtClaims, ok := r.Context().Value("jwt_claims").(map[string]interface{})
	if !ok {
		return repository.UpdateAppUserParams{}, &jwt_util.InvalidTokenError{}
	}

	idStr, ok := jwtClaims["id"].(string)
	if !ok {
		return repository.UpdateAppUserParams{}, &jwt_util.InvalidTokenError{}
	}

	userId, err := uuid.Parse(idStr)
	if err != nil {
		return repository.UpdateAppUserParams{}, &jwt_util.InvalidTokenError{}
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return repository.UpdateAppUserParams{}, err
	}

	err = json.Unmarshal(bodyBytes, &dto)
	if err != nil {
		return repository.UpdateAppUserParams{}, err
	}

	if dto.FirstName == nil {
		firstName = sql.NullString{Valid: false}
	} else {
		firstName = sql.NullString{String: *dto.FirstName, Valid: true}
	}

	if dto.LastName == nil {
		lastName = sql.NullString{Valid: false}
	} else {
		lastName = sql.NullString{String: *dto.LastName, Valid: true}
	}

	updateAppUserParams := repository.UpdateAppUserParams{
		ID:        userId,
		FirstName: firstName,
		LastName:  lastName,
	}
	return updateAppUserParams, nil
}
