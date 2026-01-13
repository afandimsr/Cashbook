package postgresql

import (
	"database/sql"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/recurring_transaction"
)

type recurringRepo struct {
	db *sql.DB
}

func NewRecurringRepo(db *sql.DB) recurring_transaction.Repository {
	return &recurringRepo{db: db}
}

func (r *recurringRepo) FindAllByUserID(userID int64) ([]recurring_transaction.RecurringTransaction, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, category_id, amount, type, note, frequency, start_date, last_processed FROM recurring_transactions WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *recurringRepo) FindDue(now time.Time) ([]recurring_transaction.RecurringTransaction, error) {
	// Simple logic: if last_processed is nil, it's due if now >= start_date.
	// If not nil, it depends on frequency.
	// This query identifies records needing processed today/due.
	rows, err := r.db.Query(`
		SELECT id, user_id, category_id, amount, type, note, frequency, start_date, last_processed 
		FROM recurring_transactions 
		WHERE (last_processed IS NULL AND start_date <= $1)
		OR (frequency = 'daily' AND last_processed <= $1 - INTERVAL '1 day')
		OR (frequency = 'weekly' AND last_processed <= $1 - INTERVAL '1 week')
		OR (frequency = 'monthly' AND last_processed <= $1 - INTERVAL '1 month')
	`, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *recurringRepo) Save(rt *recurring_transaction.RecurringTransaction) error {
	return r.db.QueryRow(`
		INSERT INTO recurring_transactions(user_id, category_id, amount, type, note, frequency, start_date, last_processed) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id
	`, rt.UserID, rt.CategoryID, rt.Amount, rt.Type, rt.Note, rt.Frequency, rt.StartDate, rt.LastProcessed).Scan(&rt.ID)
}

func (r *recurringRepo) Update(rt *recurring_transaction.RecurringTransaction) error {
	_, err := r.db.Exec(`
		UPDATE recurring_transactions 
		SET category_id = $1, amount = $2, type = $3, note = $4, frequency = $5, start_date = $6
		WHERE id = $7
	`, rt.CategoryID, rt.Amount, rt.Type, rt.Note, rt.Frequency, rt.StartDate, rt.ID)
	return err
}

func (r *recurringRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM recurring_transactions WHERE id = $1", id)
	return err
}

func (r *recurringRepo) UpdateLastProcessed(id int64, lastProcessed time.Time) error {
	_, err := r.db.Exec("UPDATE recurring_transactions SET last_processed = $1 WHERE id = $2", lastProcessed, id)
	return err
}

func (r *recurringRepo) scanRows(rows *sql.Rows) ([]recurring_transaction.RecurringTransaction, error) {
	var rts []recurring_transaction.RecurringTransaction
	for rows.Next() {
		var rt recurring_transaction.RecurringTransaction
		var lastProcessed sql.NullTime
		err := rows.Scan(&rt.ID, &rt.UserID, &rt.CategoryID, &rt.Amount, &rt.Type, &rt.Note, &rt.Frequency, &rt.StartDate, &lastProcessed)
		if err != nil {
			return nil, err
		}
		if lastProcessed.Valid {
			rt.LastProcessed = lastProcessed.Time
		}
		rts = append(rts, rt)
	}
	return rts, nil
}
