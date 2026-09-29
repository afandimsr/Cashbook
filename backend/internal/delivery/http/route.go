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
	sharedExpenseHandler *handler.SharedExpenseHandler,
	botHandler *handler.BotHandler,
	botServiceAuth gin.HandlerFunc,
	loginRateLimit gin.HandlerFunc,
	twoFARateLimit gin.HandlerFunc,
	linkCodeRateLimit gin.HandlerFunc,
) {
	api := r.Group("/api/v1")

	// auth routes (public)
	api.POST("/login", loginRateLimit, userHandler.Login)
	api.GET("/auth/google/login", userHandler.GoogleLogin)
	api.GET("/auth/google/callback", userHandler.GoogleCallback)

	// 2FA routes (public — used during login)
	api.POST("/2fa/verify", twoFARateLimit, twofaHandler.VerifyLogin)
	api.POST("/2fa/backup/verify", twoFARateLimit, twofaHandler.VerifyBackupCode)

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

		// admin managing another user's Telegram link (self-service is under /telegram below)
		users.POST("/:id/telegram/link-code", botHandler.AdminGenerateLinkCode)
		users.GET("/:id/telegram/status", botHandler.AdminGetLinkStatus)
		users.DELETE("/:id/telegram/link", botHandler.AdminUnlink)
	}

	// admin MFA settings (protected + admin only)
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminOnly())
	{
		admin.GET("/mfa-settings", mfaSettingsHandler.GetSettings)
		admin.PUT("/mfa-settings", mfaSettingsHandler.UpdateSettings)
		admin.GET("/telegram/links", botHandler.AdminListLinks)
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

	// shared expense routes
	sharedExpenseHandler.RegisterRoutes(api)

	// telegram link-code generation (protected, called from the app by an end-user)
	telegramRoutes := api.Group("/telegram")
	telegramRoutes.Use(middleware.AuthMiddleware(), middleware.RoleGuard("ADMIN", "USER"))
	{
		telegramRoutes.POST("/link-code", linkCodeRateLimit, botHandler.GenerateLinkCode)
		telegramRoutes.GET("/status", botHandler.GetLinkStatus)
		telegramRoutes.DELETE("/link", botHandler.Unlink)
	}

	// internal bot routes (service-to-service, called only by the Telegram bot)
	internalBot := api.Group("/internal/bot")
	internalBot.Use(botServiceAuth)
	{
		internalBot.POST("/link", botHandler.Link)
		internalBot.GET("/links", botHandler.ListLinks)
		internalBot.GET("/context", botHandler.GetContext)
		internalBot.POST("/transactions", botHandler.CreateTransaction)
		internalBot.GET("/summary", botHandler.GetSummary)
		internalBot.GET("/reports/spending", botHandler.GetCategorySpending)
		internalBot.GET("/reports/monthly", botHandler.GetMonthlyReport)
		internalBot.GET("/budgets", botHandler.GetBudgetStatus)
	}
}

func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
