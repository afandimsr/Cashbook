package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/delivery/http/middleware"
	"github.com/afandimsr/cashbook-backend/internal/domain/shared_expense"
	usecase "github.com/afandimsr/cashbook-backend/internal/usecase/shared_expense"
	"github.com/gin-gonic/gin"
)

type SharedExpenseHandler struct {
	usecase usecase.Usecase
}

func NewSharedExpenseHandler(u usecase.Usecase) *SharedExpenseHandler {
	return &SharedExpenseHandler{usecase: u}
}

type CreateSplitRequest struct {
	Title          string               `json:"title" binding:"required"`
	Subtotal       float64              `json:"subtotal"`
	TaxAmount      float64              `json:"tax_amount"`
	ServiceCharge  float64              `json:"service_charge"`
	OtherCharge    float64              `json:"other_charge"`
	DiscountAmount float64              `json:"discount_amount"`
	CategoryID     int64                `json:"category_id" binding:"required"`
	Date           time.Time            `json:"date" binding:"required"`
	PayerID        int64                `json:"payer_id" binding:"required"`
	SplitMethod    string               `json:"split_method"`
	Participants   []ParticipantRequest `json:"participants" binding:"required,min=1"`
	Items          []ItemRequest        `json:"items"`
}

type ParticipantRequest struct {
	UserID     *int64  `json:"user_id"`
	ShadowName string  `json:"shadow_name"`
	BaseShare  float64 `json:"base_share"`  // used by the EXACT method
	Percentage float64 `json:"percentage"`  // used by the PERCENTAGE method
}

type ItemRequest struct {
	Name               string  `json:"name"`
	Price              float64 `json:"price"`
	Quantity           int     `json:"quantity"`
	CategoryID         *int64  `json:"category_id"`
	ParticipantIndexes []int   `json:"participant_indexes"`
}

func (h *SharedExpenseHandler) Create(c *gin.Context) {
	var req CreateSplitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(int64)

	bill := shared_expense.SplitBill{
		CreatorID:      userID,
		PayerID:        req.PayerID,
		CategoryID:     req.CategoryID,
		Title:          req.Title,
		Subtotal:       req.Subtotal,
		TaxAmount:      req.TaxAmount,
		ServiceCharge:  req.ServiceCharge,
		OtherCharge:    req.OtherCharge,
		DiscountAmount: req.DiscountAmount,
		Date:           req.Date,
		SplitMethod:    shared_expense.SplitMethod(req.SplitMethod),
	}

	for _, p := range req.Participants {
		bill.Participants = append(bill.Participants, shared_expense.SplitParticipant{
			UserID:     p.UserID,
			ShadowName: p.ShadowName,
			BaseShare:  p.BaseShare,
			Percentage: p.Percentage,
		})
	}

	for _, it := range req.Items {
		bill.Items = append(bill.Items, shared_expense.SplitItem{
			Name:               it.Name,
			Price:              it.Price,
			Quantity:           it.Quantity,
			CategoryID:         it.CategoryID,
			ParticipantIndexes: it.ParticipantIndexes,
		})
	}

	err := h.usecase.CreateSplitBill(c.Request.Context(), bill)
	if err != nil {
		if errors.Is(err, shared_expense.ErrInvalidSplit) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Split bill created successfully"})
}

func (h *SharedExpenseHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	// In a real app, we'd parse filters from query params here
	filter := shared_expense.Filter{}

	bills, err := h.usecase.GetByUserID(c.Request.Context(), userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bills})
}

func (h *SharedExpenseHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	
	bill, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bill})
}

func (h *SharedExpenseHandler) Settle(c *gin.Context) {
	billID := c.Param("id")
	participantID := c.Param("participant_id")
	userID := c.MustGet("user_id").(int64)

	err := h.usecase.SettleParticipant(c.Request.Context(), billID, participantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Participant marked as paid"})
}

func (h *SharedExpenseHandler) GetSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	summary, err := h.usecase.GetSummary(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

func (h *SharedExpenseHandler) RegisterRoutes(r *gin.RouterGroup) {
	splits := r.Group("/splits")
	splits.Use(middleware.AuthMiddleware())
	{
		splits.POST("", h.Create)
		splits.GET("", h.GetAll)
		splits.GET("/summary", h.GetSummary)
		splits.GET("/:id", h.GetByID)
		splits.POST("/:id/settle/:participant_id", h.Settle)
	}
}
