package transaction

import (
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

type Usecase interface {
	GetAllByUserID(userID int64, page, limit int, search string) ([]transaction.Transaction, error)
	GetByID(id int64) (transaction.Transaction, error)
	Create(t transaction.Transaction) error
	Update(id int64, t transaction.Transaction) error
	Delete(id int64) error
	GetDashboardSummary(userID int64) (transaction.DashboardSummary, error)
}

type usecase struct {
	repo transaction.Repository
}

func New(repo transaction.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) GetAllByUserID(userID int64, page, limit int, search string) ([]transaction.Transaction, error) {
	offset := (page - 1) * limit
	return u.repo.FindAllByUserID(userID, limit, offset, search)
}

func (u *usecase) GetByID(id int64) (transaction.Transaction, error) {
	return u.repo.FindByID(id)
}

func (u *usecase) Create(t transaction.Transaction) error {
	if t.Date.IsZero() {
		t.Date = time.Now()
	}
	return u.repo.Save(&t)
}

func (u *usecase) Update(id int64, t transaction.Transaction) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	existing.CategoryID = t.CategoryID
	existing.Amount = t.Amount
	existing.Note = t.Note
	existing.Date = t.Date
	existing.Type = t.Type

	return u.repo.Update(&existing)
}

func (u *usecase) Delete(id int64) error {
	return u.repo.Delete(id)
}

func (u *usecase) GetDashboardSummary(userID int64) (transaction.DashboardSummary, error) {
	// For simplicity, fetch all (or a large range) and calculate
	// In production, this should be a DB aggregation query
	txs, err := u.repo.FindAllByUserID(userID, 1000, 0, "")
	if err != nil {
		return transaction.DashboardSummary{}, err
	}

	var summary transaction.DashboardSummary
	for _, tx := range txs {
		if tx.Type == "income" {
			summary.TotalIncome += tx.Amount
		} else {
			summary.TotalExpense += tx.Amount
		}
	}
	summary.Balance = summary.TotalIncome - summary.TotalExpense
	return summary, nil
}
