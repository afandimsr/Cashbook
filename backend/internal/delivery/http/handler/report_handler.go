package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/response"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/report"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	usecase uc.Usecase
}

func NewReportHandler(usecase uc.Usecase) *ReportHandler {
	return &ReportHandler{usecase: usecase}
}

// GetCategorySpending godoc
// @Summary      Get category spending report
// @Description  Get a breakdown of spending by category for a specific month and year
// @Tags         Reports
// @Produce      json
// @Param        month  query     int  false  "Month (1-12)"
// @Param        year   query     int  false  "Year"
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorSwaggerResponse
// @Failure      500 {object} response.ErrorSwaggerResponse
// @Router       /reports/spending [get]
func (h *ReportHandler) GetCategorySpending(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	now := time.Now()
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(now.Month())))
	yearStr := c.DefaultQuery("year", strconv.Itoa(now.Year()))

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	report, err := h.usecase.GetCategorySpending(userID, month, year)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "success", report)
}
