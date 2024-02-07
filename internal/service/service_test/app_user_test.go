package service_test

import (
	"context"
	"database/sql"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/service"
	"eau-de-go/pkg/password_util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockAppUserStore struct {
	mock.Mock
}

func (m *MockAppUserStore) UpdateAppUser(ctx context.Context, appUser repository.UpdateAppUserParams) (repository.AppUser, error) {
	args := m.Called(ctx, appUser)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserStore) UpdateAppUserPassword(ctx context.Context, appUser repository.UpdateAppUserPasswordParams) (repository.AppUser, error) {
	args := m.Called(ctx, appUser)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserStore) GetAppUserByUsername(ctx context.Context, username string) (repository.AppUser, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserStore) GetAppUserById(ctx context.Context, id uuid.UUID) (repository.AppUser, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserStore) CreateAppUser(ctx context.Context, appUser repository.CreateAppUserParams) (repository.AppUser, error) {
	args := m.Called(ctx, appUser)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func TestCreateAppUser(t *testing.T) {
	mockStore := new(MockAppUserStore)
	aps := service.NewAppUserService(mockStore)
	service.HashPasswordFunc = func(password string) ([]byte, error) {
		return []byte(password), nil
	}

	userParams := repository.CreateAppUserParams{
		Username: "testuser",
		Password: "testPassword",
		Email:    "testuser@example.com",
	}

	mockStore.On("CreateAppUser", mock.Anything, userParams).Return(repository.AppUser{}, nil)

	_, err := aps.CreateAppUser(context.Background(), userParams)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
}

func TestGetAppUserById(t *testing.T) {
	mockStore := new(MockAppUserStore)
	aps := service.NewAppUserService(mockStore)

	id := uuid.New()
	mockStore.On("GetAppUserById", mock.Anything, id).Return(repository.AppUser{}, nil)

	_, err := aps.GetAppUserById(context.Background(), id)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
}

func TestUpdateAppUser(t *testing.T) {
	mockStore := new(MockAppUserStore)
	aps := service.NewAppUserService(mockStore)

	userParams := repository.UpdateAppUserParams{
		ID:       uuid.New(),
		LastName: sql.NullString{String: "updateduser"},
	}

	mockStore.On("UpdateAppUser", mock.Anything, userParams).Return(repository.AppUser{}, nil)

	_, err := aps.UpdateAppUser(context.Background(), userParams)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
}

func TestUpdateAppUserPassword(t *testing.T) {
	mockStore := new(MockAppUserStore)
	aps := service.NewAppUserService(mockStore)

	id := uuid.New()
	oldPassword := "oldPassword"
	hashedOldPassword, _ := password_util.HashPassword(oldPassword)
	newPassword := "newPassword"

	mockStore.On("GetAppUserById", mock.Anything, id).Return(repository.AppUser{Password: string(hashedOldPassword)}, nil)
	mockStore.On("UpdateAppUserPassword", mock.Anything, mock.AnythingOfType("repository.UpdateAppUserPasswordParams")).Return(repository.AppUser{}, nil)

	_, err := aps.UpdateAppUserPassword(context.Background(), id, oldPassword, newPassword)
	assert.NoError(t, err)

	mockStore.AssertExpectations(t)
}
