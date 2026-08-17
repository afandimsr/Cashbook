package handler_test

import (
	"net/http"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/afandimsr/cashbook-backend/internal/domain/budget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBudgetUsecase struct{ mock.Mock }

func (m *MockBudgetUsecase) GetBudgets(userID int64, month, year int) ([]budget.Budget, error) {
	args := m.Called(userID, month, year)
	return args.Get(0).([]budget.Budget), args.Error(1)
}
func (m *MockBudgetUsecase) SetBudget(userID int64, b budget.Budget) error {
	return m.Called(userID, b).Error(0)
}
func (m *MockBudgetUsecase) GetBudgetByCategory(userID, categoryID int64, month, year int) (budget.Budget, error) {
	args := m.Called(userID, categoryID, month, year)
	return args.Get(0).(budget.Budget), args.Error(1)
}

func TestBudgetHandler_GetBudgets(t *testing.T) {
	uc := new(MockBudgetUsecase)
	h := handler.NewBudgetHandler(uc)
	uc.On("GetBudgets", int64(1), 5, 2026).Return([]budget.Budget{{ID: 1, Amount: 100}}, nil).Once()

	r := newRouter()
	r.GET("/budgets", h.GetBudgets)

	w := doJSON(r, http.MethodGet, "/budgets?month=5&year=2026", "")
	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestBudgetHandler_SetBudget(t *testing.T) {
	uc := new(MockBudgetUsecase)
	h := handler.NewBudgetHandler(uc)
	uc.On("SetBudget", int64(1), mock.MatchedBy(func(b budget.Budget) bool {
		return b.CategoryID == 3 && b.Amount == 900
	})).Return(nil).Once()

	r := newRouter()
	r.POST("/budgets", h.SetBudget)

	w := doJSON(r, http.MethodPost, "/budgets", `{"category_id":3,"amount":900,"month":5,"year":2026}`)
	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestBudgetHandler_SetBudget_BadJSON(t *testing.T) {
	uc := new(MockBudgetUsecase)
	h := handler.NewBudgetHandler(uc)

	r := newRouter()
	r.POST("/budgets", h.SetBudget)

	w := doJSON(r, http.MethodPost, "/budgets", `{bad`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "SetBudget", mock.Anything, mock.Anything)
}
