package bot_test

import (
	"context"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/budget"
	"github.com/afandimsr/cashbook-backend/internal/domain/category"
	"github.com/afandimsr/cashbook-backend/internal/domain/telegram"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/bot"
	reportUC "github.com/afandimsr/cashbook-backend/internal/usecase/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTelegramRepo struct{ mock.Mock }

func (m *MockTelegramRepo) FindLinkByChatID(ctx context.Context, chatID int64) (telegram.UserTelegramLink, error) {
	args := m.Called(ctx, chatID)
	return args.Get(0).(telegram.UserTelegramLink), args.Error(1)
}
func (m *MockTelegramRepo) FindLinkByUserID(ctx context.Context, userID int64) (telegram.UserTelegramLink, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(telegram.UserTelegramLink), args.Error(1)
}
func (m *MockTelegramRepo) ListActiveLinks(ctx context.Context) ([]telegram.UserTelegramLink, error) {
	args := m.Called(ctx)
	return args.Get(0).([]telegram.UserTelegramLink), args.Error(1)
}
func (m *MockTelegramRepo) SaveLink(ctx context.Context, link *telegram.UserTelegramLink) error {
	return m.Called(ctx, link).Error(0)
}
func (m *MockTelegramRepo) Unlink(ctx context.Context, userID int64) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *MockTelegramRepo) CreateLinkCode(ctx context.Context, code *telegram.LinkCode) error {
	return m.Called(ctx, code).Error(0)
}
func (m *MockTelegramRepo) FindActiveLinkCode(ctx context.Context, code string) (telegram.LinkCode, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(telegram.LinkCode), args.Error(1)
}
func (m *MockTelegramRepo) FindActiveLinkCodeByUserID(ctx context.Context, userID int64) (telegram.LinkCode, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(telegram.LinkCode), args.Error(1)
}
func (m *MockTelegramRepo) MarkLinkCodeUsed(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockTelegramRepo) WithTransaction(ctx context.Context, fn func(repo telegram.Repository) error) error {
	m.Called(ctx)
	return fn(m)
}

type MockCategoryUsecase struct{ mock.Mock }

func (m *MockCategoryUsecase) GetAllByUserID(userID int64) ([]category.Category, error) {
	args := m.Called(userID)
	return args.Get(0).([]category.Category), args.Error(1)
}
func (m *MockCategoryUsecase) GetByID(id int64) (category.Category, error) {
	args := m.Called(id)
	return args.Get(0).(category.Category), args.Error(1)
}
func (m *MockCategoryUsecase) Create(c category.Category) error { return m.Called(c).Error(0) }
func (m *MockCategoryUsecase) Update(id int64, c category.Category) error {
	return m.Called(id, c).Error(0)
}
func (m *MockCategoryUsecase) Delete(id int64) error { return m.Called(id).Error(0) }

type MockTransactionUsecase struct{ mock.Mock }

