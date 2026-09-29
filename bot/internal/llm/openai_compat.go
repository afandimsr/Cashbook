package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/afandimsr/cashbook-bot/internal/backend"
)

// maxErrorBodyBytes bounds how much of a provider's error body ends up in an
// error message (and therefore in logs).
const maxErrorBodyBytes = 300

// OpenAICompatClient talks to any OpenAI-compatible /chat/completions
// endpoint with function calling — Groq today, and e.g. OpenRouter later
// with only a new config entry.
type OpenAICompatClient struct {
	baseURL    string
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewOpenAICompat takes the API base URL without the /chat/completions
// suffix, e.g. https://api.groq.com/openai/v1. Request deadlines come from the
// caller's context (see Chain); timeout is only a safety net.
func NewOpenAICompat(baseURL, apiKey, model string, timeout time.Duration) *OpenAICompatClient {
	return &OpenAICompatClient{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		model:      model,
		httpClient: &http.Client{Timeout: timeout},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatTool struct {
	Type     string       `json:"type"`
	Function chatFunction `json:"function"`
}

type chatFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type chatRequest struct {
	Model      string        `json:"model"`
	Messages   []chatMessage `json:"messages"`
	Tools      []chatTool    `json:"tools"`
	ToolChoice string        `json:"tool_choice"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			ToolCalls []struct {
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"` // JSON-encoded object
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
	} `json:"choices"`
}

func (c *OpenAICompatClient) ParseTransaction(ctx context.Context, message string, categories []backend.Category) (*ParsedTransaction, error) {
	reqBody, err := json.Marshal(chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemInstruction},
			{Role: "user", Content: message},
		},
		Tools: []chatTool{{
			Type: "function",
			Function: chatFunction{
				Name:        toolName,
				Description: toolDescription,
				Parameters:  transactionSchema(categoryNames(categories)),
			},
		}},
		// "auto" rather than forcing the tool: not every OpenAI-compatible
		// model accepts "required", and a text answer maps cleanly to
		// ErrNoFunctionCall anyway.
		ToolChoice: "auto",
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("chat completions request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
		return nil, fmt.Errorf("chat completions: status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode chat completions response: %w", err)
	}

	for _, choice := range out.Choices {
		for _, call := range choice.Message.ToolCalls {
			if call.Function.Name != toolName {
				continue
			}
			var args map[string]any
			if err := json.Unmarshal([]byte(call.Function.Arguments), &args); err != nil {
				return nil, fmt.Errorf("%w: arguments are not valid JSON: %v", ErrInvalidOutput, err)
			}
			return parseArgs(args)
		}
	}
	return nil, ErrNoFunctionCall
}
