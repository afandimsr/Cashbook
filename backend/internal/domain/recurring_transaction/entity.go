package recurring_transaction

import (
	"context"
	"errors"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

type Frequency string

const (
	Daily   Frequency = "daily"
	Weekly  Frequency = "weekly"
	Monthly Frequency = "monthly"
)

// ErrAlreadyProcessed is returned by ProcessOne when the template is no longer
// due at commit time, e.g. a concurrent request processed it first.
var ErrAlreadyProcessed = errors.New("recurring transaction already processed")

type RecurringTransaction struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	CategoryID    int64     `json:"category_id"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"` // "income" or "expense"
	Note          string    `json:"note"`
	Frequency     Frequency `json:"frequency"`
	StartDate     time.Time `json:"start_date"`
	LastProcessed time.Time `json:"last_processed"`
}

type Repository interface {
	FindAllByUserID(userID int64) ([]RecurringTransaction, error)
	FindDueByUserID(userID int64, now time.Time) ([]RecurringTransaction, error)
	Save(rt *RecurringTransaction) error
	Update(rt *RecurringTransaction) error
	Delete(id int64) error
	// ProcessOne atomically marks template id as processed at now and saves tx.
	// It returns ErrAlreadyProcessed when the template is not due anymore.
	ProcessOne(ctx context.Context, id int64, tx *transaction.Transaction, now time.Time) error
}
