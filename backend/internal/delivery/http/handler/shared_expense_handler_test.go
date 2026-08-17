package handler_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	se "github.com/afandimsr/cashbook-backend/internal/domain/shared_expense"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSharedExpenseUsecase struct{ mock.Mock }

func (m *MockSharedExpenseUsecase) CreateSplitBill(ctx context.Context, bill se.SplitBill) error {
	return m.Called(ctx, bill).Error(0)
}
func (m *MockSharedExpenseUsecase) GetByID(ctx context.Context, id string) (se.SplitBill, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(se.SplitBill), args.Error(1)
}
func (m *MockSharedExpenseUsecase) GetByUserID(ctx context.Context, userID int64, f se.Filter) ([]se.SplitBill, error) {
	args := m.Called(ctx, userID, f)
	return args.Get(0).([]se.SplitBill), args.Error(1)
}
func (m *MockSharedExpenseUsecase) SettleParticipant(ctx context.Context, billID, participantID string, userID int64) error {
	return m.Called(ctx, billID, participantID, userID).Error(0)
}
func (m *MockSharedExpenseUsecase) GetSummary(ctx context.Context, userID int64) (se.SplitSummary, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(se.SplitSummary), args.Error(1)
}

const validSplitBody = `{"title":"Lunch","category_id":5,"date":"2026-05-01T00:00:00Z","payer_id":1,"split_method":"EQUAL","subtotal":100000,"participants":[{"user_id":1},{"user_id":2}]}`

func TestSharedExpenseHandler_Create(t *testing.T) {
	uc := new(MockSharedExpenseUsecase)
	h := handler.NewSharedExpenseHandler(uc)
	uc.On("CreateSplitBill", mock.Anything, mock.MatchedBy(func(b se.SplitBill) bool {
		return b.Title == "Lunch" && b.CreatorID == 1 && len(b.Participants) == 2
	})).Return(nil).Once()

	r := newRouter()
	r.POST("/splits", h.Create)

	w := doJSON(r, http.MethodPost, "/splits", validSplitBody)
	assert.Equal(t, http.StatusCreated, w.Code)
	uc.AssertExpectations(t)
}

func TestSharedExpenseHandler_Create_ValidationError(t *testing.T) {
	uc := new(MockSharedExpenseUsecase)
	h := handler.NewSharedExpenseHandler(uc)
	// A wrapped ErrInvalidSplit must map to 400 via errors.Is in the handler.
	uc.On("CreateSplitBill", mock.Anything, mock.Anything).
		Return(fmt.Errorf("%w: bad", se.ErrInvalidSplit)).Once()

	r := newRouter()
	r.POST("/splits", h.Create)

	w := doJSON(r, http.MethodPost, "/splits", validSplitBody)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertExpectations(t)
}

func TestSharedExpenseHandler_Create_InternalError(t *testing.T) {
	uc := new(MockSharedExpenseUsecase)
	h := handler.NewSharedExpenseHandler(uc)
	uc.On("CreateSplitBill", mock.Anything, mock.Anything).Return(errors.New("db")).Once()

	r := newRouter()
	r.POST("/splits", h.Create)

	w := doJSON(r, http.MethodPost, "/splits", validSplitBody)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestSharedExpenseHandler_Create_BadBinding(t *testing.T) {
	uc := new(MockSharedExpenseUsecase)
	h := handler.NewSharedExpenseHandler(uc)

	r := newRouter()
	r.POST("/splits", h.Create)

	// Missing title/participants → binding fails.
	w := doJSON(r, http.MethodPost, "/splits", `{"category_id":5}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	uc.AssertNotCalled(t, "CreateSplitBill", mock.Anything, mock.Anything)
}

func TestSharedExpenseHandler_GetAll(t *testing.T) {
	uc := new(MockSharedExpenseUsecase)
	h := handler.NewSharedExpenseHandler(uc)
	uc.On("GetByUserID", mock.Anything, int64(1), mock.Anything).
		Return([]se.SplitBill{{ID: "b1", Title: "Trip"}}, nil).Once()

	r := newRouter()
	r.GET("/splits", h.GetAll)

	w := doJSON(r, http.MethodGet, "/splits", "")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Trip")
	uc.AssertExpectations(t)
}

func TestSharedExpenseHandler_Settle(t *testing.T) {
	uc := new(MockSharedExpenseUsecase)
	h := handler.NewSharedExpenseHandler(uc)
	uc.On("SettleParticipant", mock.Anything, "b1", "p2", int64(1)).Return(nil).Once()

	r := newRouter()
	r.POST("/splits/:id/settle/:participant_id", h.Settle)

	w := doJSON(r, http.MethodPost, "/splits/b1/settle/p2", "")
	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}

func TestSharedExpenseHandler_GetSummary(t *testing.T) {
	uc := new(MockSharedExpenseUsecase)
	h := handler.NewSharedExpenseHandler(uc)
	uc.On("GetSummary", mock.Anything, int64(1)).
		Return(se.SplitSummary{TotalOwedToMe: 100, TotalIOwe: 20}, nil).Once()

	r := newRouter()
	r.GET("/splits/summary", h.GetSummary)

	w := doJSON(r, http.MethodGet, "/splits/summary", "")
	assert.Equal(t, http.StatusOK, w.Code)
	uc.AssertExpectations(t)
}
