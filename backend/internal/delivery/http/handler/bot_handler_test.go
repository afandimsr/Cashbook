package handler_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/afandimsr/cashbook-backend/internal/domain/audit"
	"github.com/afandimsr/cashbook-backend/internal/domain/telegram"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/bot"
	reportUC "github.com/afandimsr/cashbook-backend/internal/usecase/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBotUsecase struct{ mock.Mock }

func (m *MockBotUsecase) GenerateLinkCode(ctx context.Context, userID int64) (telegram.LinkCode, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(telegram.LinkCode), args.Error(1)
}
func (m *MockBotUsecase) LinkAccount(ctx context.Context, chatID int64, username, code string) (telegram.UserTelegramLink, error) {
	args := m.Called(ctx, chatID, username, code)
	return args.Get(0).(telegram.UserTelegramLink), args.Error(1)
}
func (m *MockBotUsecase) ListLinkedChats(ctx context.Context) ([]telegram.UserTelegramLink, error) {
	args := m.Called(ctx)
	return args.Get(0).([]telegram.UserTelegramLink), args.Error(1)
}
func (m *MockBotUsecase) GetLinkStatus(ctx context.Context, userID int64) (uc.LinkStatus, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(uc.LinkStatus), args.Error(1)
}
func (m *MockBotUsecase) Unlink(ctx context.Context, userID int64) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *MockBotUsecase) GetContext(ctx context.Context, chatID int64) (uc.Context, error) {
	args := m.Called(ctx, chatID)
	return args.Get(0).(uc.Context), args.Error(1)
}
func (m *MockBotUsecase) CreateTransaction(ctx context.Context, in uc.CreateTransactionInput) error {
	return m.Called(ctx, in).Error(0)
}
func (m *MockBotUsecase) GetSummary(ctx context.Context, chatID int64) (transaction.DashboardSummary, error) {
	args := m.Called(ctx, chatID)
	return args.Get(0).(transaction.DashboardSummary), args.Error(1)
}
func (m *MockBotUsecase) GetCategorySpending(ctx context.Context, chatID int64, month, year int) ([]reportUC.CategoryReport, error) {
	args := m.Called(ctx, chatID, month, year)
	return args.Get(0).([]reportUC.CategoryReport), args.Error(1)
}

type MockAuditRepo struct{ mock.Mock }

func (m *MockAuditRepo) Save(ctx context.Context, entry *audit.Entry) error {
	return m.Called(ctx, entry).Error(0)
}

