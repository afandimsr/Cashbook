package backend_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetContext_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/internal/bot/context", r.URL.Path)
		assert.Equal(t, "555", r.URL.Query().Get("chat_id"))
		assert.Equal(t, "test-key", r.Header.Get("X-Internal-Api-Key"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"message":"ok","data":{"user_id":7,"categories":[{"id":1,"name":"Makanan","type":"expense"}]}}`))
	}))
	defer srv.Close()

	c := backend.New(srv.URL, "test-key", 5*time.Second)
	ctx, err := c.GetContext(context.Background(), 555)
	assert.NoError(t, err)
	assert.Equal(t, int64(7), ctx.UserID)
	assert.Len(t, ctx.Categories, 1)
	assert.Equal(t, "Makanan", ctx.Categories[0].Name)
}

func TestClient_GetContext_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"success":false,"message":"telegram account is not linked to a user","error_code":"NOT_FOUND"}`))
	}))
	defer srv.Close()

	c := backend.New(srv.URL, "test-key", 5*time.Second)
	_, err := c.GetContext(context.Background(), 999)
	assert.Error(t, err)

	apiErr, ok := err.(*backend.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, "NOT_FOUND", apiErr.ErrorCode)
}

func TestClient_CreateTransaction_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/internal/bot/transactions", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success":true,"message":"transaction created"}`))
	}))
	defer srv.Close()

	c := backend.New(srv.URL, "test-key", 5*time.Second)
	err := c.CreateTransaction(context.Background(), backend.CreateTransactionRequest{
		ChatID: 555, CategoryID: 1, Amount: 50000, Type: "expense",
	})
	assert.NoError(t, err)
}

func TestClient_LinkAccount_InvalidCode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"success":false,"message":"link code is invalid, expired, or already used","error_code":"BAD_REQUEST"}`))
	}))
	defer srv.Close()

	c := backend.New(srv.URL, "test-key", 5*time.Second)
	err := c.LinkAccount(context.Background(), 555, "afandi", "BAD")
	apiErr, ok := err.(*backend.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
}

func TestClient_GetMonthlyReport_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/internal/bot/reports/monthly", r.URL.Path)
		assert.Equal(t, "555", r.URL.Query().Get("chat_id"))
		assert.Equal(t, "8", r.URL.Query().Get("month"))
		assert.Equal(t, "2026", r.URL.Query().Get("year"))
		w.Write([]byte(`{"success":true,"data":{"income":[{"category_id":2,"category_name":"Gaji","total_amount":5000000}],"expense":[{"category_id":1,"category_name":"Makanan","total_amount":100000}]}}`))
	}))
	defer srv.Close()

	c := backend.New(srv.URL, "test-key", 5*time.Second)
	report, err := c.GetMonthlyReport(context.Background(), 555, 8, 2026)
	assert.NoError(t, err)
	assert.Len(t, report.Income, 1)
	assert.Equal(t, "Gaji", report.Income[0].CategoryName)
	assert.Len(t, report.Expense, 1)
	assert.Equal(t, 100000.0, report.Expense[0].TotalAmount)
}

func TestClient_GetBudgetStatus_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/internal/bot/budgets", r.URL.Path)
		assert.Equal(t, "555", r.URL.Query().Get("chat_id"))
		assert.Equal(t, "9", r.URL.Query().Get("month"))
		assert.Equal(t, "2026", r.URL.Query().Get("year"))
		w.Write([]byte(`{"success":true,"data":[{"category_id":2,"category_name":"Hiburan","limit":300000,"spent":360000,"remaining":-60000,"percentage":120,"is_over":true}]}`))
	}))
	defer srv.Close()

	c := backend.New(srv.URL, "test-key", 5*time.Second)
	got, err := c.GetBudgetStatus(context.Background(), 555, 9, 2026)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, "Hiburan", got[0].CategoryName)
	assert.True(t, got[0].IsOver)
	assert.Equal(t, -60000.0, got[0].Remaining)
}
