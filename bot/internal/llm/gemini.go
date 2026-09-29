// Package llm turns a free-text Telegram message into a structured transaction
// using LLM function calling. The tool schema is built per-request from the
// caller's own categories (see ParseTransaction), so the model can only ever
// pick a category_name the user actually owns — it is still re-validated by
// the caller afterwards, since LLM output is never trusted directly.
//
// Each provider (gemini.go, openai_compat.go) implements the same
// ParseTransaction method; Chain (chain.go) tries them in configured order.
package llm

import (
	"context"
	"fmt"
	"log"

	"github.com/afandimsr/cashbook-bot/internal/backend"
	"google.golang.org/genai"
)

type GeminiClient struct {
	genai *genai.Client
	model string
}

func NewGemini(ctx context.Context, apiKey, model string) (*GeminiClient, error) {
	c, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create genai client: %w", err)
	}
	return &GeminiClient{genai: c, model: model}, nil
}

func (c *GeminiClient) ParseTransaction(ctx context.Context, message string, categories []backend.Category) (*ParsedTransaction, error) {
	tool := &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        toolName,
				Description: toolDescription,
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"amount": {
							Type:        genai.TypeNumber,
							Description: amountDescription,
						},
						"type": {
							Type: genai.TypeString,
							Enum: []string{"income", "expense"},
						},
						"category_name": {
							Type:        genai.TypeString,
							Enum:        categoryNames(categories),
							Description: categoryDescription,
						},
						"note": {
							Type:        genai.TypeString,
							Description: noteDescription,
						},
						"date": {
							Type:        genai.TypeString,
							Description: dateDescription,
						},
					},
					Required: []string{"amount", "type", "category_name"},
				},
			},
		},
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(systemInstruction, genai.RoleUser),
		Tools:             []*genai.Tool{tool},
		ToolConfig: &genai.ToolConfig{
			FunctionCallingConfig: &genai.FunctionCallingConfig{
				Mode:                 genai.FunctionCallingConfigModeAny,
				AllowedFunctionNames: []string{toolName},
			},
		},
		// Tambahkan ThinkingConfig untuk meredakan beban latensi & mencegah 504
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingLevel: genai.ThinkingLevelLow,
		},
	}

	resp, err := c.genai.Models.GenerateContent(ctx, c.model, genai.Text(message), config)
	if err != nil {
		return nil, fmt.Errorf("gemini generate content: %w", err)
	}

	calls := resp.FunctionCalls()
	if len(calls) == 0 {
		log.Printf("gemini: no function call for message %q", message)
		return nil, ErrNoFunctionCall
	}

	return parseArgs(calls[0].Args)
}
