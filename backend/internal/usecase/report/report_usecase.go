package report

import (
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

type CategoryReport struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
	Color        string  `json:"color"`
}

type Usecase interface {
	GetCategorySpending(userID int64, month, year int) ([]CategoryReport, error)
}

type usecase struct {
	txRepo transaction.Repository
}

func New(txRepo transaction.Repository) Usecase {
	return &usecase{txRepo: txRepo}
}

func (u *usecase) GetCategorySpending(userID int64, month, year int) ([]CategoryReport, error) {
	// For simplicity, aggregate in memory. Optimized with DB queries in production.
	txs, err := u.txRepo.GetCategorySpending(userID, 1000, 0, "")
	if err != nil {
		return nil, err
	}

	spending := make(map[int64]*CategoryReport)

	for _, tx := range txs {
		if tx.Type == "expense" && int(tx.Date.Month()) == month && tx.Date.Year() == year {
			if _, ok := spending[tx.CategoryID]; !ok {
				spending[tx.CategoryID] = &CategoryReport{
					CategoryID:   tx.CategoryID,
					CategoryName: tx.CategoryName,
					Color:        tx.Color,
					TotalAmount:  0,
				}
			}
			spending[tx.CategoryID].TotalAmount += tx.Amount
		}
	}

	var res []CategoryReport
	for _, v := range spending {
		res = append(res, *v)
	}

	return res, nil
}
