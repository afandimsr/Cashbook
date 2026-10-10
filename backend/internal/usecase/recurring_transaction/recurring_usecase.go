package recurring_transaction

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/recurring_transaction"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

type Usecase interface {
	GetRecurring(userID int64) ([]recurring_transaction.RecurringTransaction, error)
	CreateRecurring(userID int64, rt recurring_transaction.RecurringTransaction) (recurring_transaction.RecurringTransaction, error)
	DeleteRecurring(id int64) error
	ProcessDueForUser(ctx context.Context, userID int64) (ProcessResult, error)
}

// ProcessedItem is one transaction generated from a due recurring template.
type ProcessedItem struct {
	RecurringID   int64                           `json:"recurring_id"`
	TransactionID int64                           `json:"transaction_id"`
	CategoryID    int64                           `json:"category_id"`
	Amount        float64                         `json:"amount"`
	Type          string                          `json:"type"`
	Note          string                          `json:"note"`
	Frequency     recurring_transaction.Frequency `json:"frequency"`
}

// ProcessResult reports what one ProcessDueForUser run generated. Failed
// counts templates that hit an unexpected error; they stay due and are
// retried on the next run.
type ProcessResult struct {
	Processed []ProcessedItem `json:"processed"`
	Failed    int             `json:"failed"`
}

type usecase struct {
	repo recurring_transaction.Repository
}

func New(repo recurring_transaction.Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) GetRecurring(userID int64) ([]recurring_transaction.RecurringTransaction, error) {
	return u.repo.FindAllByUserID(userID)
}

func (u *usecase) CreateRecurring(userID int64, rt recurring_transaction.RecurringTransaction) (recurring_transaction.RecurringTransaction, error) {
	rt.UserID = userID
	if rt.StartDate.IsZero() {
		rt.StartDate = time.Now()
	}
	if err := u.repo.Save(&rt); err != nil {
		return recurring_transaction.RecurringTransaction{}, err
	}
	return rt, nil
}

func (u *usecase) DeleteRecurring(id int64) error {
	return u.repo.Delete(id)
}

// ProcessDueForUser generates a transaction for each of userID's due recurring
// templates. A failure on one template is logged and counted, not fatal, so
// the rest still get processed.
func (u *usecase) ProcessDueForUser(ctx context.Context, userID int64) (ProcessResult, error) {
	now := time.Now()
	due, err := u.repo.FindDueByUserID(userID, now)
	if err != nil {
		return ProcessResult{}, err
	}

	result := ProcessResult{Processed: []ProcessedItem{}}
	for _, rt := range due {
		tx := &transaction.Transaction{
			UserID:     rt.UserID,
			CategoryID: rt.CategoryID,
			Amount:     rt.Amount,
			Note:       rt.Note + " (Auto-generated)",
			Date:       now,
			Type:       rt.Type,
		}

		err := u.repo.ProcessOne(ctx, rt.ID, tx, now)
		if errors.Is(err, recurring_transaction.ErrAlreadyProcessed) {
			continue
		}
		if err != nil {
			log.Printf("process recurring transaction %d: %v", rt.ID, err)
			result.Failed++
			continue
		}

		result.Processed = append(result.Processed, ProcessedItem{
			RecurringID:   rt.ID,
			TransactionID: tx.ID,
			CategoryID:    rt.CategoryID,
			Amount:        rt.Amount,
			Type:          rt.Type,
			Note:          rt.Note,
			Frequency:     rt.Frequency,
		})
	}

	return result, nil
}
