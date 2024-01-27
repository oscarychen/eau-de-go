package app_user

import (
	"context"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AppUserStore interface {
	GetAppUserById(ctx context.Context, id uuid.UUID) (repository.AppUser, error)
	CreateAppUser(ctx context.Context, appUser repository.CreateAppUserParams) (repository.AppUser, error)
}

type AppUserService struct {
	AppUserStore AppUserStore
}

func NewAppUserService(appUserStore AppUserStore) *AppUserService {
	return &AppUserService{
		AppUserStore: appUserStore,
	}
}

func (service *AppUserService) CreateAppUser(ctx context.Context, appUserParams repository.CreateAppUserParams) (repository.AppUser, error) {
	hashedPassword, err := utils.HashPassword(appUserParams.Password)
	appUserParams.Password = string(hashedPassword)
	dao, err := service.AppUserStore.CreateAppUser(ctx, appUserParams)
	if err != nil {
		log.Error(err)
		return repository.AppUser{}, err
	}
	return dao, nil
}

func (service *AppUserService) GetAppUserById(ctx context.Context, id uuid.UUID) (repository.AppUser, error) {
	dao, err := service.AppUserStore.GetAppUserById(ctx, id)
	if err != nil {
		log.Error(err)
		return repository.AppUser{}, err
	}
	return dao, nil
}
