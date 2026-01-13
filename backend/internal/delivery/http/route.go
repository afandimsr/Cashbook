package http

import (
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	r *gin.Engine,
	userHandler *handler.UserHandler,
	categoryHandler *handler.CategoryHandler,
	transactionHandler *handler.TransactionHandler,
	budgetHandler *handler.BudgetHandler,
	reportHandler *handler.ReportHandler,
	recurringHandler *handler.RecurringHandler,
) {
	api := r.Group("/api/v1")

	// auth routes
	api.POST("/login", userHandler.Login)
	api.GET("/auth/google/login", userHandler.GoogleLogin)
	api.GET("/auth/google/callback", userHandler.GoogleCallback)

	// health check
	api.GET("/health", healthHandler)

	// user routes (protected)
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		users.GET("", userHandler.GetUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
	}

	// category routes
	categories := api.Group("/categories")
	categories.Use(middleware.AuthMiddleware())
	{
		categories.GET("", categoryHandler.GetCategories)
		categories.POST("", categoryHandler.CreateCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// transaction routes
	transactions := api.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware())
	{
		transactions.GET("", transactionHandler.GetTransactions)
		transactions.POST("", transactionHandler.CreateTransaction)
		transactions.GET("/summary", transactionHandler.GetSummary)
		transactions.PUT("/:id", transactionHandler.UpdateTransaction)
		transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
	}

	// budget routes
	budgets := api.Group("/budgets")
	budgets.Use(middleware.AuthMiddleware())
	{
		budgets.GET("", budgetHandler.GetBudgets)
		budgets.POST("", budgetHandler.SetBudget)
	}

	// report routes
	reports := api.Group("/reports")
	reports.Use(middleware.AuthMiddleware())
	{
		reports.GET("/spending", reportHandler.GetCategorySpending)
	}

	// recurring routes
	recurring := api.Group("/recurring")
	recurring.Use(middleware.AuthMiddleware())
	{
		recurring.GET("", recurringHandler.GetRecurring)
		recurring.POST("", recurringHandler.CreateRecurring)
		recurring.DELETE("/:id", recurringHandler.DeleteRecurring)
		recurring.POST("/process", recurringHandler.ProcessDue)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
