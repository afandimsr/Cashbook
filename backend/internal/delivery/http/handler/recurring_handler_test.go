package handler_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	rt "github.com/afandimsr/cashbook-backend/internal/domain/recurring_transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRecurringUsecase struct{ mock.Mock }

func (m *MockRecurringUsecase) GetRecurring(userID int64) ([]rt.RecurringTransaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]rt.RecurringTransaction), args.Error(1)
}
func (m *MockRecurringUsecase) CreateRecurring(userID int64, r rt.RecurringTransaction) (rt.RecurringTransaction, error) {
	args := m.Called(userID, r)
	return args.Get(0).(rt.RecurringTransaction), args.Error(1)
}
func (m *MockRecurringUsecase) DeleteRecurring(id int64) error   { return m.Called(id).Error(0) }
func (m *MockRecurringUsecase) ProcessDueTransactions() error    { return m.Called().Error(0) }

func TestRecurringHandler_GetRecurring(t *testing.T) {
	uc := new(MockRecurringUsecase)
	h := handler.NewRecurringHandler(uc)
	uc.On("GetRecurring", int64(1)).Return([]rt.RecurringTransaction{{ID: 1}}, nil).Once()

	r := newRouter()
	r.GET("/recurring", h.GetRecurring)

	w := doJSON(r, http.MethodGet, "/recurring", "")
	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestRecurringHandler_Create(t *testing.T) {
	uc := new(MockRecurringUsecase)
	h := handler.NewRecurringHandler(uc)
	uc.On("CreateRecurring", int64(1), mock.MatchedBy(func(r rt.RecurringTransaction) bool {
		return r.Amount == 100
	})).Return(rt.RecurringTransaction{ID: 9, Amount: 100}, nil).Once()

	r := newRouter()
	r.POST("/recurring", h.CreateRecurring)

	w := doJSON(r, http.MethodPost, "/recurring", `{"amount":100,"type":"expense","frequency":"monthly","category_id":2}`)
	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)
}

func TestRecurringHandler_Create_Error(t *testing.T) {
	uc := new(MockRecurringUsecase)
	h := handler.NewRecurringHandler(uc)
	uc.On("CreateRecurring", int64(1), mock.Anything).
		Return(rt.RecurringTransaction{}, errors.New("db")).Once()

	r := newRouter()
	r.POST("/recurring", h.CreateRecurring)

	w := doJSON(r, http.MethodPost, "/recurring", `{"amount":100}`)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRecurringHandler_Delete_InvalidID(t *testing.T) {
	uc := new(MockRecurringUsecase)
	h := handler.NewRecurringHandler(uc)

	r := newRouter()
	r.DELETE("/recurring/:id", h.DeleteRecurring)

	w := doJSON(r, http.MethodDelete, "/recurring/xx", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecurringHandler_ProcessDue(t *testing.T) {
	uc := new(MockRecurringUsecase)
	h := handler.NewRecurringHandler(uc)
	uc.On("ProcessDueTransactions").Return(nil).Once()

	r := newRouter()
	r.POST("/recurring/process", h.ProcessDue)

	w := doJSON(r, http.MethodPost, "/recurring/process", "")
	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestRecurringHandler_ProcessDue_Error(t *testing.T) {
	uc := new(MockRecurringUsecase)
	h := handler.NewRecurringHandler(uc)
	uc.On("ProcessDueTransactions").Return(errors.New("db")).Once()

	r := newRouter()
	r.POST("/recurring/process", h.ProcessDue)

	w := doJSON(r, http.MethodPost, "/recurring/process", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