func (m *MockTransactionUsecase) GetAllByUserID(userID int64, page, limit int, filter transaction.Filter) (transaction.PaginatedTransactions, error) {
	args := m.Called(userID, page, limit, filter)
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

type MockReportUsecase struct{ mock.Mock }

func (m *MockReportUsecase) GetCategorySpending(userID int64, month, year int) ([]reportUC.CategoryReport, error) {
	args := m.Called(userID, month, year)
	return args.Get(0).([]reportUC.CategoryReport), args.Error(1)
}
func (m *MockReportUsecase) GetMonthlyCategoryReport(userID int64, month, year int) (reportUC.MonthlyCategoryReport, error) {
	args := m.Called(userID, month, year)
	return args.Get(0).(reportUC.MonthlyCategoryReport), args.Error(1)
}

type MockUserUsecase struct{ mock.Mock }

func (m *MockUserUsecase) GetByID(id int64) (user.User, error) {
	args := m.Called(id)
	return args.Get(0).(user.User), args.Error(1)
}

func newTestUsecase() (uc.Usecase, *MockTelegramRepo, *MockCategoryUsecase, *MockTransactionUsecase, *MockReportUsecase, *MockUserUsecase) {
	tgRepo := new(MockTelegramRepo)
	catUC := new(MockCategoryUsecase)
	txUC := new(MockTransactionUsecase)
	repUC := new(MockReportUsecase)
	userUC := new(MockUserUsecase)
	return uc.New(tgRepo, catUC, txUC, repUC, userUC, new(MockBudgetUsecase)), tgRepo, catUC, txUC, repUC, userUC
}

type MockBudgetUsecase struct{ mock.Mock }

func (m *MockBudgetUsecase) GetBudgets(userID int64, month, year int) ([]budget.Budget, error) {
	args := m.Called(userID, month, year)
	return args.Get(0).([]budget.Budget), args.Error(1)
}

func TestBot_GenerateLinkCode_MintsNewWhenNoneActive(t *testing.T) {
	u, tgRepo, _, _, _, _ := newTestUsecase()

	tgRepo.On("FindActiveLinkCodeByUserID", mock.Anything, int64(1)).
		Return(telegram.LinkCode{}, telegram.ErrLinkCodeInvalid).Once()
	tgRepo.On("CreateLinkCode", mock.Anything, mock.MatchedBy(func(c *telegram.LinkCode) bool {
		return c.UserID == 1 && len(c.Code) == 6
	})).Return(nil).Once()

	code, err := u.GenerateLinkCode(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, code.Code, 6)
	assert.WithinDuration(t, time.Now().Add(10*time.Minute), code.ExpiresAt, 5*time.Second)
	tgRepo.AssertExpectations(t)
}

func TestBot_GenerateLinkCode_ReusesExistingActiveCode(t *testing.T) {
	u, tgRepo, _, _, _, _ := newTestUsecase()

	existing := telegram.LinkCode{ID: 5, UserID: 1, Code: "EXIST1", ExpiresAt: time.Now().Add(5 * time.Minute)}
	tgRepo.On("FindActiveLinkCodeByUserID", mock.Anything, int64(1)).Return(existing, nil).Once()

	code, err := u.GenerateLinkCode(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, existing.Code, code.Code)
	tgRepo.AssertNotCalled(t, "CreateLinkCode", mock.Anything, mock.Anything)
	tgRepo.AssertExpectations(t)
}

func TestBot_LinkAccount_Success(t *testing.T) {
	u, tgRepo, _, _, _, _ := newTestUsecase()

	lc := telegram.LinkCode{ID: 9, UserID: 42, Code: "ABC123"}
	tgRepo.On("WithTransaction", mock.Anything).Return(nil).Once()
	tgRepo.On("FindActiveLinkCode", mock.Anything, "ABC123").Return(lc, nil).Once()
	tgRepo.On("SaveLink", mock.Anything, mock.MatchedBy(func(l *telegram.UserTelegramLink) bool {
		return l.UserID == 42 && l.TelegramChatID == 555
	})).Return(nil).Once()
	tgRepo.On("MarkLinkCodeUsed", mock.Anything, int64(9)).Return(nil).Once()

	link, err := u.LinkAccount(context.Background(), 555, "afandi", "ABC123")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), link.UserID)
	tgRepo.AssertExpectations(t)
}

func TestBot_LinkAccount_InvalidCode(t *testing.T) {
	u, tgRepo, _, _, _, _ := newTestUsecase()

	tgRepo.On("WithTransaction", mock.Anything).Return(nil).Once()
	tgRepo.On("FindActiveLinkCode", mock.Anything, "BAD").Return(telegram.LinkCode{}, telegram.ErrLinkCodeInvalid).Once()

	_, err := u.LinkAccount(context.Background(), 555, "afandi", "BAD")
	assert.ErrorIs(t, err, telegram.ErrLinkCodeInvalid)
	tgRepo.AssertNotCalled(t, "SaveLink", mock.Anything, mock.Anything)
	tgRepo.AssertExpectations(t)
}

func TestBot_ListLinkedChats(t *testing.T) {
	u, tgRepo, _, _, _, _ := newTestUsecase()

	expected := []telegram.UserTelegramLink{{ID: 1, UserID: 7, TelegramChatID: 555}}
	tgRepo.On("ListActiveLinks", mock.Anything).Return(expected, nil).Once()

	got, err := u.ListLinkedChats(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	tgRepo.AssertExpectations(t)
}

func TestBot_GetLinkStatus_Linked(t *testing.T) {
	u, tgRepo, _, _, _, _ := newTestUsecase()

	linkedAt := time.Now()
	tgRepo.On("FindLinkByUserID", mock.Anything, int64(7)).
		Return(telegram.UserTelegramLink{UserID: 7, TelegramUsername: "afandi", LinkedAt: linkedAt}, nil).Once()

	status, err := u.GetLinkStatus(context.Background(), 7)
	assert.NoError(t, err)
	assert.True(t, status.Linked)
	assert.Equal(t, "afandi", status.TelegramUsername)
	tgRepo.AssertExpectations(t)
}

func TestBot_GetLinkStatus_NotLinked_IsNotAnError(t *testing.T) {
	u, tgRepo, _, _, _, _ := newTestUsecase()

	tgRepo.On("FindLinkByUserID", mock.Anything, int64(7)).
		Return(telegram.UserTelegramLink{}, telegram.ErrLinkNotFound).Once()

	status, err := u.GetLinkStatus(context.Background(), 7)
	assert.NoError(t, err)
	assert.False(t, status.Linked)
}

func TestBot_Unlink(t *testing.T) {
	u, tgRepo, _, _, _, _ := newTestUsecase()

	tgRepo.On("Unlink", mock.Anything, int64(7)).Return(nil).Once()

	err := u.Unlink(context.Background(), 7)
	assert.NoError(t, err)
	tgRepo.AssertExpectations(t)
}

func TestBot_GetContext_Success(t *testing.T) {
	u, tgRepo, catUC, _, _, userUC := newTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(555)).
		Return(telegram.UserTelegramLink{UserID: 7}, nil).Once()
	cats := []category.Category{{ID: 1, UserID: 7, Name: "Food"}}
	catUC.On("GetAllByUserID", int64(7)).Return(cats, nil).Once()
	userUC.On("GetByID", int64(7)).Return(user.User{Name: "Afandi", Email: "afandi@example.com"}, nil).Once()

	ctxData, err := u.GetContext(context.Background(), 555)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), ctxData.UserID)
	assert.Equal(t, cats, ctxData.Categories)
	assert.Equal(t, "Afandi", ctxData.Name)
	assert.Equal(t, "afandi@example.com", ctxData.Email)
	tgRepo.AssertExpectations(t)
	catUC.AssertExpectations(t)
	userUC.AssertExpectations(t)
}

