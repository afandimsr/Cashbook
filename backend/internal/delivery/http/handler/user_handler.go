package handler

import (
	"net/http"
	"strconv"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase      *uc.Usecase
	oauthUsecase uc.OAuthUsecase
}

func New(usecase *uc.Usecase, oauthUsecase uc.OAuthUsecase) *UserHandler {
	return &UserHandler{
		usecase:      usecase,
		oauthUsecase: oauthUsecase,
	}
}

// GetUsers godoc
// @Summary      Get all users
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
// @Summary      Get user by ID
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
// @Summary      Create user
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
// @Summary      Login user
// @Tags         Users
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

	token, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "400", "Username/Password Tidak Valid", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "login success", user.LoginResponse{
		Token: token,
	})
}

// UpdateUser godoc
// @Summary      Update user
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
// @Summary      Delete user
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

// GoogleLogin redirects to Google Auth URL
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

// GoogleCallback handles the callback from Google
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
		// Redirect to frontend error page or return JSON error?
		// Usually callback endpoint redirects to frontend.
		// If error, redirect to login with error param.
		frontendErrorURL := "http://localhost:5173/login?error=" + err.Error()
		c.Redirect(http.StatusTemporaryRedirect, frontendErrorURL)
		return
	}

	// Redirect to frontend with token as query param or set in cookie
	frontendURL := "http://localhost:5173/oauth/callback?token=" + token
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}
