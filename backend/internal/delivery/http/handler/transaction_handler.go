package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/transaction"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	usecase uc.Usecase
}

func NewTransactionHandler(usecase uc.Usecase) *TransactionHandler {
	return &TransactionHandler{
		usecase: usecase,
	}
}

// GetTransactions godoc
// @Summary      Browse financial transactions
// @Description  Retrieve a detailed list of financial transactions with support for pagination and business filters like category, date range, and type (income/expense).
// @Tags         Transactions
// @Produce      json
// @Param        page   query     int     false  "Page number"
// @Param        limit  query     int     false  "Items per page"
// @Param        q      query     string  false  "Search query"
// @Success      200 {object} response.SuccessTransactionResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := transaction.Filter{
		Search: c.Query("q"),
		Type:   c.Query("type"),
	}

	if catID := c.Query("category_id"); catID != "" {
		id, _ := strconv.ParseInt(catID, 10, 64)
		filter.CategoryID = id
	}

	if sd := c.Query("start_date"); sd != "" {
		if t, err := time.Parse("2006-01-02", sd); err == nil {
			filter.StartDate = t
		}
	}

	if ed := c.Query("end_date"); ed != "" {
		if t, err := time.Parse("2006-01-02", ed); err == nil {
			filter.EndDate = t
		}
	}

	txs, err := h.usecase.GetAllByUserID(userID, page, limit, filter)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", txs)
}

// CreateTransaction godoc
// @Summary      Record a new transaction
// @Description  Log a new financial entry (income or expense) to track personal cash flow.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        body body transaction.Transaction true "Transaction payload"
// @Success      201 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// Accept date as either RFC3339 datetime or date-only "2006-01-02".
	var raw struct {
		ID         int64   `json:"id"`
		UserID     int64   `json:"user_id"`
		CategoryID int64   `json:"category_id"`
		Amount     float64 `json:"amount"`
		Note       string  `json:"note"`
		Date       string  `json:"date"`
		Type       string  `json:"type"`
	}

	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request", err.Error())
		return
	}

	var parsedDate time.Time
	var parseErr error
	if raw.Date == "" {
		parsedDate = time.Time{}
	} else {
		layouts := []string{time.RFC3339, "2006-01-02"}
		for _, l := range layouts {
			parsedDate, parseErr = time.Parse(l, raw.Date)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_DATE", "invalid date format", parseErr.Error())
			return
		}
	}

	req := transaction.Transaction{
		ID:         raw.ID,
		UserID:     c.MustGet("user_id").(int64),
		CategoryID: raw.CategoryID,
		Amount:     raw.Amount,
		Note:       raw.Note,
		Date:       parsedDate,
		Type:       raw.Type,
	}

	if err := h.usecase.Create(req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "transaction created", nil)
}

// GetSummary godoc
// @Summary      Dashboard financial health summary
// @Description  Calculate and retrieve the current total balance, aggregate income, and aggregate expenses for an immediate overview of financial status.
// @Tags         Transactions
// @Produce      json
// @Success      200 {object} response.SuccessSummaryResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /transactions/summary [get]
func (h *TransactionHandler) GetSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	summary, err := h.usecase.GetDashboardSummary(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", summary)
}

// UpdateTransaction godoc
// @Summary      Amend a financial record
// @Description  Update the details of an existing transaction to ensure ledger accuracy.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Transaction ID"
// @Param        body body transaction.Transaction true "Transaction payload"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid id", err))
		return
	}

	var raw struct {
		ID         int64   `json:"id"`
		UserID     int64   `json:"user_id"`
		CategoryID int64   `json:"category_id"`
		Amount     float64 `json:"amount"`
		Note       string  `json:"note"`
		Date       string  `json:"date"`
		Type       string  `json:"type"`
	}

	if err := c.ShouldBindJSON(&raw); err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_REQUEST", "invalid request", err.Error())
		return
	}

	var parsedDate time.Time
	var parseErr error
	if raw.Date == "" {
		parsedDate = time.Time{}
	} else {
		layouts := []string{time.RFC3339, "2006-01-02"}
		for _, l := range layouts {
			parsedDate, parseErr = time.Parse(l, raw.Date)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			response.Error(c, http.StatusBadRequest, "INVALID_DATE", "invalid date format", parseErr.Error())
			return
		}
	}

	req := transaction.Transaction{
		ID:         raw.ID,
		UserID:     raw.UserID,
		CategoryID: raw.CategoryID,
		Amount:     raw.Amount,
		Note:       raw.Note,
		Date:       parsedDate,
		Type:       raw.Type,
	}

	if err := h.usecase.Update(id, req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "transaction updated", nil)
}

// DeleteTransaction godoc
// @Summary      Remove a financial record
// @Description  Permanently delete a transaction from the personal finance history.
// @Tags         Transactions
// @Produce      json
// @Param        id   path      int  true  "Transaction ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      404 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(apperror.BadRequest("invalid id", err))
		return
	}

	if err := h.usecase.Delete(id); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "transaction deleted", nil)
}
