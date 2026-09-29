package service_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/afandimsr/cashbook-bot/internal/llm"
	"github.com/afandimsr/cashbook-bot/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBackend struct{ mock.Mock }

func (m *MockBackend) GetContext(ctx context.Context, chatID int64) (*backend.Context, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*backend.Context), args.Error(1)
}
func (m *MockBackend) LinkAccount(ctx context.Context, chatID int64, username, code string) error {
	return m.Called(ctx, chatID, username, code).Error(0)
}
func (m *MockBackend) CreateTransaction(ctx context.Context, req backend.CreateTransactionRequest) error {
	return m.Called(ctx, req).Error(0)
}
func (m *MockBackend) GetSummary(ctx context.Context, chatID int64) (*backend.DashboardSummary, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*backend.DashboardSummary), args.Error(1)
}
func (m *MockBackend) GetMonthlyReport(ctx context.Context, chatID int64, month, year int) (*backend.MonthlyReport, error) {
	args := m.Called(ctx, chatID, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*backend.MonthlyReport), args.Error(1)
}
func (m *MockBackend) GetBudgetStatus(ctx context.Context, chatID int64, month, year int) ([]backend.BudgetStatus, error) {
	args := m.Called(ctx, chatID, month, year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]backend.BudgetStatus), args.Error(1)
}
func (m *MockBackend) ListLinks(ctx context.Context) ([]backend.TelegramLink, error) {
	args := m.Called(ctx)
	return args.Get(0).([]backend.TelegramLink), args.Error(1)
}

type MockLLM struct{ mock.Mock }

func (m *MockLLM) ParseTransaction(ctx context.Context, message string, categories []backend.Category) (*llm.ParsedTransaction, error) {
	args := m.Called(ctx, message, categories)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*llm.ParsedTransaction), args.Error(1)
}

var testCategories = []backend.Category{
	{ID: 1, Name: "Makanan", Type: "expense"},
	{ID: 2, Name: "Gaji", Type: "income"},
}

func TestHandleMessage_Success(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetContext", mock.Anything, int64(555)).
		Return(&backend.Context{UserID: 7, Categories: testCategories}, nil).Once()
	ai.On("ParseTransaction", mock.Anything, "Beli bensin 50rb", testCategories).
		Return(&llm.ParsedTransaction{Amount: 50000, Type: "expense", CategoryName: "makanan", Note: "bensin"}, nil).Once()
	be.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(r backend.CreateTransactionRequest) bool {
		return r.ChatID == 555 && r.CategoryID == 1 && r.Amount == 50000 && r.Type == "expense"
	})).Return(nil).Once()

	reply, err := svc.HandleMessage(context.Background(), 555, "Beli bensin 50rb")
	assert.NoError(t, err)
	assert.Contains(t, reply, "Pengeluaran")
	assert.Contains(t, reply, "Makanan")
	assert.Contains(t, reply, "Rp50.000")
	assert.NotContains(t, reply, "(-)")
	be.AssertExpectations(t)
	ai.AssertExpectations(t)
}

func TestHandleMessage_NotLinked(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetContext", mock.Anything, int64(999)).
		Return(nil, &backend.APIError{StatusCode: 404}).Once()

	reply, err := svc.HandleMessage(context.Background(), 999, "Beli bensin 50rb")
	assert.NoError(t, err)
	assert.Contains(t, reply, "belum terhubung")
	ai.AssertNotCalled(t, "ParseTransaction", mock.Anything, mock.Anything, mock.Anything)
}

func TestHandleMessage_CategoryMismatch(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetContext", mock.Anything, int64(555)).
		Return(&backend.Context{UserID: 7, Categories: testCategories}, nil).Once()
	ai.On("ParseTransaction", mock.Anything, mock.Anything, mock.Anything).
		Return(&llm.ParsedTransaction{Amount: 20000, Type: "expense", CategoryName: "Hiburan"}, nil).Once()

	reply, err := svc.HandleMessage(context.Background(), 555, "Nonton 20rb")
	assert.NoError(t, err)
	assert.Contains(t, reply, "tidak cocok")
	assert.Contains(t, reply, "Rp20.000")
	be.AssertNotCalled(t, "CreateTransaction", mock.Anything, mock.Anything)
}

func TestHandleMessage_NoParse(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetContext", mock.Anything, int64(555)).
		Return(&backend.Context{UserID: 7, Categories: testCategories}, nil).Once()
	ai.On("ParseTransaction", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, llm.ErrNoFunctionCall).Once()

	reply, err := svc.HandleMessage(context.Background(), 555, "besok jam 5 ketemu")
	assert.NoError(t, err)
	assert.Contains(t, reply, "tidak bisa memahami")
}

