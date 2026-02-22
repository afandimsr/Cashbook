package user_test

import (
	"errors"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of user.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindAll(limit, offset int) ([]user.User, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id int64) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (user.User, error) {
	args := m.Called(email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepository) FindByGoogleID(googleID string) (user.User, error) {
	args := m.Called(googleID)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockUserRepository) Save(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) Update(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) DisableTOTP(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) EnableTOTP(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateTOTPSecret(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

// MockMFASettingsRepository is a mock implementation of user.MFASettingsRepository
type MockMFASettingsRepository struct {
	mock.Mock
}

func (m *MockMFASettingsRepository) Get() (*user.MFASettings, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.MFASettings), args.Error(1)
}

func (m *MockMFASettingsRepository) Upsert(settings user.MFASettings) error {
	args := m.Called(settings)
	return args.Error(0)
}

// MockAuthService is a mock implementation of user.AuthService
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Login(email, password string) (bool, error) {
	args := m.Called(email, password)
	return args.Bool(0), args.Error(1)
}

func TestGetByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)

	mockUser := user.User{ID: 1, Name: "Test User", Email: "test@example.com"}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByID", int64(1)).Return(mockUser, nil)

		u, err := usecase.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, mockUser.ID, u.ID)
		assert.Equal(t, mockUser.Name, u.Name)

		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.On("FindByID", int64(2)).Return(user.User{}, errors.New("user not found"))

		_, err := usecase.GetByID(2)

		assert.Error(t, err)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)

	t.Run("Success", func(t *testing.T) {
		newUser := user.User{Name: "New User", Email: "new@example.com", Password: "password123"}

		// Expect Save to be called. Note: Password will be hashed, so we might need strict matching or wildcards.
		// For simplicity, we match using mock.AnythingOfType for the user argument or check fields manually if needed.
		// Since Save takes a value, testify matches equality. The password changes due to hashing.
		// We can match using a custom matcher or simply accept any User
		mockRepo.On("Save", mock.AnythingOfType("user.User")).Return(nil)

		err := usecase.Create(newUser)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MissingEmail", func(t *testing.T) {
		err := usecase.Create(user.User{Password: "123"})
		assert.Error(t, err)
	})
}

func TestResetPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)

	t.Run("Success", func(t *testing.T) {
		id := int64(1)
		newPassword := "StrongPass123!"
		mockUser := user.User{ID: id, Name: "Test", Email: "test@example.com"}

		mockRepo.On("FindByID", id).Return(mockUser, nil)
		mockRepo.On("Update", mock.AnythingOfType("user.User")).Return(nil)

		err := usecase.ResetPassword(id, newPassword)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("WeakPassword", func(t *testing.T) {
		err := usecase.ResetPassword(1, "weak")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "at least 8 characters")
	})
}

func TestValidatePassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := uc.New(mockRepo, nil)

	tests := []struct {
		name     string
		password string
		isValid  bool
	}{
		{"Valid", "StrongPass123!", true},
		{"TooShort", "Sh1!", false},
		{"NoUpper", "lowercase123!", false},
		{"NoLower", "UPPERCASE123!", false},
		{"NoNumber", "NoNumberPass!", false},
		{"NoSymbol", "NoSymbolPass123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isValid {
				mockRepo.On("FindByID", int64(1)).Return(user.User{}, nil).Once()
				mockRepo.On("Update", mock.AnythingOfType("user.User")).Return(nil).Once()
			}
			err := usecase.ResetPassword(1, tt.password)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestLoginWith2FA(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockAuth := new(MockAuthService)
	mockMFA := new(MockMFASettingsRepository)
	usecase := uc.New(mockRepo, mockAuth)
	usecase.SetMFASettingsRepo(mockMFA)

	t.Run("Success - 2FA Setup Required (TOTPSecret empty)", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"
		mockUser := user.User{
			ID:          1,
			Name:        "Test User",
			Email:       email,
			Password:    "anypassword",
			IsActive:    true,
			TOTPSecret:  "",
			TOTPEnabled: false,
		}

		mockAuth.On("Login", email, password).Return(true, nil)
		mockRepo.On("FindByEmail", email).Return(mockUser, nil)

		response, err := usecase.Login(email, password)

		assert.NoError(t, err)
		assert.True(t, response.Requires2FA)
		assert.NotEmpty(t, response.TempToken)
		assert.Empty(t, response.Token)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("Success - 2FA Verification Required (TOTPEnabled)", func(t *testing.T) {
		email := "test2@example.com"
		password := "password123"
		mockUser := user.User{
			ID:          2,
			Name:        "Test User 2",
			Email:       email,
			Password:    "anypassword",
			IsActive:    true,
			TOTPSecret:  "JBSWY3DPEHPK3PXP",
			TOTPEnabled: true,
		}

		mockAuth.On("Login", email, password).Return(true, nil)
		mockRepo.On("FindByEmail", email).Return(mockUser, nil)

		response, err := usecase.Login(email, password)

		assert.NoError(t, err)
		assert.True(t, response.Requires2FA)
		assert.NotEmpty(t, response.TempToken)
		assert.Empty(t, response.Token)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("Success - Normal Login (2FA enabled and verified)", func(t *testing.T) {
		email := "test3@example.com"
		password := "password123"
		mockUser := user.User{
			ID:          3,
			Name:        "Test User 3",
			Email:       email,
			Password:    "anypassword",
			IsActive:    true,
			TOTPSecret:  "JBSWY3DPEHPK3PXP",
			TOTPEnabled: true,
		}

		mockAuth.On("Login", email, password).Return(true, nil)
		mockRepo.On("FindByEmail", email).Return(mockUser, nil)

		response, err := usecase.Login(email, password)

		assert.NoError(t, err)
		assert.True(t, response.Requires2FA)
		assert.NotEmpty(t, response.TempToken)
		mockRepo.AssertExpectations(t)
		mockAuth.AssertExpectations(t)
	})

	t.Run("Failure - Invalid Credentials", func(t *testing.T) {
		email := "invalid@example.com"
		password := "wrongpassword"

		mockRepo.On("FindByEmail", email).Return(user.User{}, errors.New("user not found"))

		_, err := usecase.Login(email, password)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid credentials")
		mockRepo.AssertExpectations(t)
	})
}
