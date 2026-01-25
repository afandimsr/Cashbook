package bootstrap

import (
	"log"

	_ "github.com/afandimsr/cashbook-backend/docs"
	"github.com/afandimsr/cashbook-backend/internal/config"
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/middleware"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/auth"
	"github.com/afandimsr/cashbook-backend/internal/infrastructure/external"
	repo "github.com/afandimsr/cashbook-backend/internal/infrastructure/persistent/postgresql/repository"
	"github.com/afandimsr/cashbook-backend/internal/pkg/jwt"
	budgetUC "github.com/afandimsr/cashbook-backend/internal/usecase/budget"
	categoryUC "github.com/afandimsr/cashbook-backend/internal/usecase/category"
	recurringUC "github.com/afandimsr/cashbook-backend/internal/usecase/recurring_transaction"
	reportUC "github.com/afandimsr/cashbook-backend/internal/usecase/report"
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

	db, err := config.NewPostgreSQL(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	authClient := external.NewAuthClient(cfg.ClientAuthURL)
	googleAuth := auth.NewGoogleAuth(cfg)

	userRepository := repo.NewUserRepo(db)
	categoryRepository := repo.NewCategoryRepo(db)
	transactionRepository := repo.NewTransactionRepo(db)
	budgetRepository := repo.NewBudgetRepo(db)
	recurringRepository := repo.NewRecurringRepo(db)

	userUsecase := userUC.New(userRepository, authClient)
	oauthUsecase := userUC.NewOAuthUsecase(userRepository, googleAuth)
	categoryUsecase := categoryUC.New(categoryRepository)
	transactionUsecase := transactionUC.New(transactionRepository)
	budgetUsecase := budgetUC.New(budgetRepository)
	reportUsecase := reportUC.New(transactionRepository)
	recurringUsecase := recurringUC.New(recurringRepository, transactionRepository)

	userHandler := handler.New(userUsecase, oauthUsecase)
	categoryHandler := handler.NewCategoryHandler(categoryUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)
	budgetHandler := handler.NewBudgetHandler(budgetUsecase)
	reportHandler := handler.NewReportHandler(reportUsecase)
	recurringHandler := handler.NewRecurringHandler(recurringUsecase)

	r := gin.Default()
	r.Use(cors.New(middleware.Cors(cfg)))
	r.Use(middleware.ErrorHandler())

	RegisterRoutes(r, userHandler, categoryHandler, transactionHandler, budgetHandler, reportHandler, recurringHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("Running on port", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
