package report

import (
	"sort"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

// reportRowLimit caps how many transactions a report aggregates in memory.
const reportRowLimit = 1000

type CategoryReport struct {
	CategoryID   int64   `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
	Color        string  `json:"color"`
}

// MonthlyCategoryReport is the per-category income and expense breakdown for
// one month. Both slices are always non-nil (serialized as []) and sorted by
// TotalAmount descending.
type MonthlyCategoryReport struct {
	Income  []CategoryReport `json:"income"`
	Expense []CategoryReport `json:"expense"`
}

type Usecase interface {
	GetCategorySpending(userID int64, month, year int) ([]CategoryReport, error)
	GetMonthlyCategoryReport(userID int64, month, year int) (MonthlyCategoryReport, error)
}

type usecase struct {
	txRepo transaction.Repository
}

func New(txRepo transaction.Repository) Usecase {
	return &usecase{txRepo: txRepo}
}

func (u *usecase) GetCategorySpending(userID int64, month, year int) ([]CategoryReport, error) {
	// For simplicity, aggregate in memory. Optimized with DB queries in production.
	txs, err := u.txRepo.GetCategorySpending(userID, reportRowLimit, 0, transaction.Filter{})
	if err != nil {
		return nil, err
	}
	return aggregateByCategory(txs, "expense", month, year), nil
}

// GetMonthlyCategoryReport filters by date range in the query (unlike
// GetCategorySpending), so the row limit applies to the requested month only
// and older months are never silently truncated by newer activity.
func (u *usecase) GetMonthlyCategoryReport(userID int64, month, year int) (MonthlyCategoryReport, error) {
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	txs, err := u.txRepo.GetCategorySpending(userID, reportRowLimit, 0, transaction.Filter{StartDate: start, EndDate: end})
	if err != nil {
		return MonthlyCategoryReport{}, err
	}

	return MonthlyCategoryReport{
		Income:  sortedByAmount(aggregateByCategory(txs, "income", month, year)),
		Expense: sortedByAmount(aggregateByCategory(txs, "expense", month, year)),
	}, nil
}

func aggregateByCategory(txs []transaction.ReportTransaction, txType string, month, year int) []CategoryReport {
	totals := make(map[int64]*CategoryReport)

	for _, tx := range txs {
		if tx.Type == txType && int(tx.Date.Month()) == month && tx.Date.Year() == year {
			if _, ok := totals[tx.CategoryID]; !ok {
				totals[tx.CategoryID] = &CategoryReport{
					CategoryID:   tx.CategoryID,
					CategoryName: tx.CategoryName,
					Color:        tx.Color,
					TotalAmount:  0,
				}
			}
			totals[tx.CategoryID].TotalAmount += tx.Amount
		}
	}

	var res []CategoryReport
	for _, v := range totals {
		res = append(res, *v)
	}

	return res
}

func sortedByAmount(reports []CategoryReport) []CategoryReport {
	if reports == nil {
		return []CategoryReport{}
	}
	sort.Slice(reports, func(i, j int) bool { return reports[i].TotalAmount > reports[j].TotalAmount })
	return reports
}
