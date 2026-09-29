package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/afandimsr/cashbook-bot/internal/llm"
	"github.com/afandimsr/cashbook-bot/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBackend struct{ mock.Mock }

func (m *mockBackend) GetContext(ctx context.Context, chatID int64) (*backend.Context, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*backend.Context), args.Error(1)
}
func (m *mockBackend) LinkAccount(ctx context.Context, chatID int64, username, code string) error {
	return m.Called(ctx, chatID, username, code).Error(0)
}
func (m *mockBackend) CreateTransaction(ctx context.Context, req backend.CreateTransactionRequest) error {
	return m.Called(ctx, req).Error(0)
}
func (m *mockBackend) GetSummary(ctx context.Context, chatID int64) (*backend.DashboardSummary, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*backend.DashboardSummary), args.Error(1)
}
func (m *mockBackend) GetMonthlyReport(ctx context.Context, chatID int64, month, year int) (*backend.MonthlyReport, error) {
	args := m.Called(ctx, chatID, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*backend.MonthlyReport), args.Error(1)
}
func (m *mockBackend) GetBudgetStatus(ctx context.Context, chatID int64, month, year int) ([]backend.BudgetStatus, error) {
	args := m.Called(ctx, chatID, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]backend.BudgetStatus), args.Error(1)
}
func (m *mockBackend) ListLinks(ctx context.Context) ([]backend.TelegramLink, error) {
	args := m.Called(ctx)
	return args.Get(0).([]backend.TelegramLink), args.Error(1)
}

type mockLLM struct{ mock.Mock }

func (m *mockLLM) ParseTransaction(ctx context.Context, message string, categories []backend.Category) (*llm.ParsedTransaction, error) {
	args := m.Called(ctx, message, categories)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*llm.ParsedTransaction), args.Error(1)
}

type mockSender struct{ mock.Mock }

func (m *mockSender) SendMessage(chatID int64, text string) error {
	return m.Called(chatID, text).Error(0)
}

func TestShouldRunToday(t *testing.T) {
	assert.True(t, shouldRunToday(time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC)))
	assert.False(t, shouldRunToday(time.Date(2026, 9, 2, 8, 0, 0, 0, time.UTC)))
}

func TestPreviousMonth(t *testing.T) {
	month, year := previousMonth(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, 8, month)
	assert.Equal(t, 2026, year)

	month, year = previousMonth(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.Equal(t, 12, month)
	assert.Equal(t, 2025, year)
}

func TestTick_NotFirstOfMonth_DoesNothing(t *testing.T) {
	be := new(mockBackend)
	ai := new(mockLLM)
	sender := new(mockSender)
	svc := service.New(be, ai)
	s := New(be, svc, sender)
	s.now = func() time.Time { return time.Date(2026, 9, 15, 8, 0, 0, 0, time.UTC) }

	s.tick(context.Background())
	be.AssertNotCalled(t, "ListLinks", mock.Anything)
}

func TestTick_SendsOnce_ThenSkipsAlreadySentChat(t *testing.T) {
	be := new(mockBackend)
	ai := new(mockLLM)
	sender := new(mockSender)
	svc := service.New(be, ai)
	s := New(be, svc, sender)
	s.now = func() time.Time { return time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC) }

	be.On("ListLinks", mock.Anything).
		Return([]backend.TelegramLink{{UserID: 1, TelegramChatID: 555}}, nil).Twice()
	be.On("GetMonthlyReport", mock.Anything, int64(555), 8, 2026).
		Return(&backend.MonthlyReport{Expense: []backend.CategoryReport{{CategoryName: "Makanan", TotalAmount: 100000}}}, nil).Once()
	sender.On("SendMessage", int64(555), mock.Anything).Return(nil).Once()

	s.tick(context.Background()) // sends
	s.tick(context.Background()) // same month key: must skip GetMonthlyReport/SendMessage

	be.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestTick_SendFailure_DoesNotMarkAsSent(t *testing.T) {
	be := new(mockBackend)
	ai := new(mockLLM)
	sender := new(mockSender)
	svc := service.New(be, ai)
	s := New(be, svc, sender)
	s.now = func() time.Time { return time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC) }

	be.On("ListLinks", mock.Anything).
		Return([]backend.TelegramLink{{UserID: 1, TelegramChatID: 555}}, nil).Twice()
	be.On("GetMonthlyReport", mock.Anything, int64(555), 8, 2026).
		Return(&backend.MonthlyReport{Expense: []backend.CategoryReport{{CategoryName: "Makanan", TotalAmount: 100000}}}, nil).Twice()
	sender.On("SendMessage", int64(555), mock.Anything).Return(assert.AnError).Once()
	sender.On("SendMessage", int64(555), mock.Anything).Return(nil).Once()

	s.tick(context.Background()) // send fails, not marked as sent
	s.tick(context.Background()) // retried, succeeds

	be.AssertExpectations(t)
	sender.AssertExpectations(t)
}
