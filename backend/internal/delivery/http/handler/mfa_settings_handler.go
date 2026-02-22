package handler

import (
	"net/http"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/gin-gonic/gin"
)

type MFASettingsHandler struct {
	usecase *uc.MFASettingsUsecase
}

func NewMFASettingsHandler(usecase *uc.MFASettingsUsecase) *MFASettingsHandler {
	return &MFASettingsHandler{usecase: usecase}
}

type updateMFARequest struct {
	Enforce2FA bool `json:"enforce_2fa"`
}

// GetSettings godoc
// @Summary      Get MFA settings
// @Description  Retrieve the current system-wide MFA enforcement policy.
// @Tags         Admin
// @Produce      json
// @Success      200 {object} response.SuccessSingleUserResponse
// @Router       /admin/mfa-settings [get]
func (h *MFASettingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.usecase.GetSettings()
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", settings)
}

// UpdateSettings godoc
// @Summary      Update MFA settings
// @Description  Toggle system-wide 2FA enforcement. Admin only.
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        body body updateMFARequest true "MFA settings payload"
// @Success      200 {object} response.SuccessSingleUserResponse
// @Router       /admin/mfa-settings [put]
func (h *MFASettingsHandler) UpdateSettings(c *gin.Context) {
	var req updateMFARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "422", "invalid request", err.Error())
		return
	}

	adminID, _ := c.Get("user_id")

	if err := h.usecase.UpdateSettings(req.Enforce2FA, adminID.(int64)); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "MFA settings updated", nil)
}
