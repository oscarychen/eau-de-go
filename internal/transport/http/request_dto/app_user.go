package request_dto

import (
	"database/sql"
	"eau-de-go/internal/repository"
	"encoding/json"
	"io"
	"net/http"
)

type AppUserLoginRequestDto struct {
	Username string
	Password string
}

type CreateAppUserRequestDto struct {
	Username  string
	FirstName sql.NullString
	LastName  sql.NullString
	Email     string
	Password  string
	IsActive  bool
}

func (c *CreateAppUserRequestDto) UnmarshalJSON(data []byte) error {
	type Alias CreateAppUserRequestDto
	aux := &struct {
		FirstName string
		LastName  string
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.FirstName = sql.NullString{String: aux.FirstName, Valid: aux.FirstName != ""}
	c.LastName = sql.NullString{String: aux.LastName, Valid: aux.LastName != ""}
	return nil
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
