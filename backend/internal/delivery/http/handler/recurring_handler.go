package handler

import (
	"net/http"
	"strconv"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/recurring_transaction"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/recurring_transaction"
	"github.com/gin-gonic/gin"
)

type RecurringHandler struct {
	usecase uc.Usecase
}

func NewRecurringHandler(usecase uc.Usecase) *RecurringHandler {
	return &RecurringHandler{usecase: usecase}
}

// GetRecurring godoc
// @Summary      Get recurring transactions
// @Description  Get all recurring transactions for the authenticated user
// @Tags         Recurring
// @Produce      json
// @Success      200 {object} response.SuccessRecurringResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /recurring [get]
func (h *RecurringHandler) GetRecurring(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	rts, err := h.usecase.GetRecurring(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", rts)
}

// CreateRecurring godoc
// @Summary      Create recurring transaction
// @Description  Create a new recurring transaction template
// @Tags         Recurring
// @Accept       json
// @Produce      json
// @Param        body body recurring_transaction.RecurringTransaction true "Recurring Transaction payload"
// @Success      201 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /recurring [post]
func (h *RecurringHandler) CreateRecurring(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req recurring_transaction.RecurringTransaction
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request payload", err.Error())
		return
	}

	rt, err := h.usecase.CreateRecurring(userID, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to create recurring transaction", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "recurring transaction created", rt)
}

// DeleteRecurring godoc
// @Summary      Delete recurring transaction
// @Description  Delete an existing recurring transaction template
// @Tags         Recurring
// @Produce      json
// @Param        id   path      int  true  "Recurring Transaction ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /recurring/{id} [delete]
func (h *RecurringHandler) DeleteRecurring(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid id", err))
		return
	}

	if err := h.usecase.DeleteRecurring(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "deleted", nil)
}

// ProcessDue godoc
// @Summary      Process due recurring transactions
// @Description  Trigger the manual processing of due recurring transactions into transactions
// @Tags         Recurring
// @Produce      json
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /recurring/process [post]
func (h *RecurringHandler) ProcessDue(c *gin.Context) {
	if err := h.usecase.ProcessDueTransactions(); err != nil {
		response.Error(c, 400, "BAD_REQUEST", "failed to process due recurring transactions", err.Error())
		return
	}

	response.Success(c, http.StatusOK, "processed", nil)
}
