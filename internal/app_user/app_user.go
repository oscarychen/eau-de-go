package app_user

import (
	"context"
	"eau-de-go/internal/jwt_auth"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AppUserStore interface {
	GetAppUserById(ctx context.Context, id uuid.UUID) (repository.AppUser, error)
	CreateAppUser(ctx context.Context, appUser repository.CreateAppUserParams) (repository.AppUser, error)
	GetAppUserByUsername(ctx context.Context, username string) (repository.AppUser, error)
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

func (service *AppUserService) GetAppUserByUsername(ctx context.Context, username string) (repository.AppUser, error) {
	dao, err := service.AppUserStore.GetAppUserByUsername(ctx, username)
	if err != nil {
		log.Error(err)
		return repository.AppUser{}, err
	}
	return dao, nil
}

func (service *AppUserService) Login(ctx context.Context, username string, password string) (repository.AppUser, error) {
	dao, err := service.AppUserStore.GetAppUserByUsername(ctx, username)
	if err != nil {
		log.Error(err)
		return repository.AppUser{}, err
	}
	if err := utils.CheckPassword(password, []byte(dao.Password)); err != nil {
		log.Error(err)
		return repository.AppUser{}, err
	}
	return dao, nil
}

func (service *AppUserService) GetAppUserTokens(appUser repository.AppUser) (string, string, error) {
	claims := make(map[string]interface{})
	claims["id"] = appUser.ID
	claims["username"] = appUser.Username
	claims["email"] = appUser.Email
	claims["first_name"] = appUser.FirstName
	claims["last_name"] = appUser.LastName
	claims["is_active"] = appUser.IsActive
	claims["is_staff"] = appUser.IsStaff
	claims["last_login"] = appUser.LastLogin
	claims["date_joined"] = appUser.DateJoined

	refreshToken, err := jwt_auth.CreateToken(jwt_auth.Refresh, claims)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	accessToken, err := jwt_auth.CreateToken(jwt_auth.Access, claims)
	if err != nil {
		log.Error(err)
		return "", "", err
	}
	return refreshToken, accessToken, nil
}
