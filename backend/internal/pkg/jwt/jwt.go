package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

func SetSecret(secret string) {
	secretKey = []byte(secret)
}

type Claims struct {
	UserID int64    `json:"user_id"`
	Email  string   `json:"email"`
	Name   string   `json:"name,omitempty"`
	Roles  []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

type TempClaims struct {
	UserID  int64  `json:"user_id"`
	Email   string `json:"email"`
	Purpose string `json:"purpose"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int64, email string, name string, roles []string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenerateTempToken creates a short-lived token for 2FA verification.
func GenerateTempToken(userID int64, email string, purpose string) (string, error) {
	claims := &TempClaims{
		UserID:  userID,
		Email:   email,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateTempToken validates a temp token and checks its purpose.
func ValidateTempToken(tokenString string, expectedPurpose string) (*TempClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TempClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TempClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid temp token")
	}

	if claims.Purpose != expectedPurpose {
		return nil, errors.New("token purpose mismatch")
	}

	return claims, nil
}