func TestBot_GetContext_NotLinked(t *testing.T) {
	u, tgRepo, catUC, _, _, userUC := newTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(999)).
		Return(telegram.UserTelegramLink{}, telegram.ErrLinkNotFound).Once()

	_, err := u.GetContext(context.Background(), 999)
	assert.ErrorIs(t, err, telegram.ErrLinkNotFound)
	catUC.AssertNotCalled(t, "GetAllByUserID", mock.Anything)
	userUC.AssertNotCalled(t, "GetByID", mock.Anything)
}

func TestBot_CreateTransaction_Success(t *testing.T) {
	u, tgRepo, catUC, txUC, _, _ := newTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(555)).
		Return(telegram.UserTelegramLink{UserID: 7}, nil).Once()
	catUC.On("GetByID", int64(3)).Return(category.Category{ID: 3, UserID: 7}, nil).Once()
	txUC.On("Create", mock.MatchedBy(func(tx transaction.Transaction) bool {
		return tx.UserID == 7 && tx.CategoryID == 3 && tx.Amount == 50000
	})).Return(nil).Once()

	err := u.CreateTransaction(context.Background(), uc.CreateTransactionInput{
		ChatID: 555, CategoryID: 3, Amount: 50000, Type: "expense",
	})
	assert.NoError(t, err)
	tgRepo.AssertExpectations(t)
	catUC.AssertExpectations(t)
	txUC.AssertExpectations(t)
}

func TestBot_CreateTransaction_CategoryNotOwned(t *testing.T) {
	u, tgRepo, catUC, txUC, _, _ := newTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(555)).
		Return(telegram.UserTelegramLink{UserID: 7}, nil).Once()
	catUC.On("GetByID", int64(3)).Return(category.Category{ID: 3, UserID: 99}, nil).Once()

	err := u.CreateTransaction(context.Background(), uc.CreateTransactionInput{
		ChatID: 555, CategoryID: 3, Amount: 50000, Type: "expense",
	})
	assert.ErrorIs(t, err, uc.ErrCategoryNotOwned)
	txUC.AssertNotCalled(t, "Create", mock.Anything)
}

func TestBot_GetSummary_Success(t *testing.T) {
	u, tgRepo, _, txUC, _, _ := newTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(555)).
		Return(telegram.UserTelegramLink{UserID: 7}, nil).Once()
	txUC.On("GetDashboardSummary", int64(7)).
		Return(transaction.DashboardSummary{TotalIncome: 100, TotalExpense: 40, Balance: 60}, nil).Once()

	summary, err := u.GetSummary(context.Background(), 555)
	assert.NoError(t, err)
	assert.Equal(t, 60.0, summary.Balance)
	tgRepo.AssertExpectations(t)
	txUC.AssertExpectations(t)
}

func TestBot_GetCategorySpending_Success(t *testing.T) {
	u, tgRepo, _, _, repUC, _ := newTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(555)).
		Return(telegram.UserTelegramLink{UserID: 7}, nil).Once()
	expected := []reportUC.CategoryReport{{CategoryID: 1, TotalAmount: 100}}
	repUC.On("GetCategorySpending", int64(7), 9, 2026).Return(expected, nil).Once()

	got, err := u.GetCategorySpending(context.Background(), 555, 9, 2026)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	tgRepo.AssertExpectations(t)
	repUC.AssertExpectations(t)
}

