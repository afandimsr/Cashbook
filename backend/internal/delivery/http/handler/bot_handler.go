package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/audit"
	"github.com/afandimsr/cashbook-backend/internal/domain/telegram"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/bot"
	"github.com/gin-gonic/gin"
)

// BotHandler serves two audiences on two different auth mechanisms:
//   - GenerateLinkCode is called by an authenticated end-user (normal JWT) from the app.
//   - Everything else is called by the Telegram bot service (ServiceAuthMiddleware)
//     and identifies the acting user only by telegram_chat_id, never a trusted user_id.
type BotHandler struct {
	usecase   uc.Usecase
	auditRepo audit.Repository
}

func NewBotHandler(usecase uc.Usecase, auditRepo audit.Repository) *BotHandler {
	return &BotHandler{usecase: usecase, auditRepo: auditRepo}
}

// logAudit records an admin action against a target user. Best-effort: a
// broken audit log must never block the admin action itself, so failures are
// only logged to stdout, not surfaced to the caller.
func (h *BotHandler) logAudit(c *gin.Context, action string, targetUserID int64) {
	adminID, ok := c.Get("user_id")
	if !ok {
		return
	}
	entry := &audit.Entry{
		AdminUserID:  adminID.(int64),
		Action:       action,
		TargetUserID: &targetUserID,
	}
	if err := h.auditRepo.Save(c.Request.Context(), entry); err != nil {
		log.Printf("audit log failed (action=%s, target_user_id=%d): %v", action, targetUserID, err)
	}
}

// GenerateLinkCode godoc
// @Summary      Generate a Telegram account-linking code
// @Description  Create a short-lived one-time code the user sends to the CashBook bot to link their Telegram account.
// @Tags         Telegram
// @Produce      json
// @Success      201 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /telegram/link-code [post]
func (h *BotHandler) GenerateLinkCode(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	code, err := h.usecase.GenerateLinkCode(c.Request.Context(), userID)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusCreated, "link code generated", gin.H{
		"code":       code.Code,
		"expires_at": code.ExpiresAt,
	})
}

// GetLinkStatus godoc
// @Summary      Check the caller's Telegram link status
// @Description  Returns whether the authenticated user currently has a Telegram account linked.
// @Tags         Telegram
// @Produce      json
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Router       /telegram/status [get]
func (h *BotHandler) GetLinkStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	status, err := h.usecase.GetLinkStatus(c.Request.Context(), userID)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", status)
}

// Unlink godoc
// @Summary      Disconnect the caller's Telegram account
// @Description  Revokes the authenticated user's own Telegram link.
// @Tags         Telegram
// @Produce      json
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Router       /telegram/link [delete]
func (h *BotHandler) Unlink(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	if err := h.usecase.Unlink(c.Request.Context(), userID); err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "telegram account unlinked", nil)
}

// AdminGenerateLinkCode godoc
// @Summary      [Admin] Generate a Telegram link code for another user
// @Description  Same as GenerateLinkCode but for any user_id, admin-only — mirrors the /users/{id}/reset-password pattern.
// @Tags         Telegram
// @Produce      json
// @Param        id   path      int  true  "Target user ID"
// @Success      201 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      403 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Router       /users/{id}/telegram/link-code [post]
func (h *BotHandler) AdminGenerateLinkCode(c *gin.Context) {
	userID, err := parseUserIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	code, err := h.usecase.GenerateLinkCode(c.Request.Context(), userID)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	h.logAudit(c, "telegram.generate_link_code", userID)

	response.Success(c, http.StatusCreated, "link code generated", gin.H{
		"code":       code.Code,
		"expires_at": code.ExpiresAt,
	})
}

// AdminGetLinkStatus godoc
// @Summary      [Admin] Check another user's Telegram link status
// @Tags         Telegram
// @Produce      json
// @Param        id   path      int  true  "Target user ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      403 {object} response.ErrorSwaggerResponse
// @Router       /users/{id}/telegram/status [get]
func (h *BotHandler) AdminGetLinkStatus(c *gin.Context) {
	userID, err := parseUserIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	status, err := h.usecase.GetLinkStatus(c.Request.Context(), userID)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", status)
}

