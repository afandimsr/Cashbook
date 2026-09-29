package bot

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"sort"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/budget"
	"github.com/afandimsr/cashbook-backend/internal/domain/category"
	"github.com/afandimsr/cashbook-backend/internal/domain/telegram"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	categoryUC "github.com/afandimsr/cashbook-backend/internal/usecase/category"
	reportUC "github.com/afandimsr/cashbook-backend/internal/usecase/report"
	transactionUC "github.com/afandimsr/cashbook-backend/internal/usecase/transaction"
)

// ErrCategoryNotOwned is returned when a bot-resolved user tries to post a
// transaction against a category_id that belongs to a different user. The LLM
// output is never trusted directly, but this guards against a stale/forged
// category_id reaching this far too.
var ErrCategoryNotOwned = errors.New("category does not belong to the linked user")

const (
	linkCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // no 0/O/1/I to avoid typos when copied into Telegram
	linkCodeLength   = 6
	linkCodeTTL      = 10 * time.Minute
)

// Context is what the bot fetches once per incoming message: the resolved
// user plus their custom categories, so NLP parsing never has to guess.
type Context struct {
	UserID     int64               `json:"user_id"`
	Name       string              `json:"name"`
	Email      string              `json:"email"`
	Categories []category.Category `json:"categories"`
}

// userLookup is the minimal user-record accessor this usecase needs,
// satisfied by *user.Usecase; declared here so it stays mockable in tests
// without depending on the whole user usecase surface.
type userLookup interface {
	GetByID(id int64) (user.User, error)
}

// budgetLookup is the minimal budget accessor this usecase needs, satisfied
// by budget.Usecase; declared here for the same reason as userLookup.
type budgetLookup interface {
	GetBudgets(userID int64, month, year int) ([]budget.Budget, error)
}

// BudgetStatus is one category budget joined with that month's actual
// expense total, so the bot doesn't have to combine three calls itself.
type BudgetStatus struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Limit        float64 `json:"limit"`
	Spent        float64 `json:"spent"`
	Remaining    float64 `json:"remaining"` // negative when over budget
	Percentage   float64 `json:"percentage"`
	IsOver       bool    `json:"is_over"`
}

// LinkStatus reports whether a user currently has a Telegram account linked.
// GetLinkStatus treats "not linked" as a normal state, not an error.
type LinkStatus struct {
	Linked           bool       `json:"linked"`
	TelegramUsername string     `json:"telegram_username,omitempty"`
	LinkedAt         *time.Time `json:"linked_at,omitempty"`
}

type CreateTransactionInput struct {
	ChatID     int64
	CategoryID int64
	Amount     float64
	Type       string
	Note       string
	Date       time.Time
}

// Usecase is the internal-only facade the bot service talks to. Every method
// takes a Telegram identifier (chat_id / the caller), never a bot-supplied
// user_id — identity is always resolved server-side from the linked mapping.
type Usecase interface {
	GenerateLinkCode(ctx context.Context, userID int64) (telegram.LinkCode, error)
	LinkAccount(ctx context.Context, chatID int64, username, code string) (telegram.UserTelegramLink, error)
	ListLinkedChats(ctx context.Context) ([]telegram.UserTelegramLink, error)
	GetLinkStatus(ctx context.Context, userID int64) (LinkStatus, error)
	Unlink(ctx context.Context, userID int64) error
	GetContext(ctx context.Context, chatID int64) (Context, error)
	CreateTransaction(ctx context.Context, in CreateTransactionInput) error
	GetSummary(ctx context.Context, chatID int64) (transaction.DashboardSummary, error)
	GetCategorySpending(ctx context.Context, chatID int64, month, year int) ([]reportUC.CategoryReport, error)
	GetMonthlyReport(ctx context.Context, chatID int64, month, year int) (reportUC.MonthlyCategoryReport, error)
	GetBudgetStatus(ctx context.Context, chatID int64, month, year int) ([]BudgetStatus, error)
}