func TestBot_GetMonthlyReport_Success(t *testing.T) {
	u, tgRepo, _, _, repUC, _ := newTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(555)).
		Return(telegram.UserTelegramLink{UserID: 7}, nil).Once()
	expected := reportUC.MonthlyCategoryReport{
		Income:  []reportUC.CategoryReport{{CategoryID: 2, TotalAmount: 500}},
		Expense: []reportUC.CategoryReport{{CategoryID: 1, TotalAmount: 100}},
	}
	repUC.On("GetMonthlyCategoryReport", int64(7), 9, 2026).Return(expected, nil).Once()

	got, err := u.GetMonthlyReport(context.Background(), 555, 9, 2026)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	tgRepo.AssertExpectations(t)
	repUC.AssertExpectations(t)
}

func TestBot_GetMonthlyReport_NotLinked(t *testing.T) {
	u, tgRepo, _, _, repUC, _ := newTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(999)).
		Return(telegram.UserTelegramLink{}, telegram.ErrLinkNotFound).Once()

	_, err := u.GetMonthlyReport(context.Background(), 999, 9, 2026)
	assert.ErrorIs(t, err, telegram.ErrLinkNotFound)
	repUC.AssertNotCalled(t, "GetMonthlyCategoryReport", mock.Anything, mock.Anything, mock.Anything)
}

func newBudgetTestUsecase() (uc.Usecase, *MockTelegramRepo, *MockCategoryUsecase, *MockReportUsecase, *MockBudgetUsecase) {
	tgRepo := new(MockTelegramRepo)
	catUC := new(MockCategoryUsecase)
	repUC := new(MockReportUsecase)
	budUC := new(MockBudgetUsecase)
	return uc.New(tgRepo, catUC, new(MockTransactionUsecase), repUC, new(MockUserUsecase), budUC), tgRepo, catUC, repUC, budUC
}

func TestBot_GetBudgetStatus_JoinsBudgetsCategoriesAndSpend(t *testing.T) {
	u, tgRepo, catUC, repUC, budUC := newBudgetTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(555)).
		Return(telegram.UserTelegramLink{UserID: 7}, nil).Once()
	budUC.On("GetBudgets", int64(7), 9, 2026).Return([]budget.Budget{
		{CategoryID: 1, Amount: 500000},
		{CategoryID: 2, Amount: 300000},
		{CategoryID: 3, Amount: 100000}, // category deleted: skipped
		{CategoryID: 4, Amount: 0},      // no limit: skipped
	}, nil).Once()
	catUC.On("GetAllByUserID", int64(7)).Return([]category.Category{
		{ID: 1, Name: "Makanan"}, {ID: 2, Name: "Hiburan"}, {ID: 4, Name: "Lainnya"},
	}, nil).Once()
	repUC.On("GetMonthlyCategoryReport", int64(7), 9, 2026).Return(reportUC.MonthlyCategoryReport{
		Expense: []reportUC.CategoryReport{
			{CategoryID: 1, TotalAmount: 250000},
			{CategoryID: 2, TotalAmount: 360000},
		},
	}, nil).Once()

	got, err := u.GetBudgetStatus(context.Background(), 555, 9, 2026)
	assert.NoError(t, err)
	assert.Len(t, got, 2)

	// Sorted by percentage used, highest first.
	assert.Equal(t, "Hiburan", got[0].CategoryName)
	assert.True(t, got[0].IsOver)
	assert.Equal(t, -60000.0, got[0].Remaining)
	assert.Equal(t, 120.0, got[0].Percentage)

	assert.Equal(t, "Makanan", got[1].CategoryName)
	assert.False(t, got[1].IsOver)
	assert.Equal(t, 250000.0, got[1].Remaining)
	assert.Equal(t, 50.0, got[1].Percentage)
}

func TestBot_GetBudgetStatus_NoBudgets_SkipsOtherLookups(t *testing.T) {
	u, tgRepo, catUC, repUC, budUC := newBudgetTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(555)).
		Return(telegram.UserTelegramLink{UserID: 7}, nil).Once()
	budUC.On("GetBudgets", int64(7), 9, 2026).Return([]budget.Budget{}, nil).Once()

	got, err := u.GetBudgetStatus(context.Background(), 555, 9, 2026)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Empty(t, got)
	catUC.AssertNotCalled(t, "GetAllByUserID", mock.Anything)
	repUC.AssertNotCalled(t, "GetMonthlyCategoryReport", mock.Anything, mock.Anything, mock.Anything)
}

func TestBot_GetBudgetStatus_NotLinked(t *testing.T) {
	u, tgRepo, _, _, budUC := newBudgetTestUsecase()

	tgRepo.On("FindLinkByChatID", mock.Anything, int64(999)).
		Return(telegram.UserTelegramLink{}, telegram.ErrLinkNotFound).Once()

	_, err := u.GetBudgetStatus(context.Background(), 999, 9, 2026)
	assert.ErrorIs(t, err, telegram.ErrLinkNotFound)
	budUC.AssertNotCalled(t, "GetBudgets", mock.Anything, mock.Anything, mock.Anything)
}
