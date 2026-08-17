package totp_test

import (
	"testing"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/infrastructure/totp"
	otptotp "github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSecret(t *testing.T) {
	secret, qr, err := totp.GenerateSecret("user@example.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, secret)
	assert.NotEmpty(t, qr) // base64 PNG
}

func TestValidateCode_Valid(t *testing.T) {
	secret, _, err := totp.GenerateSecret("user@example.com")
	assert.NoError(t, err)

	code, err := otptotp.GenerateCode(secret, time.Now())
	assert.NoError(t, err)

	assert.True(t, totp.ValidateCode(secret, code))
}

func TestValidateCode_Invalid(t *testing.T) {
	secret, _, _ := totp.GenerateSecret("user@example.com")
	assert.False(t, totp.ValidateCode(secret, "000000"))
}
