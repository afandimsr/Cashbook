package user

import (
	"errors"
	"fmt"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/auth"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
)

type OAuthUsecase interface {
	GetGoogleAuthURL(state string) string
	HandleGoogleCallback(code string) (string, error)
}

type oauthUsecase struct {
	userRepo   user.UserRepository
	googleAuth auth.GoogleAuth
}

func NewOAuthUsecase(userRepo user.UserRepository, googleAuth auth.GoogleAuth) OAuthUsecase {
	return &oauthUsecase{
		userRepo:   userRepo,
		googleAuth: googleAuth,
	}
}

func (u *oauthUsecase) GetGoogleAuthURL(state string) string {
	return u.googleAuth.GetAuthURL(state)
}

func (u *oauthUsecase) HandleGoogleCallback(code string) (string, error) {
	token, err := u.googleAuth.ExchangeCode(code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %w", err)
	}

	googleUser, err := u.googleAuth.GetUserData(token)
	if err != nil {
		return "", fmt.Errorf("get user data failed: %w", err)
	}

	if googleUser.Email == "" {
		return "", errors.New("google email is empty")
	}

	// Find user by Google ID or Email
	existingUser, err := u.userRepo.FindByGoogleID(googleUser.ID)
	if err != nil {
		// Try finding by email
		existingUser, err = u.userRepo.FindByEmail(googleUser.Email)
		if err != nil {
			// User doesn't exist, create new
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
			existingUser = newUser
		} else {
			// Update existing user with Google ID
			existingUser.GoogleID = googleUser.ID
			existingUser.IsActive = true

			err = u.userRepo.Update(existingUser)
			if err != nil {
				return "", fmt.Errorf("failed to update user with google id: %w", err)
			}
		}
	}

	// Generate JWT
	jwtToken, err := jwt.GenerateToken(existingUser.ID, existingUser.Email, existingUser.Name, existingUser.Roles)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return jwtToken, nil
}
