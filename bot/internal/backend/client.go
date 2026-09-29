// Package backend is a thin HTTP client for the CashBook backend's
// /api/v1/internal/bot/* contract (see backend/internal/delivery/http/handler/bot_handler.go).
// It never receives or forwards an end-user JWT: every call is authenticated
// with the shared X-Internal-Api-Key and identifies the acting user only by
// Telegram chat ID, letting the backend resolve identity server-side.
package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL     string
	internalKey string
	httpClient  *http.Client
}

func New(baseURL, internalKey string, timeout time.Duration) *Client {
	return &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		internalKey: internalKey,
		httpClient:  &http.Client{Timeout: timeout},
	}
}

type Category struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"` // "income" or "expense"
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// Context is the resolved user plus their custom categories, fetched once per
// incoming Telegram message so NLP parsing never has to guess a category.
type Context struct {
	UserID     int64      `json:"user_id"`
	Name       string     `json:"name"`
	Email      string     `json:"email"`
	Categories []Category `json:"categories"`
}

type DashboardSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

type CategoryReport struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
	Color        string  `json:"color"`
}

// MonthlyReport is one month's per-category totals, split by type. The
// backend returns each slice sorted by TotalAmount descending.
type MonthlyReport struct {
	Income  []CategoryReport `json:"income"`
	Expense []CategoryReport `json:"expense"`
}

// BudgetStatus is one category budget with that month's spend. The backend
// returns them sorted by Percentage descending (most used first).
type BudgetStatus struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Limit        float64 `json:"limit"`
	Spent        float64 `json:"spent"`
	Remaining    float64 `json:"remaining"` // negative when over budget
	Percentage   float64 `json:"percentage"`
	IsOver       bool    `json:"is_over"`
}

// TelegramLink is one active chat/user pairing, used by the scheduler to know
// who to send automatic monthly reports to.
type TelegramLink struct {
	UserID           int64  `json:"user_id"`
	TelegramChatID   int64  `json:"telegram_chat_id"`
	TelegramUsername string `json:"telegram_username,omitempty"`
}

type CreateTransactionRequest struct {
	ChatID     int64   `json:"chat_id"`
	CategoryID int64   `json:"category_id"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
	Note       string  `json:"note,omitempty"`
	Date       string  `json:"date,omitempty"` // YYYY-MM-DD, empty = today
}

// APIError preserves the backend's HTTP status and error code so callers can
// tell "not linked yet" (404) apart from "invalid link code" (400) and react
// with a helpful Telegram reply instead of a generic failure message.
type APIError struct {
	StatusCode int
	Message    string
	ErrorCode  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("backend error (%d %s): %s", e.StatusCode, e.ErrorCode, e.Message)
}

func (c *Client) ListLinks(ctx context.Context) ([]TelegramLink, error) {
	var out []TelegramLink
	if err := c.doJSON(ctx, http.MethodGet, "/api/v1/internal/bot/links", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) GetContext(ctx context.Context, chatID int64) (*Context, error) {
	var out Context
	if err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/internal/bot/context?chat_id=%d", chatID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) LinkAccount(ctx context.Context, chatID int64, username, code string) error {
	req := map[string]any{"chat_id": chatID, "username": username, "code": code}
	return c.doJSON(ctx, http.MethodPost, "/api/v1/internal/bot/link", req, nil)
}

func (c *Client) CreateTransaction(ctx context.Context, req CreateTransactionRequest) error {
	return c.doJSON(ctx, http.MethodPost, "/api/v1/internal/bot/transactions", req, nil)
}

func (c *Client) GetSummary(ctx context.Context, chatID int64) (*DashboardSummary, error) {
	var out DashboardSummary
	if err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/internal/bot/summary?chat_id=%d", chatID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetMonthlyReport(ctx context.Context, chatID int64, month, year int) (*MonthlyReport, error) {
	var out MonthlyReport
	if err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/internal/bot/reports/monthly?chat_id=%d&month=%d&year=%d", chatID, month, year), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetBudgetStatus(ctx context.Context, chatID int64, month, year int) ([]BudgetStatus, error) {
	var out []BudgetStatus
	if err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/api/v1/internal/bot/budgets?chat_id=%d&month=%d&year=%d", chatID, month, year), nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, reqBody, out any) error {
	var bodyReader io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Api-Key", c.internalKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var envelope struct {
		Success   bool            `json:"success"`
		Message   string          `json:"message"`
		ErrorCode string          `json:"error_code"`
		Data      json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("decode backend response: %w", err)
	}

	if resp.StatusCode >= 400 || !envelope.Success {
		return &APIError{StatusCode: resp.StatusCode, Message: envelope.Message, ErrorCode: envelope.ErrorCode}
	}

	if out != nil && len(envelope.Data) > 0 {
		if err := json.Unmarshal(envelope.Data, out); err != nil {
			return fmt.Errorf("decode backend data: %w", err)
		}
	}
	return nil
}
