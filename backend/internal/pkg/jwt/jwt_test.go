package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	SetSecret("test-secret")

	token, err := GenerateToken(42, "user@example.com", "Alice", []string{"ADMIN"})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), claims.UserID)
	assert.Equal(t, "user@example.com", claims.Email)
	assert.Equal(t, "Alice", claims.Name)
	assert.Equal(t, []string{"ADMIN"}, claims.Roles)
}

func TestValidateToken_Garbage(t *testing.T) {
	SetSecret("test-secret")
	_, err := ValidateToken("not-a-jwt")
	assert.Error(t, err)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	SetSecret("secret-a")
	token, err := GenerateToken(1, "a@b.com", "A", nil)
	assert.NoError(t, err)

	SetSecret("secret-b")
	_, err = ValidateToken(token)
	assert.Error(t, err)
}

func TestTempToken_PurposeMatch(t *testing.T) {
	SetSecret("test-secret")

	token, err := GenerateTempToken(7, "u@e.com", "2fa")
	assert.NoError(t, err)

	claims, err := ValidateTempToken(token, "2fa")
	assert.NoError(t, err)
	assert.Equal(t, int64(7), claims.UserID)
	assert.Equal(t, "2fa", claims.Purpose)
}

func TestTempToken_PurposeMismatch(t *testing.T) {
	SetSecret("test-secret")

	token, err := GenerateTempToken(7, "u@e.com", "2fa")
	assert.NoError(t, err)

	_, err = ValidateTempToken(token, "reset-password")
	assert.Error(t, err)
}

func TestTempToken_RejectsRegularToken(t *testing.T) {
	SetSecret("test-secret")

	// A regular token has no Purpose claim, so it must fail temp validation.
	token, err := GenerateToken(1, "a@b.com", "A", nil)
	assert.NoError(t, err)

	_, err = ValidateTempToken(token, "2fa")
	assert.Error(t, err)
}
