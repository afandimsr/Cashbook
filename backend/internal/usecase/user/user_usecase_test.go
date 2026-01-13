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

func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
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
