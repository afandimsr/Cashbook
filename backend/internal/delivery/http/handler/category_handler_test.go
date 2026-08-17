package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/afandimsr/cashbook-backend/internal/domain/category"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryUsecase struct{ mock.Mock }

func (m *MockCategoryUsecase) GetAllByUserID(userID int64) ([]category.Category, error) {
	args := m.Called(userID)
	return args.Get(0).([]category.Category), args.Error(1)
}
func (m *MockCategoryUsecase) GetByID(id int64) (category.Category, error) {
	args := m.Called(id)
	return args.Get(0).(category.Category), args.Error(1)
}
func (m *MockCategoryUsecase) Create(c category.Category) error      { return m.Called(c).Error(0) }
func (m *MockCategoryUsecase) Update(id int64, c category.Category) error {
	return m.Called(id, c).Error(0)
}
func (m *MockCategoryUsecase) Delete(id int64) error { return m.Called(id).Error(0) }

func TestCategoryHandler_GetCategories(t *testing.T) {
	uc := new(MockCategoryUsecase)
	h := handler.NewCategoryHandler(uc)
	uc.On("GetAllByUserID", int64(1)).Return([]category.Category{{ID: 1, Name: "Food"}}, nil).Once()

	r := newRouter()
	r.GET("/categories", h.GetCategories)

	w := doJSON(r, http.MethodGet, "/categories", "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Food")
	uc.AssertExpectations(t)
}

func TestCategoryHandler_Create(t *testing.T) {
	uc := new(MockCategoryUsecase)
	h := handler.NewCategoryHandler(uc)
	uc.On("Create", mock.MatchedBy(func(c category.Category) bool {
		return c.Name == "Transport" && c.UserID == 1
	})).Return(nil).Once()

	r := newRouter()
	r.POST("/categories", h.CreateCategory)

	w := doJSON(r, http.MethodPost, "/categories", `{"name":"Transport","type":"expense"}`)
	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)
}

func TestCategoryHandler_Create_BadJSON(t *testing.T) {
	uc := new(MockCategoryUsecase)
	h := handler.NewCategoryHandler(uc)

	r := newRouter()
	r.POST("/categories", h.CreateCategory)

	w := doJSON(r, http.MethodPost, "/categories", `{bad`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Create", mock.Anything)
}

func TestCategoryHandler_Update_InvalidID(t *testing.T) {
	uc := new(MockCategoryUsecase)
	h := handler.NewCategoryHandler(uc)

	r := newRouter()
	r.PUT("/categories/:id", h.UpdateCategory)

	w := doJSON(r, http.MethodPut, "/categories/abc", `{"name":"x"}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCategoryHandler_Delete(t *testing.T) {
	uc := new(MockCategoryUsecase)
	h := handler.NewCategoryHandler(uc)
	uc.On("Delete", int64(7)).Return(nil).Once()

	r := newRouter()
	r.DELETE("/categories/:id", h.DeleteCategory)

	w := doJSON(r, http.MethodDelete, "/categories/7", "")
	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestCategoryHandler_GetCategories_Error(t *testing.T) {
	uc := new(MockCategoryUsecase)
	h := handler.NewCategoryHandler(uc)
	uc.On("GetAllByUserID", int64(1)).Return([]category.Category{}, errors.New("db")).Once()

	r := newRouter()
	r.GET("/categories", h.GetCategories)

	w := doJSON(r, http.MethodGet, "/categories", "")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
