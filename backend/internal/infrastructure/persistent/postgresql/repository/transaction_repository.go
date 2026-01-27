package postgresql

import (
	"database/sql"
	"strconv"

	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

type transactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) transaction.Repository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) FindAllByUserID(userID int64, limit, offset int, filter transaction.Filter) ([]transaction.Transaction, error) {
	query := "SELECT id, user_id, category_id, amount, note, date, type FROM transactions WHERE user_id = $1"
	args := []interface{}{userID}
	placeholderCount := 1

	if filter.Search != "" {
		placeholderCount++
		query += " AND note ILIKE $" + strconv.Itoa(placeholderCount)
		args = append(args, "%"+filter.Search+"%")
	}

	if filter.CategoryID != 0 {
		placeholderCount++
		query += " AND category_id = $" + strconv.Itoa(placeholderCount)
		args = append(args, filter.CategoryID)
	}

	if filter.Type != "" {
		placeholderCount++
		query += " AND type = $" + strconv.Itoa(placeholderCount)
		args = append(args, filter.Type)
	}

	if !filter.StartDate.IsZero() {
		placeholderCount++
		query += " AND date >= $" + strconv.Itoa(placeholderCount)
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		placeholderCount++
		query += " AND date <= $" + strconv.Itoa(placeholderCount)
		args = append(args, filter.EndDate)
	}

	query += " ORDER BY date DESC LIMIT $" + strconv.Itoa(placeholderCount+1) + " OFFSET $" + strconv.Itoa(placeholderCount+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []transaction.Transaction
	for rows.Next() {
		var t transaction.Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.CategoryID, &t.Amount, &t.Note, &t.Date, &t.Type)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *transactionRepo) GetCategorySpending(userID int64, limit, offset int, filter transaction.Filter) ([]transaction.ReportTransaction, error) {
	query := "SELECT t.id, t.user_id, t.category_id, c.name as category_name, c.color as color, t.amount, t.note, t.date, t.type FROM transactions as t INNER JOIN categories c ON t.category_id = c.id WHERE t.user_id = $1"
	args := []interface{}{userID}
	placeholderCount := 1

	if filter.Search != "" {
		placeholderCount++
		query += " AND note ILIKE $" + strconv.Itoa(placeholderCount)
		args = append(args, "%"+filter.Search+"%")
	}

	if filter.CategoryID != 0 {
		placeholderCount++
		query += " AND t.category_id = $" + strconv.Itoa(placeholderCount)
		args = append(args, filter.CategoryID)
	}

	if filter.Type != "" {
		placeholderCount++
		query += " AND t.type = $" + strconv.Itoa(placeholderCount)
		args = append(args, filter.Type)
	}

	if !filter.StartDate.IsZero() {
		placeholderCount++
		query += " AND t.date >= $" + strconv.Itoa(placeholderCount)
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		placeholderCount++
		query += " AND t.date <= $" + strconv.Itoa(placeholderCount)
		args = append(args, filter.EndDate)
	}

	query += " ORDER BY date DESC LIMIT $" + strconv.Itoa(placeholderCount+1) + " OFFSET $" + strconv.Itoa(placeholderCount+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []transaction.ReportTransaction
	for rows.Next() {
		var t transaction.ReportTransaction
		err := rows.Scan(&t.ID, &t.UserID, &t.CategoryID, &t.CategoryName, &t.Color, &t.Amount, &t.Note, &t.Date, &t.Type)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *transactionRepo) FindByID(id int64) (transaction.Transaction, error) {
	var t transaction.Transaction
	err := r.db.QueryRow(
		"SELECT id, user_id, category_id, amount, note, date, type FROM transactions WHERE id = $1",
		id,
	).Scan(&t.ID, &t.UserID, &t.CategoryID, &t.Amount, &t.Note, &t.Date, &t.Type)
	return t, err
}

func (r *transactionRepo) Save(t *transaction.Transaction) error {
	return r.db.QueryRow(
		"INSERT INTO transactions(user_id, category_id, amount, note, date, type) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		t.UserID, t.CategoryID, t.Amount, t.Note, t.Date, t.Type,
	).Scan(&t.ID)
}

func (r *transactionRepo) Update(t *transaction.Transaction) error {
	_, err := r.db.Exec(
		"UPDATE transactions SET category_id = $1, amount = $2, note = $3, date = $4, type = $5 WHERE id = $6",
		t.CategoryID, t.Amount, t.Note, t.Date, t.Type, t.ID,
	)
	return err
}

func (r *transactionRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM transactions WHERE id = $1", id)
	return err
}
