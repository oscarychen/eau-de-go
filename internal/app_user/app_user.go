package app_user

import (
	"context"
	"eau-de-go/internal/repository"
	"fmt"
	"github.com/google/uuid"
)

type AppUserStore interface {
	GetAppUserById(ctx context.Context, id uuid.UUID) (repository.AppUser, error)
}

type AppUserService struct {
	AppUserStore AppUserStore
}

func NewAppUserService(appUserStore AppUserStore) *AppUserService {
	return &AppUserService{
		AppUserStore: appUserStore,
	}
}

func (service *AppUserService) GetAppUserById(ctx context.Context, id uuid.UUID) (AppUserDto, error) {
	userRow, err := service.AppUserStore.GetAppUserById(ctx, id)
	if err != nil {
		fmt.Println(err)
		return AppUserDto{}, err
	}
	return convertDbRow(userRow), nil
}