type usecase struct {
	telegramRepo  telegram.Repository
	categoryUC    categoryUC.Usecase
	transactionUC transactionUC.Usecase
	reportUC      reportUC.Usecase
	userUC        userLookup
	budgetUC      budgetLookup
}

func New(telegramRepo telegram.Repository, categoryUC categoryUC.Usecase, transactionUC transactionUC.Usecase, reportUC reportUC.Usecase, userUC userLookup, budgetUC budgetLookup) Usecase {
	return &usecase{
		telegramRepo:  telegramRepo,
		categoryUC:    categoryUC,
		transactionUC: transactionUC,
		reportUC:      reportUC,
		userUC:        userUC,
		budgetUC:      budgetUC,
	}
}

func (u *usecase) resolveUserID(ctx context.Context, chatID int64) (int64, error) {
	link, err := u.telegramRepo.FindLinkByChatID(ctx, chatID)
	if err != nil {
		return 0, err
	}
	return link.UserID, nil
}

// GenerateLinkCode reuses a still-valid outstanding code for this user instead
// of always minting a new one, so repeated clicks don't multiply the number of
// valid guessable codes.
func (u *usecase) GenerateLinkCode(ctx context.Context, userID int64) (telegram.LinkCode, error) {
	if existing, err := u.telegramRepo.FindActiveLinkCodeByUserID(ctx, userID); err == nil {
		return existing, nil
	}

	code := telegram.LinkCode{
		UserID:    userID,
		Code:      generateCode(),
		ExpiresAt: time.Now().Add(linkCodeTTL),
	}
	if err := u.telegramRepo.CreateLinkCode(ctx, &code); err != nil {
		return telegram.LinkCode{}, err
	}
	return code, nil
}

func (u *usecase) LinkAccount(ctx context.Context, chatID int64, username, code string) (telegram.UserTelegramLink, error) {
	var link telegram.UserTelegramLink
	err := u.telegramRepo.WithTransaction(ctx, func(repo telegram.Repository) error {
		lc, err := repo.FindActiveLinkCode(ctx, code)
		if err != nil {
			return err
		}

		link = telegram.UserTelegramLink{
			UserID:           lc.UserID,
			TelegramChatID:   chatID,
			TelegramUsername: username,
		}
		if err := repo.SaveLink(ctx, &link); err != nil {
			return err
		}

		return repo.MarkLinkCodeUsed(ctx, lc.ID)
	})
	if err != nil {
		return telegram.UserTelegramLink{}, err
	}
	return link, nil
}

// ListLinkedChats returns every active telegram_chat_id the bot may message,
// used by the scheduler to fan out monthly reports (see bot/internal/scheduler).
func (u *usecase) ListLinkedChats(ctx context.Context) ([]telegram.UserTelegramLink, error) {
	return u.telegramRepo.ListActiveLinks(ctx)
}

// GetLinkStatus treats "no link" as a normal, non-error result — this is a
// status check for the settings UI, not a lookup that's supposed to fail.
func (u *usecase) GetLinkStatus(ctx context.Context, userID int64) (LinkStatus, error) {
	link, err := u.telegramRepo.FindLinkByUserID(ctx, userID)
	if errors.Is(err, telegram.ErrLinkNotFound) {
		return LinkStatus{Linked: false}, nil
	}
	if err != nil {
		return LinkStatus{}, err
	}
	linkedAt := link.LinkedAt
	return LinkStatus{Linked: true, TelegramUsername: link.TelegramUsername, LinkedAt: &linkedAt}, nil
}

func (u *usecase) Unlink(ctx context.Context, userID int64) error {
	return u.telegramRepo.Unlink(ctx, userID)
}

