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
	twofaHandler *handler.TwoFAHandler,
	mfaSettingsHandler *handler.MFASettingsHandler,
) {
	api := r.Group("/api/v1")

	// auth routes (public)
	api.POST("/login", userHandler.Login)
	api.GET("/auth/google/login", userHandler.GoogleLogin)
	api.GET("/auth/google/callback", userHandler.GoogleCallback)

	// 2FA routes (public — used during login)
	api.POST("/2fa/verify", twofaHandler.VerifyLogin)
	api.POST("/2fa/backup/verify", twofaHandler.VerifyBackupCode)

	// health check
	api.GET("/health", healthHandler)

	// 2FA routes (authenticated — for setup/management)
	twofa := api.Group("/2fa")
	twofa.Use(middleware.AuthMiddleware())
	{
		twofa.POST("/setup", twofaHandler.Setup)
		twofa.POST("/setup/verify", twofaHandler.VerifySetup)
		twofa.DELETE("/disable", twofaHandler.Disable)
		twofa.POST("/backup-codes", twofaHandler.GenerateBackupCodes)
	}

	// user routes (protected)
	users := api.Group("/users")
	users.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		users.GET("", userHandler.GetUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.POST("/:id/reset-password", userHandler.ResetPassword)
	}

	// admin MFA settings (protected + admin only)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("/mfa-settings", mfaSettingsHandler.GetSettings)
		admin.PUT("/mfa-settings", mfaSettingsHandler.UpdateSettings)
	}

	// user MFA settings (protected + admin only) - alternative route
	userRoutes := api.Group("/user")
	userRoutes.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		userRoutes.GET("/mfa-settings", mfaSettingsHandler.GetSettings)
		userRoutes.PUT("/mfa-settings", mfaSettingsHandler.UpdateSettings)
	}

	// category routes
	categories := api.Group("/categories")
	categories.Use(middleware.AuthMiddleware(), middleware.RoleGuard("ADMIN", "USER"))
	{
		categories.GET("", categoryHandler.GetCategories)
		categories.POST("", categoryHandler.CreateCategory)
		categories.PUT("/:id", categoryHandler.UpdateCategory)
		categories.DELETE("/:id", categoryHandler.DeleteCategory)
	}

	// transaction routes
	transactions := api.Group("/transactions")
	transactions.Use(middleware.AuthMiddleware(), middleware.RoleGuard("ADMIN", "USER"))
	{
		transactions.GET("", transactionHandler.GetTransactions)
		transactions.POST("", transactionHandler.CreateTransaction)
		transactions.GET("/summary", transactionHandler.GetSummary)
		transactions.PUT("/:id", transactionHandler.UpdateTransaction)
		transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
	}

	// budget routes
	budgets := api.Group("/budgets")
	budgets.Use(middleware.AuthMiddleware(), middleware.RoleGuard("ADMIN", "USER"))
	{
		budgets.GET("", budgetHandler.GetBudgets)
		budgets.POST("", budgetHandler.SetBudget)
	}

	// report routes
	reports := api.Group("/reports")
	reports.Use(middleware.AuthMiddleware(), middleware.RoleGuard("ADMIN", "USER"))
	{
		reports.GET("/spending", reportHandler.GetCategorySpending)
	}

	// recurring routes
	recurring := api.Group("/recurring")
	recurring.Use(middleware.AuthMiddleware(), middleware.RoleGuard("ADMIN", "USER"))
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
