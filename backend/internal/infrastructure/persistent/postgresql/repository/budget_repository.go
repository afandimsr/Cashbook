package postgresql

import (
	"database/sql"
	"errors"

	"github.com/afandimsr/cashbook-backend/internal/domain/budget"
)

type budgetRepo struct {
	db *sql.DB
}

func NewBudgetRepo(db *sql.DB) budget.Repository {
	return &budgetRepo{db: db}
}

func (r *budgetRepo) FindAllByUserID(userID int64, month, year int) ([]budget.Budget, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, category_id, amount, month, year FROM budgets WHERE user_id = $1 AND month = $2 AND year = $3",
		userID, month, year,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []budget.Budget
	for rows.Next() {
		var b budget.Budget
		if err := rows.Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Month, &b.Year); err != nil {
			return nil, err
		}
		budgets = append(budgets, b)
	}
	return budgets, nil
}

func (r *budgetRepo) FindByCategory(userID int64, categoryID int64, month, year int) (budget.Budget, error) {
	var b budget.Budget
	err := r.db.QueryRow(
		"SELECT id, user_id, category_id, amount, month, year FROM budgets WHERE user_id = $1 AND category_id = $2 AND month = $3 AND year = $4",
		userID, categoryID, month, year,
	).Scan(&b.ID, &b.UserID, &b.CategoryID, &b.Amount, &b.Month, &b.Year)
	if err != nil {
		if err == sql.ErrNoRows {
			return b, errors.New("budget not found")
		}
		return b, err
	}
	return b, nil
}

func (r *budgetRepo) Save(b *budget.Budget) error {
	return r.db.QueryRow(
		"INSERT INTO budgets(user_id, category_id, amount, month, year) VALUES($1, $2, $3, $4, $5) RETURNING id",
		b.UserID, b.CategoryID, b.Amount, b.Month, b.Year,
	).Scan(&b.ID)
}

func (r *budgetRepo) Update(b *budget.Budget) error {
	_, err := r.db.Exec(
		"UPDATE budgets SET amount = $1 WHERE id = $2",
		b.Amount, b.ID,
	)
	return err
}

func (r *budgetRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM budgets WHERE id = $1", id)
	return err
}
