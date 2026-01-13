package user

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

// type LoginRequest struct {
// 	Email    string `json:"email" binding:"required,email" example:"user@mail.com" required:"true"`
// 	Password string `json:"password" binding:"required,min=6" example:"secret" required:"true"`
// }
