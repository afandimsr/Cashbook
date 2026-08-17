package category_test

import (
	"errors"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/domain/category"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/category"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) FindAllByUserID(userID int64) ([]category.Category, error) {
	args := m.Called(userID)
	return args.Get(0).([]category.Category), args.Error(1)
}

func (m *MockCategoryRepository) FindByID(id int64) (category.Category, error) {
	args := m.Called(id)
	return args.Get(0).(category.Category), args.Error(1)
}

func (m *MockCategoryRepository) Save(c *category.Category) error {
	return m.Called(c).Error(0)
}

func (m *MockCategoryRepository) Update(c *category.Category) error {
	return m.Called(c).Error(0)
}

func (m *MockCategoryRepository) Delete(id int64) error {
	return m.Called(id).Error(0)
}

func TestCategory_GetAllByUserID(t *testing.T) {
	repo := new(MockCategoryRepository)
	u := uc.New(repo)

	expected := []category.Category{{ID: 1, Name: "Food"}}
	repo.On("FindAllByUserID", int64(9)).Return(expected, nil).Once()

	got, err := u.GetAllByUserID(9)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	repo.AssertExpectations(t)
}

func TestCategory_GetByID(t *testing.T) {
	repo := new(MockCategoryRepository)
	u := uc.New(repo)

	repo.On("FindByID", int64(3)).Return(category.Category{ID: 3, Name: "Bills"}, nil).Once()

	got, err := u.GetByID(3)
	assert.NoError(t, err)
	assert.Equal(t, "Bills", got.Name)
	repo.AssertExpectations(t)
}

func TestCategory_Create(t *testing.T) {
	repo := new(MockCategoryRepository)
	u := uc.New(repo)

	repo.On("Save", mock.MatchedBy(func(c *category.Category) bool {
		return c.Name == "Transport"
	})).Return(nil).Once()

	err := u.Create(category.Category{Name: "Transport"})
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCategory_Update_MergesFields(t *testing.T) {
	repo := new(MockCategoryRepository)
	u := uc.New(repo)

	existing := category.Category{ID: 5, Name: "Old", Type: "expense", Color: "#000", Icon: "x"}
	repo.On("FindByID", int64(5)).Return(existing, nil).Once()
	repo.On("Update", mock.MatchedBy(func(c *category.Category) bool {
		return c.ID == 5 && c.Name == "New" && c.Type == "income" && c.Color == "#fff" && c.Icon == "y"
	})).Return(nil).Once()

	err := u.Update(5, category.Category{Name: "New", Type: "income", Color: "#fff", Icon: "y"})
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCategory_Update_FindError(t *testing.T) {
	repo := new(MockCategoryRepository)
	u := uc.New(repo)

	repo.On("FindByID", int64(5)).Return(category.Category{}, errors.New("not found")).Once()

	err := u.Update(5, category.Category{Name: "New"})
	assert.Error(t, err)
	repo.AssertNotCalled(t, "Update", mock.Anything)
	repo.AssertExpectations(t)
}

func TestCategory_Delete(t *testing.T) {
	repo := new(MockCategoryRepository)
	u := uc.New(repo)

	repo.On("Delete", int64(2)).Return(nil).Once()

	assert.NoError(t, u.Delete(2))
	repo.AssertExpectations(t)
}
