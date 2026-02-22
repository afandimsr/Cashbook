package bootstrap

import (
	httpDelivery "github.com/afandimsr/cashbook-backend/internal/delivery/http"
	"github.com/afandimsr/cashbook-backend/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userHandler *handler.UserHandler, categoryHandler *handler.CategoryHandler,
	transactionHandler *handler.TransactionHandler,
	budgetHandler *handler.BudgetHandler,
	reportHandler *handler.ReportHandler,
	recurringHandler *handler.RecurringHandler,
	twofaHandler *handler.TwoFAHandler,
	mfaSettingsHandler *handler.MFASettingsHandler,
) {
	httpDelivery.RegisterRoutes(r, userHandler, categoryHandler, transactionHandler, budgetHandler, reportHandler, recurringHandler, twofaHandler, mfaSettingsHandler)
}
