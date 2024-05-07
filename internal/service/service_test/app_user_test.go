package service_test

import (
	"context"
	"database/sql"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/service"
	"eau-de-go/pkg/email_util"
	"eau-de-go/pkg/jwt_util"
	"eau-de-go/pkg/password_util"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
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

func (m *MockAppUserStore) SetUserEmailVerified(ctx context.Context, userId uuid.UUID) (repository.AppUser, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserStore) SetUserEmailUnverified(ctx context.Context, userId uuid.UUID) (repository.AppUser, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

func (m *MockAppUserStore) UpdateAppUserLastLoginNow(ctx context.Context, userId uuid.UUID) (repository.AppUser, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(repository.AppUser), args.Error(1)
}

type MockEmailVerifier struct {
	mock.Mock
}

func (m *MockEmailVerifier) CreateToken(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

func (m *MockEmailVerifier) VerifyToken(tokenString string) (string, error) {
	args := m.Called(tokenString)
	return args.String(0), args.Error(1)
}

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) makeMailBytes(toAddresses []string, mailSubject string, mailMessage string) []byte {
	args := m.Called(toAddresses, mailSubject, mailMessage)
	return args.Get(0).([]byte)
}

func (m *MockEmailSender) SendSingleEmail(recipientEmail string, mailSubject string, mailBody string) error {
	args := m.Called(recipientEmail, mailSubject, mailBody)
	return args.Error(0)
}

func (m *MockEmailSender) SendMassEmail(recipientEmails []string, mailSubject string, mailBody string) error {
	args := m.Called(recipientEmails, mailSubject, mailBody)
	return args.Error(0)
}

type MockJwtUtil struct {
	mock.Mock
}

func (m *MockJwtUtil) DecodeToken(tokenType jwt_util.TokenType, token string) (map[string]interface{}, error) {
	args := m.Called(tokenType, token)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockJwtUtil) CreateAccessToken(claims map[string]interface{}) (string, map[string]interface{}, error) {
	args := m.Called(claims)
	return args.String(0), args.Get(1).(map[string]interface{}), args.Error(2)
}

func (m *MockJwtUtil) CreateRefreshToken(claims map[string]interface{}) (string, map[string]interface{}, error) {
	args := m.Called(claims)
	return args.String(0), args.Get(1).(map[string]interface{}), args.Error(2)
}

func (m *MockJwtUtil) CopyTokenClaims(claims map[string]interface{}) map[string]interface{} {
	args := m.Called(claims)
	return args.Get(0).(map[string]interface{})
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

func TestLoginWithValidCredentials(t *testing.T) {
	mockStore := new(MockAppUserStore)
	s := service.AppUserService{AppUserStore: mockStore}
	username := "user"
	password := "P4ssword!123"
	passwordHash, err := password_util.HashPassword(password)
	assert.NoError(t, err)

	user := repository.AppUser{ID: uuid.New(), Username: username, Password: string(passwordHash), IsActive: true}
	mockStore.On("GetAppUserByUsername", mock.Anything, username).Return(user, nil)
	mockStore.On("UpdateAppUserLastLoginNow", mock.Anything, user.ID).Return(user, nil)

	result, err := s.Login(context.Background(), username, password)

	assert.NoError(t, err)
	assert.Equal(t, username, result.Username)
	mockStore.AssertExpectations(t)
}

func TestLoginWithInvalidUsername(t *testing.T) {
	mockStore := new(MockAppUserStore)
	s := service.AppUserService{AppUserStore: mockStore}

	mockStore.On("GetAppUserByUsername", mock.Anything, "invalid").Return(repository.AppUser{}, errors.New("user not found"))

	_, err := s.Login(context.Background(), "invalid", "password")

	assert.Error(t, err)
	mockStore.AssertExpectations(t)
	mockStore.AssertNotCalled(t, "UpdateAppUserLastLoginNow")
}

func TestLoginWithInvalidPassword(t *testing.T) {
	mockStore := new(MockAppUserStore)
	s := service.AppUserService{AppUserStore: mockStore}
	username := "user"
	password := "P4ssword!123"
	passwordHash, err := password_util.HashPassword(password)
	assert.NoError(t, err)

	user := repository.AppUser{ID: uuid.New(), Username: username, Password: string(passwordHash)}
	mockStore.On("GetAppUserByUsername", mock.Anything, username).Return(user, nil)

	result, err := s.Login(context.Background(), username, "wrong password")

	assert.Empty(t, result)
	assert.Error(t, err)
	mockStore.AssertExpectations(t)
	mockStore.AssertNotCalled(t, "UpdateAppUserLastLoginNow")
}

func TestInactiveUserCannotLogIn(t *testing.T) {
	mockStore := new(MockAppUserStore)
	s := service.AppUserService{AppUserStore: mockStore}
	username := "user"
	password := "P4ssword!123"
	passwordHash, err := password_util.HashPassword(password)
	assert.NoError(t, err)

	user := repository.AppUser{ID: uuid.New(), Username: username, Password: string(passwordHash), IsActive: false}
	mockStore.On("GetAppUserByUsername", mock.Anything, username).Return(user, nil)
	mockStore.On("UpdateAppUserLastLoginNow", mock.Anything, user.ID).Return(user, nil)

	result, err := s.Login(context.Background(), username, password)

	assert.Error(t, err)
	assert.Empty(t, result)
	mockStore.AssertExpectations(t)
}

func TestRefreshToken(t *testing.T) {
	mockStore := new(MockAppUserStore)
	mockJwtUtil := new(MockJwtUtil)
	s := service.AppUserService{AppUserStore: mockStore, JwtUtil: mockJwtUtil}

	user := repository.AppUser{ID: uuid.New(), Username: "test", IsActive: true}
	mockStore.On("GetAppUserById", mock.Anything, user.ID).Return(user, nil)
	mockStore.On("UpdateAppUserLastLoginNow", mock.Anything, user.ID).Return(user, nil)
	mockJwtUtil.On("DecodeToken", jwt_util.Refresh, "validToken").Return(map[string]interface{}{"id": user.ID.String()}, nil)
	mockJwtUtil.On("CreateAccessToken", mock.Anything).Return("newAccessToken", map[string]interface{}{}, nil)

	_, _, _, err := s.RefreshToken(context.Background(), "validToken")

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
	mockJwtUtil.AssertExpectations(t)
}

func TestRefreshTokenInvalidToken(t *testing.T) {
	mockStore := new(MockAppUserStore)
	mockJwtUtil := new(MockJwtUtil)
	s := service.AppUserService{AppUserStore: mockStore, JwtUtil: mockJwtUtil}

	mockJwtUtil.On("DecodeToken", jwt_util.Refresh, "invalidToken").Return(map[string]interface{}{}, errors.New("invalid token"))
	_, _, _, err := s.RefreshToken(context.Background(), "invalidToken")

	assert.Error(t, err)
	mockStore.AssertExpectations(t)
	mockJwtUtil.AssertExpectations(t)
}

func TestRefreshTokenUserNotFound(t *testing.T) {
	mockStore := new(MockAppUserStore)
	mockJwtUtil := new(MockJwtUtil)
	s := service.AppUserService{AppUserStore: mockStore, JwtUtil: mockJwtUtil}

	userId := uuid.New()
	mockJwtUtil.On("DecodeToken", jwt_util.Refresh, "validToken").Return(map[string]interface{}{"id": userId.String()}, nil)
	mockStore.On("GetAppUserById", mock.Anything, userId).Return(repository.AppUser{}, errors.New("user not found"))

	_, _, _, err := s.RefreshToken(context.Background(), "validToken")

	assert.Error(t, err)
	mockStore.AssertExpectations(t)
	mockJwtUtil.AssertExpectations(t)
}

func TestRefreshTokenInactiveUser(t *testing.T) {
	mockStore := new(MockAppUserStore)
	mockJwtUtil := new(MockJwtUtil)
	s := service.AppUserService{AppUserStore: mockStore, JwtUtil: mockJwtUtil}

	user := repository.AppUser{ID: uuid.New(), Username: "test", IsActive: false}
	mockJwtUtil.On("DecodeToken", jwt_util.Refresh, "validToken").Return(map[string]interface{}{"id": user.ID.String()}, nil)
	mockStore.On("GetAppUserById", mock.Anything, user.ID).Return(user, nil)

	_, _, _, err := s.RefreshToken(context.Background(), "validToken")

	assert.Error(t, err)
	mockStore.AssertExpectations(t)
	mockJwtUtil.AssertExpectations(t)
}

func TestVerifyEmailVerificationToken(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	userEmailAddress := "test@example.com"
	emailVerifier := email_util.NewEmailTokenVerifier()
	token, _ := emailVerifier.CreateToken(userEmailAddress)

	mockStore := new(MockAppUserStore)
	mockStore.On("SetUserEmailVerified", ctx, userId).Return(repository.AppUser{}, nil)

	s := service.NewAppUserService(mockStore)

	_, err := s.VerifyEmailVerificationToken(ctx, userId, userEmailAddress, token)

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
	mockStore.AssertCalled(t, "SetUserEmailVerified", ctx, userId)
}

func TestVerifyEmailVerificationToken_InvalidToken(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	userEmailAddress := "test@example.com"
	token := "bad_token"

	mockStore := new(MockAppUserStore)

	s := service.NewAppUserService(mockStore)

	_, err := s.VerifyEmailVerificationToken(ctx, userId, userEmailAddress, token)

	assert.Error(t, err)
	mockStore.AssertExpectations(t)
	mockStore.AssertNotCalled(t, "SetUserEmailVerified", ctx, userId)
}

func TestVerifyEmailVerificationToken_EmailMismatch(t *testing.T) {
	ctx := context.Background()
	userId := uuid.New()
	userEmailAddress := "test@example.com"
	emailVerifier := email_util.NewEmailTokenVerifier()
	token, _ := emailVerifier.CreateToken("wrong@email.com")

	mockStore := new(MockAppUserStore)

	s := service.NewAppUserService(mockStore)

	_, err := s.VerifyEmailVerificationToken(ctx, userId, userEmailAddress, token)

	assert.Error(t, err)
	mockStore.AssertExpectations(t)
	mockStore.AssertNotCalled(t, "SetUserEmailVerified", ctx, userId)
}

func TestSendEmailVerification_Success(t *testing.T) {
	emailAddress := "test@example.com"
	mockVerifier := new(MockEmailVerifier)
	mockVerifier.On("CreateToken", emailAddress).Return("token", nil)

	mockSender := new(MockEmailSender)
	mockSender.On("SendSingleEmail", emailAddress, "Email Verification", url.QueryEscape("token")).Return(nil)

	s := service.NewAppUserService(nil)

	s.EmailVerifier = mockVerifier
	s.EmailSender = mockSender

	err := s.SendUserEmailVerification(context.Background(), emailAddress)

	assert.NoError(t, err)
	mockVerifier.AssertExpectations(t)
	mockSender.AssertExpectations(t)

}

func TestSendEmailVerification_TokenCreationError(t *testing.T) {
	emailAddress := "test@example.com"
	mockVerifier := new(MockEmailVerifier)
	mockVerifier.On("CreateToken", emailAddress).Return("", errors.New("token creation error"))

	mockSender := new(MockEmailSender)

	s := service.NewAppUserService(nil)
	s.EmailVerifier = mockVerifier
	s.EmailSender = mockSender

	err := s.SendUserEmailVerification(context.Background(), emailAddress)

	assert.Error(t, err)
	mockVerifier.AssertExpectations(t)
}

func TestSendEmailVerification_EmailSendingError(t *testing.T) {
	emailAddress := "test@example.com"
	mockVerifier := new(MockEmailVerifier)
	mockVerifier.On("CreateToken", emailAddress).Return("token", nil)

	mockSender := new(MockEmailSender)
	mockSender.On("SendSingleEmail", emailAddress, "Email Verification", mock.Anything).Return(errors.New("email sending error"))

	s := service.NewAppUserService(nil)
	s.EmailVerifier = mockVerifier
	s.EmailSender = mockSender

	err := s.SendUserEmailVerification(context.Background(), emailAddress)

	assert.Error(t, err)
	mockVerifier.AssertExpectations(t)
	mockSender.AssertExpectations(t)
}
