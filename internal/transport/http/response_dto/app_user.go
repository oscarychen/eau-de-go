package response_dto

import (
	"eau-de-go/internal/repository"
	"github.com/google/uuid"
)

type AppUserDto struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	Email     string    `json:"email"`
	LastLogin *string   `json:"last_login,omitempty"`
	IsActive  bool      `json:"is_active"`
}

type AppUserLoginResponse struct {
	AppUserDto
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
