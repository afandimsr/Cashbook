package recurring_transaction

import (
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/recurring_transaction"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

type Usecase interface {
	GetRecurring(userID int64) ([]recurring_transaction.RecurringTransaction, error)
	CreateRecurring(userID int64, rt recurring_transaction.RecurringTransaction) error
	DeleteRecurring(id int64) error
	ProcessDueTransactions() error
}

type usecase struct {
	repo   recurring_transaction.Repository
	txRepo transaction.Repository
}

func New(repo recurring_transaction.Repository, txRepo transaction.Repository) Usecase {
	return &usecase{repo: repo, txRepo: txRepo}
}

func (u *usecase) GetRecurring(userID int64) ([]recurring_transaction.RecurringTransaction, error) {
	return u.repo.FindAllByUserID(userID)
}

func (u *usecase) CreateRecurring(userID int64, rt recurring_transaction.RecurringTransaction) error {
	rt.UserID = userID
	if rt.StartDate.IsZero() {
		rt.StartDate = time.Now()
	}
	return u.repo.Save(&rt)
}

func (u *usecase) DeleteRecurring(id int64) error {
	return u.repo.Delete(id)
}

func (u *usecase) ProcessDueTransactions() error {
	now := time.Now()
	due, err := u.repo.FindDue(now)
	if err != nil {
		return err
	}

	for _, rt := range due {
		// Create a real transaction
		tx := &transaction.Transaction{
			UserID:     rt.UserID,
			CategoryID: rt.CategoryID,
			Amount:     rt.Amount,
			Note:       rt.Note + " (Auto-generated)",
			Date:       now,
			Type:       rt.Type,
		}

		if err := u.txRepo.Save(tx); err != nil {
			// In production, log error and continue
			continue
		}

		// Update last processed date
		if err := u.repo.UpdateLastProcessed(rt.ID, now); err != nil {
			continue
		}
	}

	return nil
}
