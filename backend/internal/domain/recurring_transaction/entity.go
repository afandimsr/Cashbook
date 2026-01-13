package recurring_transaction

import "time"

type Frequency string

const (
	Daily   Frequency = "daily"
	Weekly  Frequency = "weekly"
	Monthly Frequency = "monthly"
)

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
	FindDue(now time.Time) ([]RecurringTransaction, error)
	Save(rt *RecurringTransaction) error
	Update(rt *RecurringTransaction) error
	Delete(id int64) error
	UpdateLastProcessed(id int64, lastProcessed time.Time) error
}