func TestBotHandler_GenerateLinkCode(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GenerateLinkCode", mock.Anything, int64(1)).
		Return(telegram.LinkCode{Code: "ABC123"}, nil).Once()

	r := newRouter()
	r.POST("/telegram/link-code", h.GenerateLinkCode)

	w := doJSON(r, http.MethodPost, "/telegram/link-code", "")
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_Link_Success(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("LinkAccount", mock.Anything, int64(555), "afandi", "ABC123").
		Return(telegram.UserTelegramLink{UserID: 1, TelegramChatID: 555}, nil).Once()

	r := newRouter()
	r.POST("/internal/bot/link", h.Link)

	body := `{"chat_id":555,"username":"afandi","code":"ABC123"}`
	w := doJSON(r, http.MethodPost, "/internal/bot/link", body)
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_Link_InvalidCode(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("LinkAccount", mock.Anything, int64(555), "afandi", "BAD").
		Return(telegram.UserTelegramLink{}, telegram.ErrLinkCodeInvalid).Once()

	r := newRouter()
	r.POST("/internal/bot/link", h.Link)

	body := `{"chat_id":555,"username":"afandi","code":"BAD"}`
	w := doJSON(r, http.MethodPost, "/internal/bot/link", body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBotHandler_Link_AlreadyLinkedToAnotherUser(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("LinkAccount", mock.Anything, int64(555), "afandi", "ABC123").
		Return(telegram.UserTelegramLink{}, telegram.ErrChatAlreadyLinked).Once()

	r := newRouter()
	r.POST("/internal/bot/link", h.Link)

	body := `{"chat_id":555,"username":"afandi","code":"ABC123"}`
	w := doJSON(r, http.MethodPost, "/internal/bot/link", body)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestBotHandler_GetLinkStatus(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetLinkStatus", mock.Anything, int64(1)).
		Return(uc.LinkStatus{Linked: true, TelegramUsername: "afandi"}, nil).Once()

	r := newRouter()
	r.GET("/telegram/status", h.GetLinkStatus)

	w := doJSON(r, http.MethodGet, "/telegram/status", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_Unlink(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("Unlink", mock.Anything, int64(1)).Return(nil).Once()

	r := newRouter()
	r.DELETE("/telegram/link", h.Unlink)

	w := doJSON(r, http.MethodDelete, "/telegram/link", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_AdminGenerateLinkCode(t *testing.T) {
	m := new(MockBotUsecase)
	auditRepo := new(MockAuditRepo)
	h := handler.NewBotHandler(m, auditRepo)
	m.On("GenerateLinkCode", mock.Anything, int64(42)).
		Return(telegram.LinkCode{Code: "ABC123"}, nil).Once()
	auditRepo.On("Save", mock.Anything, mock.MatchedBy(func(e *audit.Entry) bool {
		return e.AdminUserID == 1 && e.Action == "telegram.generate_link_code" && e.TargetUserID != nil && *e.TargetUserID == 42
	})).Return(nil).Once()

	r := newRouter()
	r.POST("/users/:id/telegram/link-code", h.AdminGenerateLinkCode)

	w := doJSON(r, http.MethodPost, "/users/42/telegram/link-code", "")
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
	auditRepo.AssertExpectations(t)
}

func TestBotHandler_AdminGenerateLinkCode_UserNotFound(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GenerateLinkCode", mock.Anything, int64(999999)).
		Return(telegram.LinkCode{}, telegram.ErrUserNotFound).Once()

	r := newRouter()
	r.POST("/users/:id/telegram/link-code", h.AdminGenerateLinkCode)

	w := doJSON(r, http.MethodPost, "/users/999999/telegram/link-code", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBotHandler_AdminGetLinkStatus(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetLinkStatus", mock.Anything, int64(42)).
		Return(uc.LinkStatus{Linked: false}, nil).Once()

	r := newRouter()
	r.GET("/users/:id/telegram/status", h.AdminGetLinkStatus)

	w := doJSON(r, http.MethodGet, "/users/42/telegram/status", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_AdminUnlink(t *testing.T) {
	m := new(MockBotUsecase)
	auditRepo := new(MockAuditRepo)
	h := handler.NewBotHandler(m, auditRepo)
	m.On("Unlink", mock.Anything, int64(42)).Return(nil).Once()
	auditRepo.On("Save", mock.Anything, mock.MatchedBy(func(e *audit.Entry) bool {
		return e.AdminUserID == 1 && e.Action == "telegram.unlink" && e.TargetUserID != nil && *e.TargetUserID == 42
	})).Return(nil).Once()

	r := newRouter()
	r.DELETE("/users/:id/telegram/link", h.AdminUnlink)

	w := doJSON(r, http.MethodDelete, "/users/42/telegram/link", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
	auditRepo.AssertExpectations(t)
}

func TestBotHandler_AdminListLinks(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("ListLinkedChats", mock.Anything).
		Return([]telegram.UserTelegramLink{{ID: 1, UserID: 7, TelegramChatID: 555}}, nil).Once()

	r := newRouter()
	r.GET("/admin/telegram/links", h.AdminListLinks)

	w := doJSON(r, http.MethodGet, "/admin/telegram/links", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_ListLinks_Success(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("ListLinkedChats", mock.Anything).
		Return([]telegram.UserTelegramLink{{ID: 1, UserID: 7, TelegramChatID: 555}}, nil).Once()

	r := newRouter()
	r.GET("/internal/bot/links", h.ListLinks)

	w := doJSON(r, http.MethodGet, "/internal/bot/links", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_GetContext_Success(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetContext", mock.Anything, int64(555)).
		Return(uc.Context{UserID: 7}, nil).Once()

	r := newRouter()
	r.GET("/internal/bot/context", h.GetContext)

	req := doJSON(r, http.MethodGet, "/internal/bot/context?chat_id=555", "")
	assert.Equal(t, http.StatusOK, req.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_GetContext_NotLinked(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetContext", mock.Anything, int64(999)).
		Return(uc.Context{}, telegram.ErrLinkNotFound).Once()

	r := newRouter()
	r.GET("/internal/bot/context", h.GetContext)

	w := doJSON(r, http.MethodGet, "/internal/bot/context?chat_id=999", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBotHandler_GetContext_MissingChatID(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))

	r := newRouter()
	r.GET("/internal/bot/context", h.GetContext)

	w := doJSON(r, http.MethodGet, "/internal/bot/context", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "GetContext", mock.Anything, mock.Anything)
}

func TestBotHandler_CreateTransaction_Success(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("CreateTransaction", mock.Anything, mock.MatchedBy(func(in uc.CreateTransactionInput) bool {
		return in.ChatID == 555 && in.CategoryID == 3 && in.Amount == 50000 && in.Type == "expense"
	})).Return(nil).Once()

	r := newRouter()
	r.POST("/internal/bot/transactions", h.CreateTransaction)

	body := `{"chat_id":555,"category_id":3,"amount":50000,"type":"expense","note":"bensin"}`
	w := doJSON(r, http.MethodPost, "/internal/bot/transactions", body)
	assert.Equal(t, http.StatusCreated, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_CreateTransaction_CategoryNotOwned(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("CreateTransaction", mock.Anything, mock.Anything).Return(uc.ErrCategoryNotOwned).Once()

	r := newRouter()
	r.POST("/internal/bot/transactions", h.CreateTransaction)

	body := `{"chat_id":555,"category_id":3,"amount":50000,"type":"expense"}`
	w := doJSON(r, http.MethodPost, "/internal/bot/transactions", body)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestBotHandler_GetSummary_Success(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetSummary", mock.Anything, int64(555)).
		Return(transaction.DashboardSummary{Balance: 100}, nil).Once()

	r := newRouter()
	r.GET("/internal/bot/summary", h.GetSummary)

	w := doJSON(r, http.MethodGet, "/internal/bot/summary?chat_id=555", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_GetCategorySpending_Success(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetCategorySpending", mock.Anything, int64(555), 9, 2026).
		Return([]reportUC.CategoryReport{{CategoryID: 1}}, nil).Once()

	r := newRouter()
	r.GET("/internal/bot/reports/spending", h.GetCategorySpending)

	w := doJSON(r, http.MethodGet, "/internal/bot/reports/spending?chat_id=555&month=9&year=2026", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func (m *MockBotUsecase) GetMonthlyReport(ctx context.Context, chatID int64, month, year int) (reportUC.MonthlyCategoryReport, error) {
	args := m.Called(ctx, chatID, month, year)
	return args.Get(0).(reportUC.MonthlyCategoryReport), args.Error(1)
}

func TestBotHandler_GetMonthlyReport_Success(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetMonthlyReport", mock.Anything, int64(555), 9, 2026).
		Return(reportUC.MonthlyCategoryReport{
			Income:  []reportUC.CategoryReport{{CategoryID: 2}},
			Expense: []reportUC.CategoryReport{{CategoryID: 1}},
		}, nil).Once()

	r := newRouter()
	r.GET("/internal/bot/reports/monthly", h.GetMonthlyReport)

	w := doJSON(r, http.MethodGet, "/internal/bot/reports/monthly?chat_id=555&month=9&year=2026", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_GetMonthlyReport_NotLinked(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetMonthlyReport", mock.Anything, int64(999), 9, 2026).
		Return(reportUC.MonthlyCategoryReport{}, telegram.ErrLinkNotFound).Once()

	r := newRouter()
	r.GET("/internal/bot/reports/monthly", h.GetMonthlyReport)

	w := doJSON(r, http.MethodGet, "/internal/bot/reports/monthly?chat_id=999&month=9&year=2026", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBotHandler_GetMonthlyReport_InvalidMonth(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))

	r := newRouter()
	r.GET("/internal/bot/reports/monthly", h.GetMonthlyReport)

	w := doJSON(r, http.MethodGet, "/internal/bot/reports/monthly?chat_id=555&month=13&year=2026", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "GetMonthlyReport", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func (m *MockBotUsecase) GetBudgetStatus(ctx context.Context, chatID int64, month, year int) ([]uc.BudgetStatus, error) {
	args := m.Called(ctx, chatID, month, year)
	return args.Get(0).([]uc.BudgetStatus), args.Error(1)
}

func TestBotHandler_GetBudgetStatus_Success(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetBudgetStatus", mock.Anything, int64(555), 9, 2026).
		Return([]uc.BudgetStatus{{CategoryID: 1, CategoryName: "Makanan", Limit: 500000, Spent: 250000}}, nil).Once()

	r := newRouter()
	r.GET("/internal/bot/budgets", h.GetBudgetStatus)

	w := doJSON(r, http.MethodGet, "/internal/bot/budgets?chat_id=555&month=9&year=2026", "")
	assert.Equal(t, http.StatusOK, w.Code)
	m.AssertExpectations(t)
}

func TestBotHandler_GetBudgetStatus_NotLinked(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))
	m.On("GetBudgetStatus", mock.Anything, int64(999), 9, 2026).
		Return([]uc.BudgetStatus{}, telegram.ErrLinkNotFound).Once()

	r := newRouter()
	r.GET("/internal/bot/budgets", h.GetBudgetStatus)

	w := doJSON(r, http.MethodGet, "/internal/bot/budgets?chat_id=999&month=9&year=2026", "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBotHandler_GetBudgetStatus_InvalidMonth(t *testing.T) {
	m := new(MockBotUsecase)
	h := handler.NewBotHandler(m, new(MockAuditRepo))

	r := newRouter()
	r.GET("/internal/bot/budgets", h.GetBudgetStatus)

	w := doJSON(r, http.MethodGet, "/internal/bot/budgets?chat_id=555&month=0", "")
	assert.Equal(t, http.StatusBadRequest, w.Code)
	m.AssertNotCalled(t, "GetBudgetStatus", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
