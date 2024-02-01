package service

import (
	"context"
	"eau-de-go/internal/repository"
	"eau-de-go/pkg/jwt"
	passwordPkg "eau-de-go/pkg/password"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	hashedPassword, err := passwordPkg.HashPassword(appUserParams.Password)
	appUserParams.Password = string(hashedPassword)
	dao, err := service.AppUserStore.CreateAppUser(ctx, appUserParams)
	if err != nil { // TODO: Refactor this error handling
		var dbErr *pq.Error
		if errors.As(err, &dbErr) {
			if dbErr.Code.Name() == "unique_violation" {
				return repository.AppUser{}, &repository.DuplicateKeyError{Key: "Duplicate user already exist."}
			}
		} else {
			log.Error(err)
			return repository.AppUser{}, err
		}
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
	if err := passwordPkg.CheckPassword(password, []byte(dao.Password)); err != nil {
		log.Error(err)
		return repository.AppUser{}, err
	}
	return dao, nil
}

func (service *AppUserService) GetAppUserTokens(appUser repository.AppUser) (string, map[string]interface{}, string, map[string]interface{}, error) {
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

	refreshToken, refreshTokenClaims, err := jwt.CreateRefreshToken(claims)
	if err != nil {
		log.Error(err)
		return "", nil, "", nil, err
	}

	accessToken, accessTokenClaims, err := jwt.CreateAccessToken(claims)
	if err != nil {
		log.Error(err)
		return "", nil, "", nil, err
	}
	return refreshToken, refreshTokenClaims, accessToken, accessTokenClaims, nil
}

func (service *AppUserService) RefreshToken(ctx context.Context, refreshToken string) (string, map[string]interface{}, error) {
	return jwt.CreateAccessTokenFromRefreshToken(refreshToken)
}
