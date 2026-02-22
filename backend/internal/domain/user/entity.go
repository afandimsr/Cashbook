package user

import "time"

type User struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Password    string   `json:"password,omitempty"`
	GoogleID    string   `json:"google_id,omitempty"`
	Roles       []string `json:"roles"`
	IsActive    bool     `json:"is_active"`
	TOTPSecret  string   `json:"-"`
	TOTPEnabled bool     `json:"totp_enabled"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string `json:"token,omitempty"`
	Requires2FA bool   `json:"requires_2fa,omitempty"`
	TempToken   string `json:"temp_token,omitempty"`
}

type PasswordResetRequest struct {
	Password string `json:"password" binding:"required,min=8"`
}

type TwoFASetupResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

type TwoFAVerifyRequest struct {
	TempToken string `json:"temp_token" binding:"required"`
	Code      string `json:"code" binding:"required,min=6,max=9"`
}

type TwoFASetupVerifyRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

type MFASettings struct {
	ID         int64     `json:"id"`
	Enforce2FA bool      `json:"enforce_2fa"`
	UpdatedBy  int64     `json:"updated_by"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type MFABackupCode struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	CodeHash  string     `json:"-"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
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
