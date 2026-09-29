package llm

import (
	"errors"
	"fmt"

	"github.com/afandimsr/cashbook-bot/internal/backend"
)

// ErrNoFunctionCall means the model answered in plain text instead of calling
// create_transaction — typically because the message wasn't a parseable
// transaction. The caller should ask the user to rephrase rather than guess,
// and Chain does not retry it on another provider.
var ErrNoFunctionCall = errors.New("model did not return a create_transaction call")

// ErrInvalidOutput means the model called create_transaction but with missing
// or malformed arguments. Chain retries it on the next provider.
var ErrInvalidOutput = errors.New("invalid create_transaction arguments in model output")

const systemInstruction = `You are a transaction parser for CashBook, an Indonesian personal finance app.
Given a short natural-language message describing money coming in or going out, call create_transaction
with the amount (plain number, no currency symbol), whether it's income or expense, the closest matching
category_name from the provided list (always pick the closest one, even if imperfect — you can only choose
from that list), an optional short note, and an optional date (YYYY-MM-DD, omit for today). Amounts like
"50rb" or "50k" mean 50000; "1jt" or "1 juta" means 1000000. Only skip calling the function if the message
clearly does not describe a financial transaction at all.`

const (
	toolName        = "create_transaction"
	toolDescription = "Record a single income or expense transaction parsed from the user's message."

	amountDescription   = "The transaction amount as a plain number, e.g. 50000."
	categoryDescription = "Must be exactly one of the user's existing category names."
	noteDescription     = "A short note describing the transaction, e.g. the merchant or item."
	dateDescription     = "Date in YYYY-MM-DD format. Omit if the message doesn't specify one."
)

type ParsedTransaction struct {
	Amount       float64
	Type         string // "income" or "expense"
	CategoryName string
	Note         string
	Date         string // YYYY-MM-DD, empty = today
}

// transactionSchema is the create_transaction parameters as plain JSON
// Schema, for providers that take a raw schema (see openai_compat.go).
// Gemini builds the equivalent genai.Schema from the same descriptions.
func transactionSchema(categoryNames []string) map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"amount":        map[string]any{"type": "number", "description": amountDescription},
			"type":          map[string]any{"type": "string", "enum": []string{"income", "expense"}},
			"category_name": map[string]any{"type": "string", "enum": categoryNames, "description": categoryDescription},
			"note":          map[string]any{"type": "string", "description": noteDescription},
			"date":          map[string]any{"type": "string", "description": dateDescription},
		},
		"required": []string{"amount", "type", "category_name"},
	}
}

func categoryNames(categories []backend.Category) []string {
	names := make([]string, len(categories))
	for i, c := range categories {
		names[i] = c.Name
	}
	return names
}

func parseArgs(args map[string]any) (*ParsedTransaction, error) {
	amount, ok := args["amount"].(float64)
	if !ok {
		return nil, fmt.Errorf("%w: missing or invalid amount", ErrInvalidOutput)
	}
	txType, _ := args["type"].(string)
	categoryName, _ := args["category_name"].(string)
	if txType == "" || categoryName == "" {
		return nil, fmt.Errorf("%w: missing type or category_name", ErrInvalidOutput)
	}
	note, _ := args["note"].(string)
	date, _ := args["date"].(string)

	return &ParsedTransaction{
		Amount:       amount,
		Type:         txType,
		CategoryName: categoryName,
		Note:         note,
		Date:         date,
	}, nil
}
