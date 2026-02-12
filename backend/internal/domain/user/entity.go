package user

import "time"

type User struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password,omitempty"`
	GoogleID string   `json:"google_id,omitempty"`
	Roles    []string `json:"roles"`
	IsActive bool     `json:"is_active"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type PasswordResetRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type OauthState struct {
	ID            string     `json:"id"`
	State         string     `json:"state"`
	Provider      string     `json:"provider"`
	ClientID      string     `json:"client_id,omitempty"`
	RedirectURI   string     `json:"redirect_uri,omitempty"`
	IPHash        string     `json:"ip_hash,omitempty"`
	UserAgentHash string     `json:"user_agent_hash,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	ExpiresAt     time.Time  `json:"expires_at"`
	UsedAt        *time.Time `json:"used_at,omitempty"`
}
