package response_dto

import (
	"eau-de-go/internal/repository"
	"github.com/google/uuid"
)

type AppUserDto struct {
	ID        uuid.UUID
	Username  string
	FirstName *string `json:",omitempty"`
	LastName  *string `json:",omitempty"`
	Email     string
	LastLogin *string `json:",omitempty"`
	IsActive  bool
}

type AppUserLoginResponse struct {
	AppUserDto
	AccessToken  string
	RefreshToken string
}

func ConvertDbRow(user repository.AppUser) AppUserDto {
	var lastLogin *string

	if user.LastLogin.Valid {
		lastLoginStr := user.LastLogin.Time.String()
		lastLogin = &lastLoginStr
	}

	return AppUserDto{
		ID:        user.ID,
		Username:  user.Username,
		FirstName: &user.FirstName,
		LastName:  &user.LastName,
		Email:     user.Email,
		LastLogin: lastLogin,
		IsActive:  user.IsActive,
	}
}
