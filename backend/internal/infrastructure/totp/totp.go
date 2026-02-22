package totp

import (
	"bytes"
	"encoding/base64"
	"image/png"

	"github.com/pquerna/otp/totp"
)

// GenerateSecret creates a new TOTP secret and returns the secret string and a QR code as base64-encoded PNG.
func GenerateSecret(email string) (secret string, qrCodeBase64 string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "CashBook",
		AccountName: email,
	})
	if err != nil {
		return "", "", err
	}

	// Generate QR code image
	img, err := key.Image(200, 200)
	if err != nil {
		return "", "", err
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", "", err
	}

	qrCodeBase64 = base64.StdEncoding.EncodeToString(buf.Bytes())
	return key.Secret(), qrCodeBase64, nil
}

// ValidateCode checks a TOTP code against a secret.
func ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
