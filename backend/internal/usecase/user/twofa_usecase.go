package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/totp"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type TwoFAUsecase struct {
	userRepo       user.UserRepository
	backupCodeRepo user.MFABackupCodeRepository
}

func NewTwoFAUsecase(userRepo user.UserRepository, backupCodeRepo user.MFABackupCodeRepository) *TwoFAUsecase {
	return &TwoFAUsecase{
		userRepo:       userRepo,
		backupCodeRepo: backupCodeRepo,
	}
}

// Setup generates a new TOTP secret and QR code for the user.
// It stores the secret but does NOT enable 2FA until VerifySetup is called.
func (u *TwoFAUsecase) Setup(userID int64) (*user.TwoFASetupResponse, error) {
	existingUser, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if existingUser.TOTPEnabled {
		return nil, apperror.BadRequest("2FA is already enabled", nil)
	}

	secret, qrCode, err := totp.GenerateSecret(existingUser.Email)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	// Store secret (but don't enable yet — user must verify first)
	existingUser.TOTPSecret = secret
	if err := u.userRepo.UpdateTOTPSecret(existingUser); err != nil {
		return nil, apperror.Internal(err)
	}

	return &user.TwoFASetupResponse{
		Secret: secret,
		QRCode: qrCode,
	}, nil
}

// VerifySetup confirms the TOTP setup by validating the initial code.
func (u *TwoFAUsecase) VerifySetup(userID int64, code string) error {
	existingUser, err := u.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	if existingUser.TOTPSecret == "" {
		return apperror.BadRequest("2FA setup not initiated", nil)
	}

	if !totp.ValidateCode(existingUser.TOTPSecret, code) {
		return apperror.BadRequest("invalid TOTP code", nil)
	}

	existingUser.TOTPEnabled = true
	if err := u.userRepo.EnableTOTP(existingUser); err != nil {
		return apperror.Internal(err)
	}

	return nil
}

// Disable turns off 2FA for the user and clears the secret.
func (u *TwoFAUsecase) Disable(userID int64) error {
	existingUser, err := u.userRepo.FindByID(userID)
	if err != nil {
		return err
	}

	existingUser.TOTPSecret = ""
	existingUser.TOTPEnabled = false
	if err := u.userRepo.DisableTOTP(existingUser); err != nil {
		return apperror.Internal(err)
	}

	// Also clear backup codes
	_ = u.backupCodeRepo.DeleteByUserID(userID)

	return nil
}

// VerifyLogin validates the TOTP code during login and returns the full JWT.
func (u *TwoFAUsecase) VerifyLogin(tempToken, code string) (string, error) {
	claims, err := jwt.ValidateTempToken(tempToken, "verify")
	if err != nil {
		return "", apperror.Unauthorized("invalid or expired 2FA token", err)
	}

	existingUser, err := u.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", apperror.Unauthorized("user not found", err)
	}

	if !totp.ValidateCode(existingUser.TOTPSecret, code) {
		return "", apperror.Unauthorized("invalid TOTP code", nil)
	}

	token, err := jwt.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Name, existingUser.Roles)
	if err != nil {
		return "", apperror.Internal(err)
	}

	return token, nil
}

// GenerateBackupCodes creates 10 new one-time backup codes.
func (u *TwoFAUsecase) GenerateBackupCodes(userID int64) ([]string, error) {
	existingUser, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if !existingUser.TOTPEnabled {
		return nil, apperror.BadRequest("2FA must be enabled to generate backup codes", nil)
	}

	// Delete old codes
	_ = u.backupCodeRepo.DeleteByUserID(userID)

	var plainCodes []string
	var hashes []string

	for i := 0; i < 10; i++ {
		code, err := generateBackupCode()
		if err != nil {
			return nil, apperror.Internal(err)
		}
		plainCodes = append(plainCodes, code)

		hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		if err != nil {
			return nil, apperror.Internal(err)
		}
		hashes = append(hashes, string(hash))
	}

	if err := u.backupCodeRepo.SaveBatch(userID, hashes); err != nil {
		return nil, apperror.Internal(err)
	}

	return plainCodes, nil
}

// VerifyBackupCode validates a backup code during login and returns the full JWT.
func (u *TwoFAUsecase) VerifyBackupCode(tempToken, code string) (string, error) {
	claims, err := jwt.ValidateTempToken(tempToken, "verify")
	if err != nil {
		return "", apperror.Unauthorized("invalid or expired backup code token", err)
	}

	codes, err := u.backupCodeRepo.FindByUserID(claims.UserID)
	if err != nil {
		return "", apperror.Internal(err)
	}

	for _, bc := range codes {
		if bc.UsedAt != nil {
			continue // Already used
		}
		if err := bcrypt.CompareHashAndPassword([]byte(bc.CodeHash), []byte(code)); err == nil {
			// Match found — mark as used
			if err := u.backupCodeRepo.MarkUsed(bc.ID); err != nil {
				return "", apperror.Internal(err)
			}

			existingUser, err := u.userRepo.FindByID(claims.UserID)
			if err != nil {
				return "", apperror.Internal(err)
			}

			token, err := jwt.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Name, existingUser.Roles)
			if err != nil {
				return "", apperror.Internal(err)
			}

			return token, nil
		}
	}

	return "", apperror.Unauthorized("invalid backup code", nil)
}

func generateBackupCode() (string, error) {
	b := make([]byte, 4) // 8 hex chars
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	code := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s", code[:4], code[4:]), nil
}
