package http_test

import (
	"bytes"
	"context"
	"eau-de-go/internal/repository"
	transportHttp "eau-de-go/internal/transport/http"
	"eau-de-go/internal/transport/http/request_dto"
	"eau-de-go/internal/transport/http/response_dto"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var refreshTokenCookieName = "refresh"

type MockAppUserService struct {
	mock.Mock
}

func (m *MockAppUserService) SendUserEmailVerification(ctx context.Context, emailAddress string) error {
	args := m.Called(ctx, emailAddress)
	return args.Error(0)
}

func (m *MockAppUserService) VerifyEmailVerificationToken(ctx context.Context, userId uuid.UUID, emailAddress string, token string) (bool, error) {
	args := m.Called(ctx, userId, emailAddress, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockAppUserService) UpdateAppUserPassword(ctx context.Context, userId uuid.UUID, oldPassword string, newPassword string) (repository.AppUser, error) {
	args := m.Called(ctx, userId, oldPassword, newPassword)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserService) Login(ctx context.Context, username string, password string) (repository.AppUser, error) {
	args := m.Called(ctx, username, password)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserService) GetAppUserById(ctx context.Context, ID uuid.UUID) (repository.AppUser, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserService) CreateAppUser(ctx context.Context, appUserParams repository.CreateAppUserParams) (repository.AppUser, error) {
	args := m.Called(ctx, appUserParams)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserService) UpdateAppUser(ctx context.Context, appUserParams repository.UpdateAppUserParams) (repository.AppUser, error) {
	args := m.Called(ctx, appUserParams)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserService) GetAppUserTokens(appUser repository.AppUser) (string, map[string]interface{}, string, map[string]interface{}, error) {
	args := m.Called(appUser)
	return args.String(0), args.Get(1).(map[string]interface{}), args.String(2), args.Get(3).(map[string]interface{}), args.Error(4)
}

func (m *MockAppUserService) RefreshToken(ctx context.Context, refreshToken string) (string, map[string]interface{}, repository.AppUser, error) {
	args := m.Called(ctx, refreshToken)
	return args.String(0), args.Get(1).(map[string]interface{}), args.Get(2).(repository.AppUser), args.Error(3)
}

func TestLoginSuccessful(t *testing.T) {
	mockService := new(MockAppUserService)
	handler := transportHttp.Handler{AppUserService: mockService}

	loginDto := request_dto.AppUserLoginRequestDto{Username: "test", Password: "test"}
	loginDtoBytes, _ := json.Marshal(loginDto)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(loginDtoBytes))

	expectedUser := repository.AppUser{ID: uuid.New(), Username: "test"}

	var mockExp int64 = 1707105923
	mockService.On("GetAppUserTokens", expectedUser).Return("refreshToken", map[string]interface{}{"exp": mockExp}, "accessToken", map[string]interface{}{"exp": 123}, nil)
	mockService.On("Login", mock.Anything, "test", "test").Return(expectedUser, nil)

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response response_dto.AppUserLoginResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, "test", response.AppUserDto.Username)
}

func TestLoginInvalidCredentials(t *testing.T) {
	mockService := new(MockAppUserService)
	handler := transportHttp.Handler{AppUserService: mockService}

	loginDto := request_dto.AppUserLoginRequestDto{Username: "test", Password: "wrong"}
	loginDtoBytes, _ := json.Marshal(loginDto)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginDtoBytes))

	mockService.On("Login", mock.Anything, "test", "wrong").Return(repository.AppUser{}, errors.New("invalid credentials"))

	rr := httptest.NewRecorder()
	handler.Login(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestTokenRefreshSuccessful(t *testing.T) {
	mockService := new(MockAppUserService)
	handler := transportHttp.Handler{AppUserService: mockService}

	req, _ := http.NewRequest("POST", "/auth/token-refresh", nil)
	req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: "validRefreshToken"})

	expectedUser := repository.AppUser{ID: uuid.New(), Username: "test"}
	mockService.On("RefreshToken", mock.Anything, "validRefreshToken").Return("newAccessToken", make(map[string]interface{}), expectedUser, nil)

	rr := httptest.NewRecorder()
	handler.TokenRefresh(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response response_dto.AppUserLoginResponse
	_ = json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, "test", response.AppUserDto.Username)
	assert.Equal(t, "newAccessToken", response.AccessToken)
}

func TestTokenRefreshInvalidToken(t *testing.T) {
	mockService := new(MockAppUserService)
	handler := transportHttp.Handler{AppUserService: mockService}

	req, _ := http.NewRequest("POST", "/auth/token-refresh", nil)
	req.AddCookie(&http.Cookie{Name: refreshTokenCookieName, Value: "invalidRefreshToken"})

	mockService.On("RefreshToken", mock.Anything, "invalidRefreshToken").Return("", make(map[string]interface{}), repository.AppUser{}, errors.New("invalid token"))

	rr := httptest.NewRecorder()
	handler.TokenRefresh(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetAppUserByIdSuccessful(t *testing.T) {
	mockService := new(MockAppUserService)
	handler := transportHttp.Handler{AppUserService: mockService}

	expectedUser := repository.AppUser{ID: uuid.New(), Username: "test"}
	mockService.On("GetAppUserById", mock.Anything, expectedUser.ID).Return(expectedUser, nil)

	req, _ := http.NewRequest("GET", "/users/"+expectedUser.ID.String(), nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handler.GetAppUserById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response response_dto.AppUserDto
	_ = json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Equal(t, "test", response.Username)
}

func TestGetAppUserByIdNotFound(t *testing.T) {
	mockService := new(MockAppUserService)
	handler := transportHttp.Handler{AppUserService: mockService}

	nonExistentUserId := uuid.New()
	mockService.On("GetAppUserById", mock.Anything, nonExistentUserId).Return(repository.AppUser{}, errors.New("not found"))

	req, _ := http.NewRequest("GET", "/users/"+nonExistentUserId.String(), nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", handler.GetAppUserById)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestSendUserEmailVerification_Success(t *testing.T) {
	emailAddress := "test@example.com"
	mockService := new(MockAppUserService)
	mockService.On("SendUserEmailVerification", mock.Anything, emailAddress).Return(nil)

	req, _ := http.NewRequest("POST", "/api/user/send-email-verification/", strings.NewReader(""))
	req = req.WithContext(context.WithValue(req.Context(), "jwt_claims", map[string]interface{}{
		"email_verified": false,
		"email":          emailAddress,
	}))

	recorder := httptest.NewRecorder()
	handler := transportHttp.Handler{AppUserService: mockService}
	handler.SendUserEmailVerification(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	mockService.AssertExpectations(t)
}

func TestSendUserEmailVerification_EmailAlreadyVerified(t *testing.T) {
	emailAddress := "test@example.com"
	mockService := new(MockAppUserService)

	req, _ := http.NewRequest("POST", "/api/user/send-email-verification/", strings.NewReader(""))
	req = req.WithContext(context.WithValue(req.Context(), "jwt_claims", map[string]interface{}{
		"email_verified": true,
		"email":          emailAddress,
	}))

	recorder := httptest.NewRecorder()
	handler := transportHttp.Handler{AppUserService: mockService}
	handler.SendUserEmailVerification(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	mockService.AssertExpectations(t)
}
