package response_dto_test

import (
	"database/sql"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/transport/http/response_dto"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertDbRow_HappyPath(t *testing.T) {
	user := repository.AppUser{
		ID:        uuid.New(),
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Email:     "testuser@example.com",
		LastLogin: sql.NullTime{Time: time.Now(), Valid: true},
		IsActive:  true,
	}

	dto := response_dto.ConvertDbRow(user)

	assert.Equal(t, user.ID, dto.ID)
	assert.Equal(t, user.Username, dto.Username)
	assert.Equal(t, user.FirstName, *dto.FirstName)
	assert.Equal(t, user.LastName, *dto.LastName)
	assert.Equal(t, user.Email, dto.Email)
	assert.Equal(t, user.LastLogin.Time.String(), *dto.LastLogin)
	assert.Equal(t, user.IsActive, dto.IsActive)
}

func TestConvertDbRow_NullLastLogin(t *testing.T) {
	user := repository.AppUser{
		ID:        uuid.New(),
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Email:     "testuser@example.com",
		LastLogin: sql.NullTime{Valid: false},
		IsActive:  true,
	}

	dto := response_dto.ConvertDbRow(user)

	assert.Equal(t, user.ID, dto.ID)
	assert.Equal(t, user.Username, dto.Username)
	assert.Equal(t, user.FirstName, *dto.FirstName)
	assert.Equal(t, user.LastName, *dto.LastName)
	assert.Equal(t, user.Email, dto.Email)
	assert.Nil(t, dto.LastLogin)
	assert.Equal(t, user.IsActive, dto.IsActive)
}
