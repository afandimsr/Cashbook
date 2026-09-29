package llm_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"github.com/afandimsr/cashbook-bot/internal/llm"
	"github.com/stretchr/testify/assert"
)

var compatCategories = []backend.Category{
	{ID: 1, Name: "Makanan", Type: "expense"},
	{ID: 2, Name: "Gaji", Type: "income"},
}

func TestOpenAICompat_ParsesToolCall(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/chat/completions", r.URL.Path)
		assert.Equal(t, "Bearer test-key", r.Header.Get("Authorization"))

		var body struct {
			Model    string `json:"model"`
			Messages []struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"messages"`
			Tools []struct {
				Function struct {
					Name       string         `json:"name"`
					Parameters map[string]any `json:"parameters"`
				} `json:"function"`
			} `json:"tools"`
		}
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "openai/gpt-oss-120b", body.Model)
		assert.Equal(t, "beli kopi 25rb", body.Messages[len(body.Messages)-1].Content)
		assert.Equal(t, "create_transaction", body.Tools[0].Function.Name)
		enum := body.Tools[0].Function.Parameters["properties"].(map[string]any)["category_name"].(map[string]any)["enum"]
		assert.Equal(t, []any{"Makanan", "Gaji"}, enum)

		w.Write([]byte(`{"choices":[{"message":{"tool_calls":[{"type":"function","function":{"name":"create_transaction","arguments":"{\"amount\":25000,\"type\":\"expense\",\"category_name\":\"Makanan\",\"note\":\"kopi\"}"}}]}}]}`))
	}))
	defer srv.Close()

	c := llm.NewOpenAICompat(srv.URL+"/v1/", "test-key", "openai/gpt-oss-120b", 5*time.Second)
	got, err := c.ParseTransaction(context.Background(), "beli kopi 25rb", compatCategories)
	assert.NoError(t, err)
	assert.Equal(t, 25000.0, got.Amount)
	assert.Equal(t, "Makanan", got.CategoryName)
	assert.Equal(t, "kopi", got.Note)
}

func TestOpenAICompat_TextAnswerIsNoFunctionCall(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"choices":[{"message":{"content":"Itu bukan transaksi."}}]}`))
	}))
	defer srv.Close()

	c := llm.NewOpenAICompat(srv.URL, "test-key", "openai/gpt-oss-120b", 5*time.Second)
	_, err := c.ParseTransaction(context.Background(), "jam 5", compatCategories)
	assert.ErrorIs(t, err, llm.ErrNoFunctionCall)
}

func TestOpenAICompat_MalformedArgumentsIsInvalidOutput(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"choices":[{"message":{"tool_calls":[{"function":{"name":"create_transaction","arguments":"{not json"}}]}}]}`))
	}))
	defer srv.Close()

	c := llm.NewOpenAICompat(srv.URL, "test-key", "openai/gpt-oss-120b", 5*time.Second)
	_, err := c.ParseTransaction(context.Background(), "beli 50rb", compatCategories)
	assert.ErrorIs(t, err, llm.ErrInvalidOutput)
}

func TestOpenAICompat_RateLimitedIsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error":{"message":"rate limited"}}`))
	}))
	defer srv.Close()

	c := llm.NewOpenAICompat(srv.URL, "test-key", "openai/gpt-oss-120b", 5*time.Second)
	_, err := c.ParseTransaction(context.Background(), "beli 50rb", compatCategories)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "status 429")
	assert.NotErrorIs(t, err, llm.ErrNoFunctionCall)
}
