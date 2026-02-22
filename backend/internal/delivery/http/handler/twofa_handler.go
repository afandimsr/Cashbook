package handler

import (
	"net/http"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type TwoFAHandler struct {
	usecase *uc.TwoFAUsecase
}

func NewTwoFAHandler(usecase *uc.TwoFAUsecase) *TwoFAHandler {
	return &TwoFAHandler{usecase: usecase}
}

// Setup2FA godoc
// @Summary      Initialize 2FA setup
// @Description  Generate TOTP secret and QR code for the authenticated user.
// @Tags         2FA
// @Produce      json
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Router       /2fa/setup [post]
func (h *TwoFAHandler) Setup(c *gin.Context) {
	userID, _ := c.Get("user_id")

	result, err := h.usecase.Setup(userID.(int64))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "2FA setup initiated", result)
}

// VerifySetup godoc
// @Summary      Confirm 2FA setup
// @Description  Validate the initial TOTP code to enable 2FA.
// @Tags         2FA
// @Accept       json
// @Produce      json
// @Param        body body user.TwoFASetupVerifyRequest true "Verify payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Router       /2fa/setup/verify [post]
func (h *TwoFAHandler) VerifySetup(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req user.TwoFASetupVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "422", "invalid request", err.Error())
		return
	}

	if err := h.usecase.VerifySetup(userID.(int64), req.Code); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "2FA enabled successfully", nil)
}

// VerifyLogin godoc
// @Summary      Verify 2FA during login
// @Description  Validate TOTP code with temporary token to complete login.
// @Tags         2FA
// @Accept       json
// @Produce      json
// @Param        body body user.TwoFAVerifyRequest true "Verify payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Router       /2fa/verify [post]
func (h *TwoFAHandler) VerifyLogin(c *gin.Context) {
	var req user.TwoFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "422", "invalid request", err.Error())
		return
	}

	token, err := h.usecase.VerifyLogin(req.TempToken, req.Code)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "login success", user.LoginResponse{
		Token: token,
	})
}

// Disable2FA godoc
// @Summary      Disable 2FA
// @Description  Turn off two-factor authentication for the current user.
// @Tags         2FA
// @Produce      json
// @Success      200 {object} response.SuccessSingleUserResponse
// @Router       /2fa/disable [delete]
func (h *TwoFAHandler) Disable(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := h.usecase.Disable(userID.(int64)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "2FA disabled", nil)
}

// GenerateBackupCodes godoc
// @Summary      Generate backup codes
// @Description  Generate 10 one-time backup codes for account recovery.
// @Tags         2FA
// @Produce      json
// @Success      200 {object} response.SuccessSingleUserResponse
// @Router       /2fa/backup-codes [post]
func (h *TwoFAHandler) GenerateBackupCodes(c *gin.Context) {
	userID, _ := c.Get("user_id")

	codes, err := h.usecase.GenerateBackupCodes(userID.(int64))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "backup codes generated", codes)
}

// VerifyBackupCode godoc
// @Summary      Verify backup code during login
// @Description  Use a one-time backup code to complete 2FA login.
// @Tags         2FA
// @Accept       json
// @Produce      json
// @Param        body body user.TwoFAVerifyRequest true "Backup verify payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Router       /2fa/backup/verify [post]
func (h *TwoFAHandler) VerifyBackupCode(c *gin.Context) {
	var req user.TwoFAVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "422", "invalid request", err.Error())
		return
	}

	token, err := h.usecase.VerifyBackupCode(req.TempToken, req.Code)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "login success", user.LoginResponse{
		Token: token,
	})
}