func TestWhoAmI_Linked(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetContext", mock.Anything, int64(555)).
		Return(&backend.Context{UserID: 7, Name: "Afandi", Email: "afandi@example.com", Categories: testCategories}, nil).Once()

	reply, err := svc.WhoAmI(context.Background(), 555, "Mohamad", "Afandi")
	assert.NoError(t, err)
	assert.Contains(t, reply, "Halo Mohamad Afandi!")
	assert.Contains(t, reply, "afandi@example.com")
}

func TestWhoAmI_NotLinked(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetContext", mock.Anything, int64(999)).
		Return(nil, &backend.APIError{StatusCode: 404}).Once()

	reply, err := svc.WhoAmI(context.Background(), 999, "Budi", "")
	assert.NoError(t, err)
	assert.Contains(t, reply, "Halo Budi!") // no trailing space when last name is empty
	assert.Contains(t, reply, "belum terhubung")
}

func TestLink_InvalidCode(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("LinkAccount", mock.Anything, int64(555), "afandi", "BAD").
		Return(&backend.APIError{StatusCode: 400}).Once()

	reply, err := svc.Link(context.Background(), 555, "afandi", "BAD")
	assert.NoError(t, err)
	assert.Contains(t, reply, "tidak valid")
}

func TestMonthlyReport_Success(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetMonthlyReport", mock.Anything, int64(555), 8, 2026).
		Return(&backend.MonthlyReport{
			Income:  []backend.CategoryReport{{CategoryName: "Gaji", TotalAmount: 5000000}},
			Expense: []backend.CategoryReport{{CategoryName: "Makanan", TotalAmount: 100000}},
		}, nil).Once()

	reply, err := svc.MonthlyReport(context.Background(), 555, 8, 2026)
	assert.NoError(t, err)
	assert.Contains(t, reply, "- Gaji: Rp5.000.000")
	assert.Contains(t, reply, "Total pemasukan: Rp5.000.000")
	assert.Contains(t, reply, "- Makanan: Rp100.000")
	assert.Contains(t, reply, "Total pengeluaran: Rp100.000")
	assert.Contains(t, reply, "Selisih: Rp4.900.000")
}

func TestMonthlyReport_OnlyExpense(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetMonthlyReport", mock.Anything, int64(555), 8, 2026).
		Return(&backend.MonthlyReport{
			Income:  []backend.CategoryReport{},
			Expense: []backend.CategoryReport{{CategoryName: "Makanan", TotalAmount: 100000}},
		}, nil).Once()

	reply, err := svc.MonthlyReport(context.Background(), 555, 8, 2026)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Pemasukan:\n- tidak ada")
	assert.Contains(t, reply, "Selisih: -Rp100.000")
}

func TestMonthlyReport_Empty(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetMonthlyReport", mock.Anything, int64(555), 8, 2026).
		Return(&backend.MonthlyReport{}, nil).Once()

	reply, err := svc.MonthlyReport(context.Background(), 555, 8, 2026)
	assert.NoError(t, err)
	assert.Contains(t, reply, "tidak ada transaksi tercatat")
}

func TestMonthlyReport_NotLinked(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetMonthlyReport", mock.Anything, int64(999), 8, 2026).
		Return(nil, &backend.APIError{StatusCode: 404}).Once()

	reply, err := svc.MonthlyReport(context.Background(), 999, 8, 2026)
	assert.NoError(t, err)
	assert.Contains(t, reply, "belum terhubung")
}

func TestSummary_Success(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetSummary", mock.Anything, int64(555)).
		Return(&backend.DashboardSummary{TotalIncome: 100, TotalExpense: 40, Balance: 60}, nil).Once()

	reply, err := svc.Summary(context.Background(), 555)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Saldo")
	assert.Contains(t, reply, "Rp100")
	assert.Contains(t, reply, "Rp60")
}

func TestHandleMessage_NonTransaction_AnsweredLocally(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	reply, err := svc.HandleMessage(context.Background(), 555, "halo, makasih ya")
	assert.NoError(t, err)
	assert.Equal(t, service.HelpMessage, reply)
	be.AssertNotCalled(t, "GetContext", mock.Anything, mock.Anything)
	ai.AssertNotCalled(t, "ParseTransaction", mock.Anything, mock.Anything, mock.Anything)
}

func TestHandleMessage_AllProvidersDown_IsServerError(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetContext", mock.Anything, int64(555)).
		Return(&backend.Context{UserID: 7, Categories: testCategories}, nil).Once()
	ai.On("ParseTransaction", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("all llm providers failed: gemini: status 429")).Once()

	reply, err := svc.HandleMessage(context.Background(), 555, "beli kopi 25rb")
	assert.Error(t, err)
	assert.Empty(t, reply)
	be.AssertNotCalled(t, "CreateTransaction", mock.Anything, mock.Anything)
}

