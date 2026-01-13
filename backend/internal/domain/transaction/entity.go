package transaction

import "time"

type Transaction struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	CategoryID int64     `json:"category_id"`
	Amount     float64   `json:"amount"`
	Note       string    `json:"note"`
	Date       time.Time `json:"date"`
	Type       string    `json:"type"` // "income" or "expense"
}

type ReportTransaction struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	CategoryID   int64     `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Color        string    `json:"color"`
	Amount       float64   `json:"amount"`
	Note         string    `json:"note"`
	Date         time.Time `json:"date"`
	Type         string    `json:"type"` // "income" or "expense"
}

type DashboardSummary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

type Repository interface {
	FindAllByUserID(userID int64, limit, offset int, search string) ([]Transaction, error)
	GetCategorySpending(userID int64, limit, offset int, search string) ([]ReportTransaction, error)
	FindByID(id int64) (Transaction, error)
	Save(transaction *Transaction) error
	Update(transaction *Transaction) error
	Delete(id int64) error
}
