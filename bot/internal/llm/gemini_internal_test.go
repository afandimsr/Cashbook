package llm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs_Success(t *testing.T) {
	args := map[string]any{
		"amount":        50000.0,
		"type":          "expense",
		"category_name": "Makanan",
		"note":          "bensin",
	}

	parsed, err := parseArgs(args)
	assert.NoError(t, err)
	assert.Equal(t, 50000.0, parsed.Amount)
	assert.Equal(t, "expense", parsed.Type)
	assert.Equal(t, "Makanan", parsed.CategoryName)
	assert.Equal(t, "bensin", parsed.Note)
}

func TestParseArgs_MissingAmount(t *testing.T) {
	args := map[string]any{
		"type":          "expense",
		"category_name": "Makanan",
	}

	_, err := parseArgs(args)
	assert.Error(t, err)
}

func TestParseArgs_MissingCategoryName(t *testing.T) {
	args := map[string]any{
		"amount": 50000.0,
		"type":   "expense",
	}

	_, err := parseArgs(args)
	assert.Error(t, err)
}
