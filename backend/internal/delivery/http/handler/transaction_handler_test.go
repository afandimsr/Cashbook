package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionUsecase struct{ mock.Mock }

func (m *MockTransactionUsecase) GetAllByUserID(userID int64, page, limit int, f transaction.Filter) (transaction.PaginatedTransactions, error) {
	args := m.Called(userID, page, limit, f)
	return args.Get(0).(transaction.PaginatedTransactions), args.Error(1)
}
func (m *MockTransactionUsecase) GetByID(id int64) (transaction.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}
func (m *MockTransactionUsecase) Create(t transaction.Transaction) error { return m.Called(t).Error(0) }
func (m *MockTransactionUsecase) Update(id int64, t transaction.Transaction) error {
	return m.Called(id, t).Error(0)
}
func (m *MockTransactionUsecase) Delete(id int64) error { return m.Called(id).Error(0) }
func (m *MockTransactionUsecase) GetDashboardSummary(userID int64) (transaction.DashboardSummary, error) {
	args := m.Called(userID)
	return args.Get(0).(transaction.DashboardSummary), args.Error(1)
}

func TestTransactionHandler_GetTransactions(t *testing.T) {
	uc := new(MockTransactionUsecase)
	h := handler.NewTransactionHandler(uc)
	uc.On("GetAllByUserID", int64(1), 1, 10, mock.Anything).
		Return(transaction.PaginatedTransactions{Total: 0}, nil).Once()

	r := newRouter()
	r.GET("/transactions", h.GetTransactions)

	w := doJSON(r, http.MethodGet, "/transactions", "")
	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestTransactionHandler_Create(t *testing.T) {
	uc := new(MockTransactionUsecase)
	h := handler.NewTransactionHandler(uc)
	uc.On("Create", mock.MatchedBy(func(tx transaction.Transaction) bool {
		return tx.UserID == 1 && tx.Amount == 500 && tx.Type == "expense"
	})).Return(nil).Once()

	r := newRouter()
	r.POST("/transactions", h.CreateTransaction)

	w := doJSON(r, http.MethodPost, "/transactions", `{"amount":500,"type":"expense","date":"2026-05-01","category_id":3}`)
	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)
}

func TestTransactionHandler_Create_BadDate(t *testing.T) {
	uc := new(MockTransactionUsecase)
	h := handler.NewTransactionHandler(uc)

	r := newRouter()
	r.POST("/transactions", h.CreateTransaction)

	w := doJSON(r, http.MethodPost, "/transactions", `{"amount":1,"date":"31-31-2026"}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "Create", mock.Anything)
}

func TestTransactionHandler_Summary(t *testing.T) {
	uc := new(MockTransactionUsecase)
	h := handler.NewTransactionHandler(uc)
	uc.On("GetDashboardSummary", int64(1)).
		Return(transaction.DashboardSummary{TotalIncome: 100, TotalExpense: 40, Balance: 60}, nil).Once()

	r := newRouter()
	r.GET("/transactions/summary", h.GetSummary)

	w := doJSON(r, http.MethodGet, "/transactions/summary", "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "60")
	uc.AssertExpectations(t)
}

func TestTransactionHandler_Delete_InvalidID(t *testing.T) {
	uc := new(MockTransactionUsecase)
	h := handler.NewTransactionHandler(uc)

	r := newRouter()
	r.DELETE("/transactions/:id", h.DeleteTransaction)

	w := doJSON(r, http.MethodDelete, "/transactions/xx", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTransactionHandler_Delete_Error(t *testing.T) {
	uc := new(MockTransactionUsecase)
	h := handler.NewTransactionHandler(uc)
	uc.On("Delete", int64(5)).Return(errors.New("db")).Once()

	r := newRouter()
	r.DELETE("/transactions/:id", h.DeleteTransaction)

	w := doJSON(r, http.MethodDelete, "/transactions/5", "")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
