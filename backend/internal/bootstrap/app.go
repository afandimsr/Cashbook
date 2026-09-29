package bootstrap

import (
	"log"
	"time"

	_ "github.com/afandimsr/cashbook-backend/docs"
	"github.com/afandimsr/cashbook-backend/internal/config"
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/middleware"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/apm"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/auth"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/external"
	repo "github.com/afandimsr/cashbook-backend/internal/infrastructure/persistent/postgresql/repository"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
	"github.com/afandimsr/cashbook-backend/internal/pkg/ratelimit"
	botUC "github.com/afandimsr/cashbook-backend/internal/usecase/bot"
	budgetUC "github.com/afandimsr/cashbook-backend/internal/usecase/budget"
	categoryUC "github.com/afandimsr/cashbook-backend/internal/usecase/category"
	recurringUC "github.com/afandimsr/cashbook-backend/internal/usecase/recurring_transaction"
	reportUC "github.com/afandimsr/cashbook-backend/internal/usecase/report"
	sharedExpenseUC "github.com/afandimsr/cashbook-backend/internal/usecase/shared_expense"
	transactionUC "github.com/afandimsr/cashbook-backend/internal/usecase/transaction"
	userUC "github.com/afandimsr/cashbook-backend/internal/usecase/user"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Run() {
	cfg := config.Load()
	jwt.SetSecret(cfg.JWTSecret)

	// set gin mode
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := config.NewPostgreSQL(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	// initialize APM
	apm.Init(cfg)

	authClient := external.NewAuthClient(cfg.ClientAuthURL)
	googleAuth := auth.NewGoogleAuth(cfg)

	// Repositories
	userRepository := repo.NewUserRepo(db)
	oauthStateRepository := repo.NewOauthStateRepo(db)
	categoryRepository := repo.NewCategoryRepo(db)
	transactionRepository := repo.NewTransactionRepo(db)
	budgetRepository := repo.NewBudgetRepo(db)
	recurringRepository := repo.NewRecurringRepo(db)
	mfaSettingsRepository := repo.NewMFASettingsRepo(db)
	mfaBackupCodeRepository := repo.NewMFABackupCodeRepo(db)
	sharedExpenseRepository := repo.NewSharedExpenseRepo(db)
	telegramRepository := repo.NewTelegramRepo(db)
	auditRepository := repo.NewAuditRepo(db)

	// Use cases
	userUsecase := userUC.New(userRepository, authClient)
	userUsecase.SetMFASettingsRepo(mfaSettingsRepository)
	oauthUsecase := userUC.NewOAuthUsecase(userRepository, oauthStateRepository, googleAuth)
	categoryUsecase := categoryUC.New(categoryRepository)
	transactionUsecase := transactionUC.New(transactionRepository)
	budgetUsecase := budgetUC.New(budgetRepository)
	reportUsecase := reportUC.New(transactionRepository)
	recurringUsecase := recurringUC.New(recurringRepository, transactionRepository)
	twofaUsecase := userUC.NewTwoFAUsecase(userRepository, mfaBackupCodeRepository)
	mfaSettingsUsecase := userUC.NewMFASettingsUsecase(mfaSettingsRepository)
	sharedExpenseUsecase := sharedExpenseUC.New(sharedExpenseRepository)
	botUsecase := botUC.New(telegramRepository, categoryUsecase, transactionUsecase, reportUsecase, userUsecase, budgetUsecase)

	// Rate limiters. In-memory, per-process — fine for this app's current
	// single-instance deployment; a horizontally-scaled deployment would need
	// a shared store (e.g. Redis) instead. See ratelimit package docs.
	loginAccountLimiter := ratelimit.New(5, 15*time.Minute) // per-email: 5 consecutive failed logins / 15 min
	loginIPLimiter := ratelimit.New(10, time.Minute)        // per-IP on /login
	twoFAIPLimiter := ratelimit.New(10, time.Minute)        // per-IP on /2fa/verify + /2fa/backup/verify
	linkCodeIPLimiter := ratelimit.New(10, time.Minute)     // per-IP on POST /telegram/link-code
	loginRateLimit := middleware.RateLimit(loginIPLimiter, middleware.ClientIPKey)
	twoFARateLimit := middleware.RateLimit(twoFAIPLimiter, middleware.ClientIPKey)
	linkCodeRateLimit := middleware.RateLimit(linkCodeIPLimiter, middleware.ClientIPKey)

	// Handlers
	userHandler := handler.New(cfg, userUsecase, oauthUsecase, loginAccountLimiter)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)
	budgetHandler := handler.NewBudgetHandler(budgetUsecase)
	reportHandler := handler.NewReportHandler(reportUsecase)
	recurringHandler := handler.NewRecurringHandler(recurringUsecase)
	twofaHandler := handler.NewTwoFAHandler(twofaUsecase)
	mfaSettingsHandler := handler.NewMFASettingsHandler(mfaSettingsUsecase)
	sharedExpenseHandler := handler.NewSharedExpenseHandler(sharedExpenseUsecase)
	botHandler := handler.NewBotHandler(botUsecase, auditRepository)
	botServiceAuth := middleware.ServiceAuthMiddleware(cfg.BotInternalAPIKey)

	r := gin.Default()
	r.SetTrustedProxies(nil) // Trust proxies for ClientIP() to work behind Nginx
	r.Use(
		cors.New(middleware.Cors(cfg)),
		apm.GinMiddleware(), // Gin Elastic APM
		middleware.ErrorHandler(),
	)

	RegisterRoutes(r, userHandler, categoryHandler, transactionHandler, budgetHandler, reportHandler, recurringHandler, twofaHandler, mfaSettingsHandler, sharedExpenseHandler, botHandler, botServiceAuth, loginRateLimit, twoFARateLimit, linkCodeRateLimit)
	if gin.Mode() != gin.ReleaseMode {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	log.Println("Running on port", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
