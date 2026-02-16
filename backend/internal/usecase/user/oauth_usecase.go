package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/auth"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
)

type OAuthUsecase interface {
	GetGoogleAuthURL(ip, userAgent string) (string, error)
	HandleGoogleCallback(code, state, ip, userAgent string) (string, error)
}

type oauthUsecase struct {
	userRepo       user.UserRepository
	oauthStateRepo user.OauthStateRepository
	googleAuth     auth.GoogleAuth
}

func NewOAuthUsecase(userRepo user.UserRepository, oauthStateRepo user.OauthStateRepository, googleAuth auth.GoogleAuth) OAuthUsecase {
	return &oauthUsecase{
		userRepo:       userRepo,
		oauthStateRepo: oauthStateRepo,
		googleAuth:     googleAuth,
	}
}

func (u *oauthUsecase) GetGoogleAuthURL(ip, userAgent string) (string, error) {
	// Generate random state
	state, err := generateRandomString(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	// Generate ID
	uuid, err := generateUUID()
	if err != nil {
		return "", fmt.Errorf("failed to generate uuid: %w", err)
	}

	// Hash IP and User Agent
	ipHash := hashString(ip)
	uaHash := hashString(userAgent)

	oauthState := user.OauthState{
		ID:            uuid,
		State:         state,
		Provider:      "google",
		IPHash:        ipHash,
		UserAgentHash: uaHash,
		CreatedAt:     time.Now(),
		ExpiresAt:     time.Now().Add(10 * time.Minute), // 10 minutes expiration
	}

	if err := u.oauthStateRepo.Save(oauthState); err != nil {
		return "", fmt.Errorf("failed to save oauth state: %w", err)
	}

	return u.googleAuth.GetAuthURL(state), nil
}

func (u *oauthUsecase) HandleGoogleCallback(code, state, ip, userAgent string) (string, error) {
	// 1. Verify State
	storedState, err := u.oauthStateRepo.FindByState(state)
	if err != nil {
		return "", errors.New("invalid oauth state")
	}

	// 2. Check Expiration
	if time.Now().After(storedState.ExpiresAt) {
		return "", errors.New("oauth state expired")
	}

	// 3. Check Usage (Replay Attack Protection)
	if storedState.UsedAt != nil {
		return "", errors.New("oauth state already used")
	}

	// 4. Verify IP and User Agent Binding
	// Note: IP check can be flaky if user switches networks (WiFi -> 4G).
	// For strict security, we enforce it. For better UX, might consider relaxing it or logging warning.
	// Given the user provided the design with IP binding, I will enforce it.
	// 4. Verify IP and User Agent Binding
	currentIPHash := hashString(ip)
	currentUAHash := hashString(userAgent)

	if storedState.IPHash != "" && storedState.IPHash != currentIPHash {
		return "", errors.New("ip address mismatch")
	}
	if storedState.UserAgentHash != "" && storedState.UserAgentHash != currentUAHash {
		return "", errors.New("user agent mismatch")
	}

	// 5. Mark State as Used
	now := time.Now()
	storedState.UsedAt = &now
	if err := u.oauthStateRepo.Update(*storedState); err != nil {
		// Log error but proceed? Or fail? Better fail to be safe against concurrency issues?
		// If update fails, it might mean another request used it.
		return "", fmt.Errorf("failed to mark state as used: %w", err)
	}

	// 6. Exchange Code
	token, err := u.googleAuth.ExchangeCode(code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %w", err)
	}

	// 7. Get User Info
	googleUser, err := u.googleAuth.GetUserData(token)
	if err != nil {
		return "", fmt.Errorf("get user data failed: %w", err)
	}

	if googleUser.Email == "" {
		return "", errors.New("google email is empty")
	}

	// 8. Find or Create User (Logic remains largely same)
	existingUser, err := u.userRepo.FindByGoogleID(googleUser.ID)
	if err != nil {
		existingUser, err = u.userRepo.FindByEmail(googleUser.Email)
		if err != nil {
			newUser := user.User{
				Name:     googleUser.Name,
				Email:    googleUser.Email,
				GoogleID: googleUser.ID,
				IsActive: true,
			}
			err = u.userRepo.Save(newUser)
			if err != nil {
				return "", fmt.Errorf("failed to save new user: %w", err)
			}
			// Fetch again to get ID
			existingUser, _ = u.userRepo.FindByGoogleID(googleUser.ID)
		} else {
			existingUser.GoogleID = googleUser.ID
			existingUser.IsActive = true
			if err := u.userRepo.Update(existingUser); err != nil {
				return "", fmt.Errorf("failed to update user with google id: %w", err)
			}
		}
	}

	jwtToken, err := jwt.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Name, existingUser.Roles)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return jwtToken, nil
}

func hashString(s string) string {
	h := sha256.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func generateUUID() (string, error) {
	idBytes := make([]byte, 16)
	if _, err := rand.Read(idBytes); err != nil {
		return "", err
	}
	idBytes[6] = (idBytes[6] & 0x0f) | 0x40 // Version 4
	idBytes[8] = (idBytes[8] & 0x3f) | 0x80 // Variant 10
	return fmt.Sprintf("%x-%x-%x-%x-%x", idBytes[0:4], idBytes[4:6], idBytes[6:8], idBytes[8:10], idBytes[10:]), nil
}
