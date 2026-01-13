package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	"github.com/afandimsr/cashbook-backend/internal/domain/apperror"
	"github.com/afandimsr/cashbook-backend/internal/domain/budget"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/budget"
	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	usecase uc.Usecase
}

func NewBudgetHandler(usecase uc.Usecase) *BudgetHandler {
	return &BudgetHandler{usecase: usecase}
}

// GetBudgets godoc
// @Summary      Get budgets
// @Description  Get budgets for a specific month and year
// @Tags         Budgets
// @Produce      json
// @Param        month  query     int  false  "Month (1-12)"
// @Param        year   query     int  false  "Year"
// @Success      200 {object} response.SuccessBudgetResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /budgets [get]
func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	now := time.Now()
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(now.Year()))

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	budgets, err := h.usecase.GetBudgets(userID, month, year)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", budgets)
}

// SetBudget godoc
// @Summary      Set or update budget
// @Description  Set or update a budget for a category, month, and year
// @Tags         Budgets
// @Accept       json
// @Produce      json
// @Param        body body budget.Budget true "Budget payload"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorSwaggerResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /budgets [post]
func (h *BudgetHandler) SetBudget(c *gin.Context) {
	userID := c.MustGet("userID").(int64)

	var req budget.Budget
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.BadRequest("invalid request", err))
		return
	}

	if err := h.usecase.SetBudget(userID, req); err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "budget updated", nil)
}