// AdminUnlink godoc
// @Summary      [Admin] Disconnect another user's Telegram account
// @Tags         Telegram
// @Produce      json
// @Param        id   path      int  true  "Target user ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      403 {object} response.ErrorSwaggerResponse
// @Router       /users/{id}/telegram/link [delete]
func (h *BotHandler) AdminUnlink(c *gin.Context) {
	userID, err := parseUserIDParam(c)
	if err != nil {
		c.Error(err)
		return
	}

	if err := h.usecase.Unlink(c.Request.Context(), userID); err != nil {
		c.Error(mapBotError(err))
		return
	}

	h.logAudit(c, "telegram.unlink", userID)

	response.Success(c, http.StatusOK, "telegram account unlinked", nil)
}

// AdminListLinks godoc
// @Summary      [Admin] List every active Telegram link
// @Tags         Telegram
// @Produce      json
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      403 {object} response.ErrorSwaggerResponse
// @Router       /admin/telegram/links [get]
func (h *BotHandler) AdminListLinks(c *gin.Context) {
	links, err := h.usecase.ListLinkedChats(c.Request.Context())
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", links)
}

// Link godoc
// @Summary      [Internal] Link a Telegram chat to a user
// @Description  Consumes a one-time code generated in the app to bind a telegram_chat_id to a user_id. Internal service auth only.
// @Tags         Telegram
// @Accept       json
// @Produce      json
// @Success      201 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Router       /internal/bot/link [post]
func (h *BotHandler) Link(c *gin.Context) {
	var req struct {
		ChatID   int64  `json:"chat_id" binding:"required"`
		Username string `json:"username"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	link, err := h.usecase.LinkAccount(c.Request.Context(), req.ChatID, req.Username, req.Code)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusCreated, "telegram account linked", link)
}

// ListLinks godoc
// @Summary      [Internal] List every active Telegram link
// @Description  Returns every user_telegram_links row with is_active=true, used by the bot's monthly report scheduler to know who to message. Internal service auth only.
// @Tags         Telegram
// @Produce      json
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Router       /internal/bot/links [get]
func (h *BotHandler) ListLinks(c *gin.Context) {
	links, err := h.usecase.ListLinkedChats(c.Request.Context())
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", links)
}

// GetContext godoc
// @Summary      [Internal] Resolve a chat and its categories
// @Description  Resolves telegram_chat_id to a user and returns their custom categories in one call. Internal service auth only.
// @Tags         Telegram
// @Produce      json
// @Param        chat_id query int true "Telegram chat ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Router       /internal/bot/context [get]
func (h *BotHandler) GetContext(c *gin.Context) {
	chatID, err := parseChatID(c)
	if err != nil {
		c.Error(err)
		return
	}

	ctxData, err := h.usecase.GetContext(c.Request.Context(), chatID)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", ctxData)
}

// CreateTransaction godoc
// @Summary      [Internal] Record a transaction on behalf of a linked user
// @Description  Resolves telegram_chat_id to a user and creates a transaction. Internal service auth only.
// @Tags         Telegram
// @Accept       json
// @Produce      json
// @Success      201 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Router       /internal/bot/transactions [post]
func (h *BotHandler) CreateTransaction(c *gin.Context) {
	var req struct {
		ChatID     int64   `json:"chat_id" binding:"required"`
		CategoryID int64   `json:"category_id" binding:"required"`
		Amount     float64 `json:"amount" binding:"required"`
		Type       string  `json:"type" binding:"required"`
		Note       string  `json:"note"`
		Date       string  `json:"date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	var date time.Time
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.Error(apperror.BadRequest("invalid date format, expected YYYY-MM-DD", err))
			return
		}
		date = parsed
	}

	err := h.usecase.CreateTransaction(c.Request.Context(), uc.CreateTransactionInput{
		ChatID:     req.ChatID,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Type:       req.Type,
		Note:       req.Note,
		Date:       date,
	})
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusCreated, "transaction created", nil)
}

// GetSummary godoc
// @Summary      [Internal] Dashboard summary for a linked user
// @Description  Resolves telegram_chat_id to a user and returns their dashboard totals. Internal service auth only.
// @Tags         Telegram
// @Produce      json
// @Param        chat_id query int true "Telegram chat ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Router       /internal/bot/summary [get]
func (h *BotHandler) GetSummary(c *gin.Context) {
	chatID, err := parseChatID(c)
	if err != nil {
		c.Error(err)
		return
	}

	summary, err := h.usecase.GetSummary(c.Request.Context(), chatID)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", summary)
}

