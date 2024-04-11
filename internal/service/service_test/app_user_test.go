package service_test

import (
	"context"
	"database/sql"
	"eau-de-go/internal/repository"
	"eau-de-go/internal/service"
	"eau-de-go/pkg/email_util"
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

func (m *MockAppUserStore) UpdateAppUserLastLogin(ctx context.Context, data repository.UpdateAppUserLastLoginParams) (repository.AppUser, error) {
	args := m.Called(ctx, data)
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
