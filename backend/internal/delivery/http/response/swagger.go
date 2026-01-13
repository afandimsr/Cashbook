package response

import (
	"github.com/afandimsr/cashbook-backend/internal/domain/budget"
	"github.com/afandimsr/cashbook-backend/internal/domain/category"
	"github.com/afandimsr/cashbook-backend/internal/domain/recurring_transaction"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	"github.com/afandimsr/cashbook-backend/internal/domain/user"
)

// Generic success response for swagger
type SuccessUserResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"success"`
	Data    []user.User `json:"data"`
}

type SuccessSingleUserResponse struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"success"`
	Data    user.User `json:"data"`
}

type SuccessCategoryResponse struct {
	Success bool                `json:"success" example:"true"`
	Message string              `json:"message" example:"success"`
	Data    []category.Category `json:"data"`
}

type SuccessTransactionResponse struct {
	Success bool                      `json:"success" example:"true"`
	Message string                    `json:"message" example:"success"`
	Data    []transaction.Transaction `json:"data"`
}

type SuccessBudgetResponse struct {
	Success bool            `json:"success" example:"true"`
	Message string          `json:"message" example:"success"`
	Data    []budget.Budget `json:"data"`
}

type SuccessRecurringResponse struct {
	Success bool                                         `json:"success" example:"true"`
	Message string                                       `json:"message" example:"success"`
	Data    []recurring_transaction.RecurringTransaction `json:"data"`
}

type SuccessSummaryResponse struct {
	Success bool                         `json:"success" example:"true"`
	Message string                       `json:"message" example:"success"`
	Data    transaction.DashboardSummary `json:"data"`
}

type ErrorSwaggerResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"error"`
	Errors  string `json:"errors"`
}
