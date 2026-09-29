package telegram

import (
	"context"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/afandimsr/cashbook-bot/internal/llm"
	"github.com/afandimsr/cashbook-bot/internal/ratelimit"
	"github.com/afandimsr/cashbook-bot/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func TestDispatch_LinkCommand_ThrottledAfterMaxAttempts(t *testing.T) {
	be := new(mockBackend)
	ai := new(mockLLM)
	svc := service.New(be, ai)
	b := &Bot{svc: svc, linkLimiter: ratelimit.New(linkAttemptsMax, linkAttemptsWindow)}

	be.On("LinkAccount", mock.Anything, int64(555), "afandi", "BAD").
		Return(&backend.APIError{StatusCode: 400}).Times(linkAttemptsMax)

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 555},
		From:     &tgbotapi.User{UserName: "afandi"},
		Text:     "/link BAD",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}},
	}

	for i := 0; i < linkAttemptsMax; i++ {
		reply, err := b.dispatch(context.Background(), 555, msg)
		assert.NoError(t, err)
		assert.NotContains(t, reply, "Terlalu banyak")
	}

	reply, err := b.dispatch(context.Background(), 555, msg)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Terlalu banyak percobaan")

	be.AssertExpectations(t)
	be.AssertNumberOfCalls(t, "LinkAccount", linkAttemptsMax)
}

func TestDispatch_StartCommand_ReportsLinkedAccount(t *testing.T) {
	be := new(mockBackend)
	ai := new(mockLLM)
	svc := service.New(be, ai)
	b := &Bot{svc: svc, linkLimiter: ratelimit.New(linkAttemptsMax, linkAttemptsWindow)}

	be.On("GetContext", mock.Anything, int64(555)).
		Return(&backend.Context{UserID: 7, Name: "Afandi", Email: "afandi@example.com"}, nil).Once()

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 555},
		From:     &tgbotapi.User{FirstName: "Afandi"},
		Text:     "/start",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 6}},
	}

	reply, err := b.dispatch(context.Background(), 555, msg)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Afandi")
	assert.Contains(t, reply, "afandi@example.com")
	be.AssertExpectations(t)
}

func TestDispatch_StartCommand_PromptsLinkWhenNotLinked(t *testing.T) {
	be := new(mockBackend)
	ai := new(mockLLM)
	svc := service.New(be, ai)
	b := &Bot{svc: svc, linkLimiter: ratelimit.New(linkAttemptsMax, linkAttemptsWindow)}

	be.On("GetContext", mock.Anything, int64(555)).
		Return(nil, &backend.APIError{StatusCode: 404}).Once()

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 555},
		From:     &tgbotapi.User{FirstName: "Afandi"},
		Text:     "/start",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 6}},
	}

	reply, err := b.dispatch(context.Background(), 555, msg)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Afandi")
	assert.Contains(t, reply, "belum terhubung")
}

func TestDispatch_LinkCommand_OtherChatUnaffected(t *testing.T) {
	be := new(mockBackend)
	ai := new(mockLLM)
	svc := service.New(be, ai)
	b := &Bot{svc: svc, linkLimiter: ratelimit.New(1, time.Minute)}

	be.On("LinkAccount", mock.Anything, int64(1), "u1", "CODE").
		Return(&backend.APIError{StatusCode: 400}).Once()
	be.On("LinkAccount", mock.Anything, int64(2), "u2", "CODE").
		Return(&backend.APIError{StatusCode: 400}).Once()

	msg1 := &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 1}, From: &tgbotapi.User{UserName: "u1"}, Text: "/link CODE",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}}}
	msg2 := &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: 2}, From: &tgbotapi.User{UserName: "u2"}, Text: "/link CODE",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}}}

	_, err := b.dispatch(context.Background(), 1, msg1)
	assert.NoError(t, err)
	reply2, err := b.dispatch(context.Background(), 2, msg2)
	assert.NoError(t, err)
	assert.NotContains(t, reply2, "Terlalu banyak")

	be.AssertExpectations(t)
}

func TestDispatch_HelpCommand_AnsweredLocally(t *testing.T) {
	be := new(mockBackend)
	ai := new(mockLLM)
	svc := service.New(be, ai)
	b := &Bot{svc: svc, linkLimiter: ratelimit.New(linkAttemptsMax, linkAttemptsWindow)}

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 555},
		From:     &tgbotapi.User{FirstName: "Afandi"},
		Text:     "/help",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}},
	}

	reply, err := b.dispatch(context.Background(), 555, msg)
	assert.NoError(t, err)
	assert.Equal(t, service.HelpMessage, reply)
	be.AssertNotCalled(t, "GetContext", mock.Anything, mock.Anything)
}

func TestDispatch_KategoriCommand_ListsCategories(t *testing.T) {
	be := new(mockBackend)
	svc := service.New(be, new(mockLLM))
	b := &Bot{svc: svc, linkLimiter: ratelimit.New(linkAttemptsMax, linkAttemptsWindow)}

	be.On("GetContext", mock.Anything, int64(555)).Return(&backend.Context{Categories: []backend.Category{
		{Name: "Makanan", Type: "expense"},
	}}, nil).Once()

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 555},
		From:     &tgbotapi.User{FirstName: "Afandi"},
		Text:     "/kategori",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 9}},
	}

	reply, err := b.dispatch(context.Background(), 555, msg)
	assert.NoError(t, err)
	assert.Contains(t, reply, "- Makanan")
	be.AssertExpectations(t)
}

func TestDispatch_BudgetCommand_UsesCurrentMonth(t *testing.T) {
	be := new(mockBackend)
	svc := service.New(be, new(mockLLM))
	b := &Bot{svc: svc, linkLimiter: ratelimit.New(linkAttemptsMax, linkAttemptsWindow)}

	now := time.Now()
	be.On("GetBudgetStatus", mock.Anything, int64(555), int(now.Month()), now.Year()).
		Return([]backend.BudgetStatus{{CategoryName: "Makanan", Limit: 500000, Spent: 100000, Remaining: 400000, Percentage: 20}}, nil).Once()

	msg := &tgbotapi.Message{
		Chat:     &tgbotapi.Chat{ID: 555},
		From:     &tgbotapi.User{FirstName: "Afandi"},
		Text:     "/budget",
		Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 7}},
	}

	reply, err := b.dispatch(context.Background(), 555, msg)
	assert.NoError(t, err)
	assert.Contains(t, reply, "Makanan")
	be.AssertExpectations(t)
}
