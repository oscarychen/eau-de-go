package service

import (
	"context"
	"eau-de-go/internal/repository"
	"eau-de-go/pkg/email_util"
	"eau-de-go/pkg/jwt_util"
	"eau-de-go/pkg/password_util"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/url"
)

type AppUserStore interface {
	GetAppUserById(ctx context.Context, id uuid.UUID) (repository.AppUser, error)
	CreateAppUser(ctx context.Context, appUser repository.CreateAppUserParams) (repository.AppUser, error)
	UpdateAppUser(ctx context.Context, appUser repository.UpdateAppUserParams) (repository.AppUser, error)
	UpdateAppUserPassword(ctx context.Context, appUser repository.UpdateAppUserPasswordParams) (repository.AppUser, error)
	GetAppUserByUsername(ctx context.Context, username string) (repository.AppUser, error)
	SetUserEmailVerified(ctx context.Context, userId uuid.UUID) (repository.AppUser, error)
	SetUserEmailUnverified(ctx context.Context, userId uuid.UUID) (repository.AppUser, error)
	UpdateAppUserLastLogin(ctx context.Context, data repository.UpdateAppUserLastLoginParams) (repository.AppUser, error)
}

type AppUserService struct {
	AppUserStore  AppUserStore
	EmailVerifier email_util.EmailTokenVerifier
	EmailSender   email_util.EmailSender
}

func NewAppUserService(appUserStore AppUserStore) *AppUserService {
	return &AppUserService{
		AppUserStore:  appUserStore,
		EmailVerifier: email_util.NewEmailTokenVerifier(),
		EmailSender:   email_util.NewEmailSender(),
	}
}

var HashPasswordFunc = password_util.HashPassword

func (service *AppUserService) CreateAppUser(ctx context.Context, appUserParams repository.CreateAppUserParams) (repository.AppUser, error) {
	hashedPassword, err := HashPasswordFunc(appUserParams.Password)
	if err != nil {
		return repository.AppUser{}, err
	}
	appUserParams.Password = string(hashedPassword)

	validatedEmail, err := email_util.ValidateEmailAddress(appUserParams.Email)
	if err != nil {
		return repository.AppUser{}, err
	}
	appUserParams.Email = validatedEmail

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

func (service *AppUserService) SendUserEmailVerification(ctx context.Context, emailAddress string) error {

	token, err := service.EmailVerifier.CreateToken(emailAddress)
	urlSafeToken := url.QueryEscape(token)
	if err != nil {
		log.Error(err)
		return err
	}
	//TODO: front end url from settings
	err = service.EmailSender.SendSingleEmail(emailAddress, "Email Verification", urlSafeToken)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (service *AppUserService) VerifyEmailVerificationToken(ctx context.Context, userId uuid.UUID, userEmailAddress string, token string) (bool, error) {
	verifiedEmail, err := service.EmailVerifier.VerifyToken(token)
	if err != nil {
		log.Error(err)
		return false, err
	}

	if verifiedEmail != userEmailAddress {
		return false, errors.New("email address does not match")
	}

	_, err = service.AppUserStore.SetUserEmailVerified(ctx, userId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	return true, nil
}

func (service *AppUserService) UpdateAppUser(ctx context.Context, appUserParams repository.UpdateAppUserParams) (repository.AppUser, error) {
	dao, err := service.AppUserStore.UpdateAppUser(ctx, appUserParams)
	if err != nil {
		log.Error(err)
		return repository.AppUser{}, err
	}
	return dao, nil
}

func (service *AppUserService) UpdateAppUserPassword(ctx context.Context, userId uuid.UUID, oldPassword string, newPassword string) (repository.AppUser, error) {

	if oldPassword == newPassword {
		return repository.AppUser{}, &password_util.SamePasswordError{}
	}

	dao, err := service.AppUserStore.GetAppUserById(ctx, userId)
	if err != nil {
		log.Error(err)
		return repository.AppUser{}, err
	}

	err = password_util.CheckPassword(oldPassword, []byte(dao.Password))
	if err != nil {
		return repository.AppUser{}, &repository.IncorrectUserCredentialError{}
	}
	hashedNewPassword, err := HashPasswordFunc(newPassword)
	if err != nil {
		return repository.AppUser{}, err
	}
	appUserParams := repository.UpdateAppUserPasswordParams{
		ID:       userId,
		Password: string(hashedNewPassword),
	}
	dao, err = service.AppUserStore.UpdateAppUserPassword(ctx, appUserParams)
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
		return repository.AppUser{}, &repository.IncorrectUserCredentialError{}
	}
	if err := password_util.CheckPassword(password, []byte(dao.Password)); err != nil {
		return repository.AppUser{}, &repository.IncorrectUserCredentialError{}
	}
	return dao, nil
}

func (service *AppUserService) makeTokenClaimMap(appUser repository.AppUser) map[string]interface{} {
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
	claims["email_verified"] = appUser.EmailVerified
	return claims
}

func (service *AppUserService) GetAppUserTokens(appUser repository.AppUser) (string, map[string]interface{}, string, map[string]interface{}, error) {
	claims := service.makeTokenClaimMap(appUser)

	refreshToken, refreshTokenClaims, err := jwt_util.CreateRefreshToken(claims)
	if err != nil {
		log.Error(err)
		return "", nil, "", nil, err
	}

	accessToken, accessTokenClaims, err := jwt_util.CreateAccessToken(claims)
	if err != nil {
		log.Error(err)
		return "", nil, "", nil, err
	}
	return refreshToken, refreshTokenClaims, accessToken, accessTokenClaims, nil
}

func (service *AppUserService) RefreshToken(ctx context.Context, refreshToken string) (string, map[string]interface{}, repository.AppUser, error) {
	refreshTokenClaims, err := jwt_util.DecodeToken(jwt_util.Refresh, refreshToken)
	if err != nil {
		log.Error(err)
		return "", nil, repository.AppUser{}, &jwt_util.InvalidTokenError{}
	}

	idStr, ok := refreshTokenClaims["id"].(string)
	if !ok {
		return "", nil, repository.AppUser{}, &jwt_util.InvalidTokenError{}
	}

	userId, err := uuid.Parse(idStr)
	if err != nil {
		log.Error(err)
		return "", nil, repository.AppUser{}, &jwt_util.InvalidTokenError{}
	}

	appUser, err := service.GetAppUserById(ctx, userId)
	tokenClaims := service.makeTokenClaimMap(appUser)
	accessToken, claims, err := jwt_util.CreateAccessToken(tokenClaims)
	if err != nil {
		log.Error(err)
		return "", nil, repository.AppUser{}, err
	}

	return accessToken, claims, appUser, nil
}
