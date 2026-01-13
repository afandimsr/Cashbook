package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/afandimsr/cashbook-backend/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
}

type GoogleAuth interface {
	GetAuthURL(state string) string
	ExchangeCode(code string) (*oauth2.Token, error)
	GetUserData(token *oauth2.Token) (*GoogleUser, error)
}

type googleAuth struct {
	config *oauth2.Config
}

func NewGoogleAuth(cfg *config.Config) GoogleAuth {
	return &googleAuth{
		config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			RedirectURL:  cfg.GoogleRedirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func (g *googleAuth) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *googleAuth) ExchangeCode(code string) (*oauth2.Token, error) {
	return g.config.Exchange(context.Background(), code)
}

func (g *googleAuth) GetUserData(token *oauth2.Token) (*GoogleUser, error) {
	client := g.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: %s", resp.Status)
	}

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