// GetCategorySpending godoc
// @Summary      [Internal] Category spending breakdown for a linked user
// @Description  Resolves telegram_chat_id to a user and returns their per-category spending for the given month/year. Internal service auth only.
// @Tags         Telegram
// @Produce      json
// @Param        chat_id query int true "Telegram chat ID"
// @Param        month   query int false "Month (1-12)"
// @Param        year    query int false "Year"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Router       /internal/bot/reports/spending [get]
func (h *BotHandler) GetCategorySpending(c *gin.Context) {
	chatID, err := parseChatID(c)
	if err != nil {
		c.Error(err)
		return
	}

	now := time.Now()
	month, _ := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))

	report, err := h.usecase.GetCategorySpending(c.Request.Context(), chatID, month, year)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", report)
}

// GetMonthlyReport godoc
// @Summary      [Internal] Monthly income and expense breakdown for a linked user
// @Description  Resolves telegram_chat_id to a user and returns their per-category income and expense totals for the given month/year. Internal service auth only.
// @Tags         Telegram
// @Produce      json
// @Param        chat_id query int true "Telegram chat ID"
// @Param        month   query int false "Month (1-12)"
// @Param        year    query int false "Year"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Router       /internal/bot/reports/monthly [get]
func (h *BotHandler) GetMonthlyReport(c *gin.Context) {
	chatID, err := parseChatID(c)
	if err != nil {
		c.Error(err)
		return
	}

	month, year, err := parseMonthYear(c)
	if err != nil {
		c.Error(err)
		return
	}

	report, err := h.usecase.GetMonthlyReport(c.Request.Context(), chatID, month, year)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", report)
}

// GetBudgetStatus godoc
// @Summary      [Internal] Budget usage for a linked user
// @Description  Resolves telegram_chat_id to a user and returns each category budget for the given month/year with its spend, remaining amount, and percentage used. Internal service auth only.
// @Tags         Telegram
// @Produce      json
// @Param        chat_id query int true "Telegram chat ID"
// @Param        month   query int false "Month (1-12)"
// @Param        year    query int false "Year"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Router       /internal/bot/budgets [get]
func (h *BotHandler) GetBudgetStatus(c *gin.Context) {
	chatID, err := parseChatID(c)
	if err != nil {
		c.Error(err)
		return
	}

	month, year, err := parseMonthYear(c)
	if err != nil {
		c.Error(err)
		return
	}

	statuses, err := h.usecase.GetBudgetStatus(c.Request.Context(), chatID, month, year)
	if err != nil {
		c.Error(mapBotError(err))
		return
	}

	response.Success(c, http.StatusOK, "success", statuses)
}

// parseMonthYear reads optional month/year query params, defaulting to the
// current month and rejecting out-of-range values with a 400.
func parseMonthYear(c *gin.Context) (int, int, error) {
	now := time.Now()
	month, err := strconv.Atoi(c.DefaultQuery("month", strconv.Itoa(int(now.Month()))))
	if err != nil || month < 1 || month > 12 {
		return 0, 0, apperror.BadRequest("invalid month, expected 1-12", err)
	}
	year, err := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(now.Year())))
	if err != nil || year < 1 {
		return 0, 0, apperror.BadRequest("invalid year", err)
	}
	return month, year, nil
}

func parseChatID(c *gin.Context) (int64, error) {
	chatID, err := strconv.ParseInt(c.Query("chat_id"), 10, 64)
	if err != nil {
		return 0, apperror.BadRequest("invalid or missing chat_id", err)
	}
	return chatID, nil
}

func parseUserIDParam(c *gin.Context) (int64, error) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return 0, apperror.BadRequest("invalid user id", err)
	}
	return userID, nil
}

// mapBotError converts bot-domain sentinel errors into the right HTTP status,
// mirroring the sentinel -> errors.Is -> 4xx convention used elsewhere (see
// shared_expense_handler.go).
func mapBotError(err error) error {
	switch {
	case errors.Is(err, telegram.ErrLinkNotFound):
		return apperror.NotFound("telegram account is not linked to a user", err)
	case errors.Is(err, telegram.ErrLinkCodeInvalid):
		return apperror.BadRequest("link code is invalid, expired, or already used", err)
	case errors.Is(err, telegram.ErrChatAlreadyLinked):
		return apperror.Conflict("this telegram account is already linked to another user", err)
	case errors.Is(err, telegram.ErrUserNotFound):
		return apperror.NotFound("user not found", err)
	case errors.Is(err, uc.ErrCategoryNotOwned):
		return apperror.BadRequest("category does not belong to the linked user", err)
	default:
		return err
	}
}
