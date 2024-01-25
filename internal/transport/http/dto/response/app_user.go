package response

import (
	"eau-de-go/internal/repository"
	"github.com/google/uuid"
)

type AppUserDto struct {
	ID        uuid.UUID
	FirstName *string
	LastName  *string
	Email     string
	LastLogin *string
	IsActive  bool
}

func ConvertDbRow(user repository.AppUser) AppUserDto {
	var firstName, lastName, lastLogin *string
	if user.FirstName.Valid {
		firstName = &user.FirstName.String
	}
	if user.LastName.Valid {
		lastName = &user.LastName.String
	}

	if user.LastLogin.Valid {
		lastLoginStr := user.LastLogin.Time.String()
		lastLogin = &lastLoginStr
	}

	return AppUserDto{
		ID:        user.ID,
		FirstName: firstName,
		LastName:  lastName,
		Email:     user.Email,
		LastLogin: lastLogin,
		IsActive:  user.IsActive,
	}
}