func TestHandleMessage_InvalidOutput_AsksToRephrase(t *testing.T) {
	be := new(MockBackend)
	ai := new(MockLLM)
	svc := service.New(be, ai)

	be.On("GetContext", mock.Anything, int64(555)).
		Return(&backend.Context{UserID: 7, Categories: testCategories}, nil).Once()
	ai.On("ParseTransaction", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, fmt.Errorf("all llm providers failed: %w", llm.ErrInvalidOutput)).Once()

	reply, err := svc.HandleMessage(context.Background(), 555, "beli kopi 25rb")
	assert.NoError(t, err)
	assert.Contains(t, reply, "tidak bisa memahami")
}

func TestCategories_GroupsAndSorts(t *testing.T) {
	be := new(MockBackend)
	svc := service.New(be, new(MockLLM))

	be.On("GetContext", mock.Anything, int64(555)).Return(&backend.Context{Categories: []backend.Category{
		{Name: "transportasi", Type: "expense"},
		{Name: "Gaji", Type: "income"},
		{Name: "Makanan", Type: "expense"},
		{Name: "Bonus", Type: "income"},
	}}, nil).Once()

	reply, err := svc.Categories(context.Background(), 555)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Pemasukan:\n- Bonus\n- Gaji\n")
	assert.Contains(t, reply, "Pengeluaran:\n- Makanan\n- transportasi\n")
	assert.Less(t, strings.Index(reply, "Pemasukan:"), strings.Index(reply, "Pengeluaran:"))
}

func TestCategories_OnlyExpense(t *testing.T) {
	be := new(MockBackend)
	svc := service.New(be, new(MockLLM))

	be.On("GetContext", mock.Anything, int64(555)).Return(&backend.Context{Categories: []backend.Category{
		{Name: "Makanan", Type: "expense"},
	}}, nil).Once()

	reply, err := svc.Categories(context.Background(), 555)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Pemasukan:\n- tidak ada")
	assert.Contains(t, reply, "- Makanan")
}

func TestCategories_Empty(t *testing.T) {
	be := new(MockBackend)
	svc := service.New(be, new(MockLLM))

	be.On("GetContext", mock.Anything, int64(555)).Return(&backend.Context{}, nil).Once()

	reply, err := svc.Categories(context.Background(), 555)
	assert.NoError(t, err)
	assert.Contains(t, reply, "belum punya kategori")
}

func TestCategories_NotLinked(t *testing.T) {
	be := new(MockBackend)
	svc := service.New(be, new(MockLLM))

	be.On("GetContext", mock.Anything, int64(999)).
		Return(nil, &backend.APIError{StatusCode: 404}).Once()

	reply, err := svc.Categories(context.Background(), 999)
	assert.NoError(t, err)
	assert.Contains(t, reply, "belum terhubung")
}

func TestBudgets_ShowsUsageAndFlagsOver(t *testing.T) {
	be := new(MockBackend)
	svc := service.New(be, new(MockLLM))

	be.On("GetBudgetStatus", mock.Anything, int64(555), 9, 2026).Return([]backend.BudgetStatus{
		{CategoryName: "Hiburan", Limit: 300000, Spent: 360000, Remaining: -60000, Percentage: 120, IsOver: true},
		{CategoryName: "Transportasi", Limit: 200000, Spent: 170000, Remaining: 30000, Percentage: 85},
		{CategoryName: "Makanan", Limit: 500000, Spent: 250000, Remaining: 250000, Percentage: 50},
	}, nil).Once()

	reply, err := svc.Budgets(context.Background(), 555, 9, 2026)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Anggaran September 2026")
	assert.Contains(t, reply, "🔴 Hiburan: Rp360.000 / Rp300.000 (120%)\n   lewat Rp60.000")
	assert.Contains(t, reply, "🟡 Transportasi: Rp170.000 / Rp200.000 (85%)\n   sisa Rp30.000")
	assert.Contains(t, reply, "🟢 Makanan: Rp250.000 / Rp500.000 (50%)\n   sisa Rp250.000")
	assert.Contains(t, reply, "Total terpakai: Rp780.000 dari Rp1.000.000")
}

func TestBudgets_Empty(t *testing.T) {
	be := new(MockBackend)
	svc := service.New(be, new(MockLLM))

	be.On("GetBudgetStatus", mock.Anything, int64(555), 9, 2026).Return([]backend.BudgetStatus{}, nil).Once()

	reply, err := svc.Budgets(context.Background(), 555, 9, 2026)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Belum ada anggaran untuk September 2026")
}

func TestBudgets_NotLinked(t *testing.T) {
	be := new(MockBackend)
	svc := service.New(be, new(MockLLM))

	be.On("GetBudgetStatus", mock.Anything, int64(999), 9, 2026).
		Return(nil, &backend.APIError{StatusCode: 404}).Once()

	reply, err := svc.Budgets(context.Background(), 999, 9, 2026)
	assert.NoError(t, err)
	assert.Contains(t, reply, "belum terhubung")
}