func (u *usecase) GetContext(ctx context.Context, chatID int64) (Context, error) {
	userID, err := u.resolveUserID(ctx, chatID)
	if err != nil {
		return Context{}, err
	}

	categories, err := u.categoryUC.GetAllByUserID(userID)
	if err != nil {
		return Context{}, err
	}

	usr, err := u.userUC.GetByID(userID)
	if err != nil {
		return Context{}, err
	}

	return Context{UserID: userID, Name: usr.Name, Email: usr.Email, Categories: categories}, nil
}

func (u *usecase) CreateTransaction(ctx context.Context, in CreateTransactionInput) error {
	userID, err := u.resolveUserID(ctx, in.ChatID)
	if err != nil {
		return err
	}

	cat, err := u.categoryUC.GetByID(in.CategoryID)
	if err != nil {
		return err
	}
	if cat.UserID != userID {
		return ErrCategoryNotOwned
	}

	return u.transactionUC.Create(transaction.Transaction{
		UserID:     userID,
		CategoryID: in.CategoryID,
		Amount:     in.Amount,
		Note:       in.Note,
		Date:       in.Date,
		Type:       in.Type,
	})
}

func (u *usecase) GetSummary(ctx context.Context, chatID int64) (transaction.DashboardSummary, error) {
	userID, err := u.resolveUserID(ctx, chatID)
	if err != nil {
		return transaction.DashboardSummary{}, err
	}
	return u.transactionUC.GetDashboardSummary(userID)
}

func (u *usecase) GetCategorySpending(ctx context.Context, chatID int64, month, year int) ([]reportUC.CategoryReport, error) {
	userID, err := u.resolveUserID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return u.reportUC.GetCategorySpending(userID, month, year)
}

func (u *usecase) GetMonthlyReport(ctx context.Context, chatID int64, month, year int) (reportUC.MonthlyCategoryReport, error) {
	userID, err := u.resolveUserID(ctx, chatID)
	if err != nil {
		return reportUC.MonthlyCategoryReport{}, err
	}
	return u.reportUC.GetMonthlyCategoryReport(userID, month, year)
}

// GetBudgetStatus returns every budget set for the month (limit > 0) with
// its spend, most-used first. Budgets whose category no longer exists are
// skipped. It mirrors the web app's budget page (frontend useBudgets), but
// spend comes from the date-filtered monthly report.
func (u *usecase) GetBudgetStatus(ctx context.Context, chatID int64, month, year int) ([]BudgetStatus, error) {
	userID, err := u.resolveUserID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	budgets, err := u.budgetUC.GetBudgets(userID, month, year)
	if err != nil {
		return nil, err
	}
	statuses := []BudgetStatus{}
	if len(budgets) == 0 {
		return statuses, nil
	}

	categories, err := u.categoryUC.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}
	names := make(map[int64]string, len(categories))
	for _, c := range categories {
		names[c.ID] = c.Name
	}

	report, err := u.reportUC.GetMonthlyCategoryReport(userID, month, year)
	if err != nil {
		return nil, err
	}
	spent := make(map[int64]float64, len(report.Expense))
	for _, r := range report.Expense {
		spent[r.CategoryID] = r.TotalAmount
	}

	for _, b := range budgets {
		name, ok := names[b.CategoryID]
		if !ok || b.Amount <= 0 {
			continue
		}
		s := spent[b.CategoryID]
		statuses = append(statuses, BudgetStatus{
			CategoryID:   b.CategoryID,
			CategoryName: name,
			Limit:        b.Amount,
			Spent:        s,
			Remaining:    b.Amount - s,
			Percentage:   s / b.Amount * 100,
			IsOver:       s > b.Amount,
		})
	}
	sort.Slice(statuses, func(i, j int) bool { return statuses[i].Percentage > statuses[j].Percentage })
	return statuses, nil
}

func generateCode() string {
	b := make([]byte, linkCodeLength)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(linkCodeAlphabet))))
		b[i] = linkCodeAlphabet[n.Int64()]
	}
	return string(b)
}
