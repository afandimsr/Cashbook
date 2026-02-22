package handler

import (
	"net/http"
	"strconv"

	"github.com/afandimsr/cashbook-backend/internal/config"
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	cfg          *config.Config
	usecase      *uc.Usecase
	oauthUsecase uc.OAuthUsecase
}

func New(cfg *config.Config, usecase *uc.Usecase, oauthUsecase uc.OAuthUsecase) *UserHandler {
	return &UserHandler{
		cfg:          cfg,
		usecase:      usecase,
		oauthUsecase: oauthUsecase,
	}
}

// GetUsers godoc
// @Summary      Administrative user list
// @Description  Retrieve a paginated list of all registered users in the application.
// @Tags         Users
// @Produce      json
// @Success      200 {object} response.SuccessUserResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, err := h.usecase.GetAll(page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_FETCH_ERROR", "failed to get users", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "success", users)
}

// GetUser godoc
// @Summary      Retrieve user profile
// @Description  Fetch comprehensive details of a specific user account by their unique identifier.
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid user id", err))
		return
	}

	u, err := h.usecase.GetByID(id)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", u)
}

// CreateUser godoc
// @Summary      Provision new user account
// @Description  Create a new user entry with assigned roles and initial status (Admin privileged).
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body body user.User true "User payload"
// @Success      201 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req user.User

	// RBAC: only admin can create new users
	role, _ := c.Get("role")
	if roleStr, ok := role.(string); !ok || roleStr != "ADMIN" {
		response.Error(c, http.StatusForbidden, "FORBIDDEN", "only admin can create users", "")
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "422", "invalid request", err.Error())
		return
	}

	// Ensure roles and is_active have sensible defaults if not provided
	if len(req.Roles) == 0 {
		req.Roles = []string{"USER"}
	}

	// Default is_active true (keeps previous behavior)
	if !req.IsActive {
		req.IsActive = true
	}

	if err := h.usecase.Create(req); err != nil {
		response.Error(c, http.StatusInternalServerError, "USER_CREATE_ERROR", "failed to create user", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "user created", nil)
}

// Login godoc
// @Summary      Authenticate user session
// @Description  Validate user credentials. If 2FA is enabled, returns a temporary token for 2FA verification.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body body user.LoginRequest true "Login payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req user.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "422", "invalid request", err.Error())
		return
	}

	loginResp, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "400", "Username/Password Tidak Valid", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "login success", loginResp)
}

// UpdateUser godoc
// @Summary      Update user information
// @Description  Modify existing user account details and preferences.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        body body user.User true "User payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid user id", err))
		return
	}

	var req user.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	if err := h.usecase.Update(id, req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "user updated", nil)
}

// DeleteUser godoc
// @Summary      Deactivate user account
// @Description  Permanently remove a user identity and access from the system.
// @Tags         Users
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		// c.Error(apperror.BadRequest("invalid user id", err))
		response.Error(c, 400, "INVALID_USER_ID", "Invalid user ID", err.Error())
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		c.Error(err)
		response.Error(c, 500, "ERROR_DELETE_USER", "User gagal dihapus", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "user deleted", nil)
}

// GoogleLogin godoc
// @Summary      Initiate Google OAuth2 flow
// @Description  Redirects the client to Google's OAuth2 authorization page to begin the authentication process.
// @Tags         Auth
// @Produce      json
// @Success      307 "Temporary Redirect to Google Auth URL"
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /auth/google/login [get]
func (h *UserHandler) GoogleLogin(c *gin.Context) {
	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	url, err := h.oauthUsecase.GetGoogleAuthURL(ip, userAgent)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "OAUTH_INIT_ERROR", "failed to initialize oauth", err.Error())
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
// @Summary      Google OAuth2 callback
// @Description  Handles the redirection from Google after user authorization, exchanges code for token, and authenticates user.
// @Tags         Auth
// @Produce      json
// @Param        code   query     string  true  "OAuth2 Code"
// @Param        state  query     string  true  "OAuth2 State"
// @Success      307 "Temporary Redirect to Frontend"
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /auth/google/callback [get]
func (h *UserHandler) GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		response.Error(c, http.StatusBadRequest, "OAUTH_INVALID_REQUEST", "code or state is missing", "")
		return
	}

	ip := c.ClientIP()
	userAgent := c.Request.UserAgent()

	token, err := h.oauthUsecase.HandleGoogleCallback(code, state, ip, userAgent)
	if err != nil {
		c.Error(err)
		// Redirect to login with error param using dynamic ClientAuthURL
		frontendErrorURL := h.cfg.FrontendURL + "/login?error=" + err.Error()
		c.Redirect(http.StatusTemporaryRedirect, frontendErrorURL)
		return
	}

	// Redirect to frontend with token as query param using dynamic ClientAuthURL
	frontendURL := h.cfg.FrontendURL + "/oauth/callback?token=" + token
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

// ResetPassword godoc
// @Summary      Enforce password reset
// @Description  Administrative utility to securely reset a user's password following ISO security standards.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        body body user.PasswordResetRequest true "Reset payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      403 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /users/{id}/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid user id", err))
		return
	}

	var req user.PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "422", "invalid request", err.Error())
		return
	}

	if err := h.usecase.ResetPassword(id, req.Password); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "password reset successfully", nil)
}
