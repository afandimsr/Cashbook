package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/shared_expense"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
)

type sharedExpenseQueryExecutor interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type sharedExpenseRepo struct {
	db sharedExpenseQueryExecutor
}

func NewSharedExpenseRepo(db *sql.DB) shared_expense.Repository {
	return &sharedExpenseRepo{db: db}
}

func (r *sharedExpenseRepo) FindAllByUserID(ctx context.Context, userID int64, filter shared_expense.Filter) ([]shared_expense.SplitBill, error) {
	query := `
		SELECT DISTINCT b.id, b.creator_id, b.payer_id, b.category_id, b.title,
		       b.subtotal, b.tax_amount, b.service_charge, b.other_charge, b.discount_amount, b.total_amount,
		       b.date, b.status, b.split_method, b.created_at, b.updated_at
		FROM split_bills b
		LEFT JOIN split_participants p ON b.id = p.split_bill_id
		WHERE b.payer_id = $1 OR p.user_id = $1
		ORDER BY b.date DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bills []shared_expense.SplitBill
	for rows.Next() {
		var b shared_expense.SplitBill
		err := rows.Scan(&b.ID, &b.CreatorID, &b.PayerID, &b.CategoryID, &b.Title,
			&b.Subtotal, &b.TaxAmount, &b.ServiceCharge, &b.OtherCharge, &b.DiscountAmount, &b.TotalAmount,
			&b.Date, &b.Status, &b.SplitMethod, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bills = append(bills, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load participants and items for each bill so the UI can render the breakdown.
	for i := range bills {
		participants, err := r.FindParticipantsByBillID(ctx, bills[i].ID)
		if err != nil {
			return nil, err
		}
		bills[i].Participants = participants

		items, err := r.FindItemsByBillID(ctx, bills[i].ID)
		if err != nil {
			return nil, err
		}
		bills[i].Items = items
	}

	return bills, nil
}

func (r *sharedExpenseRepo) FindByID(ctx context.Context, id string) (shared_expense.SplitBill, error) {
	var b shared_expense.SplitBill
	query := `
		SELECT id, creator_id, payer_id, category_id, title,
		       subtotal, tax_amount, service_charge, other_charge, discount_amount, total_amount,
		       date, status, split_method, created_at, updated_at
		FROM split_bills WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(&b.ID, &b.CreatorID, &b.PayerID, &b.CategoryID, &b.Title,
		&b.Subtotal, &b.TaxAmount, &b.ServiceCharge, &b.OtherCharge, &b.DiscountAmount, &b.TotalAmount,
		&b.Date, &b.Status, &b.SplitMethod, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return b, err
	}

	participants, err := r.FindParticipantsByBillID(ctx, id)
	if err != nil {
		return b, err
	}
	b.Participants = participants

	items, err := r.FindItemsByBillID(ctx, id)
	if err != nil {
		return b, err
	}
	b.Items = items

	return b, nil
}

func (r *sharedExpenseRepo) Save(ctx context.Context, b *shared_expense.SplitBill) error {
	query := `
		INSERT INTO split_bills (creator_id, payer_id, category_id, title, subtotal, tax_amount, service_charge, other_charge, discount_amount, total_amount, date, status, split_method)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, b.CreatorID, b.PayerID, b.CategoryID, b.Title, b.Subtotal, b.TaxAmount, b.ServiceCharge, b.OtherCharge, b.DiscountAmount, b.TotalAmount, b.Date, b.Status, b.SplitMethod).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *sharedExpenseRepo) UpdateStatus(ctx context.Context, id string, status shared_expense.SplitBillStatus) error {
	query := `UPDATE split_bills SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *sharedExpenseRepo) SaveParticipant(ctx context.Context, p *shared_expense.SplitParticipant) error {
	query := `
		INSERT INTO split_participants (split_bill_id, user_id, shadow_name, base_share, share_amount, is_paid)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, p.SplitBillID, p.UserID, p.ShadowName, p.BaseShare, p.ShareAmount, p.IsPaid).Scan(&p.ID, &p.CreatedAt)
}

func (r *sharedExpenseRepo) UpdateParticipantSettlement(ctx context.Context, id string, isPaid bool, settledAt *time.Time, transactionID *int64) error {
	query := `UPDATE split_participants SET is_paid = $1, settled_at = $2, transaction_id = $3 WHERE id = $4`
	_, err := r.db.Exec(query, isPaid, settledAt, transactionID, id)
	return err
}

func (r *sharedExpenseRepo) FindParticipantsByBillID(ctx context.Context, billID string) ([]shared_expense.SplitParticipant, error) {
	query := `
		SELECT p.id, p.split_bill_id, p.user_id, p.shadow_name, COALESCE(u.name, ''),
		       p.base_share, p.share_amount, p.is_paid, p.settled_at, p.transaction_id, p.created_at
		FROM split_participants p
		LEFT JOIN users u ON p.user_id = u.id
		WHERE p.split_bill_id = $1
		ORDER BY p.created_at
	`
	rows, err := r.db.Query(query, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []shared_expense.SplitParticipant
	for rows.Next() {
		var p shared_expense.SplitParticipant
		err := rows.Scan(&p.ID, &p.SplitBillID, &p.UserID, &p.ShadowName, &p.UserName, &p.BaseShare, &p.ShareAmount, &p.IsPaid, &p.SettledAt, &p.TransactionID, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *sharedExpenseRepo) SaveItem(ctx context.Context, item *shared_expense.SplitItem) error {
	query := `
		INSERT INTO split_items (split_bill_id, name, price, quantity, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, item.SplitBillID, item.Name, item.Price, item.Quantity, item.CategoryID).Scan(&item.ID, &item.CreatedAt)
}

func (r *sharedExpenseRepo) AddItemParticipant(ctx context.Context, itemID, participantID string) error {
	query := `
		INSERT INTO split_item_participants (item_id, participant_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`
	_, err := r.db.Exec(query, itemID, participantID)
	return err
}

func (r *sharedExpenseRepo) FindItemsByBillID(ctx context.Context, billID string) ([]shared_expense.SplitItem, error) {
	query := `
		SELECT id, split_bill_id, name, price, quantity, category_id, created_at
		FROM split_items WHERE split_bill_id = $1
		ORDER BY created_at
	`
	rows, err := r.db.Query(query, billID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []shared_expense.SplitItem
	byID := make(map[string]*shared_expense.SplitItem)
	for rows.Next() {
		var it shared_expense.SplitItem
		if err := rows.Scan(&it.ID, &it.SplitBillID, &it.Name, &it.Price, &it.Quantity, &it.CategoryID, &it.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return items, nil
	}
	for i := range items {
		byID[items[i].ID] = &items[i]
	}

	// Load the participant assignments for this bill's items in one query.
	linkQuery := `
		SELECT ip.item_id, ip.participant_id
		FROM split_item_participants ip
		JOIN split_items i ON ip.item_id = i.id
		WHERE i.split_bill_id = $1
	`
	linkRows, err := r.db.Query(linkQuery, billID)
	if err != nil {
		return nil, err
	}
	defer linkRows.Close()

	for linkRows.Next() {
		var itemID, participantID string
		if err := linkRows.Scan(&itemID, &participantID); err != nil {
			return nil, err
		}
		if it, ok := byID[itemID]; ok {
			it.ParticipantIDs = append(it.ParticipantIDs, participantID)
		}
	}
	if err := linkRows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// SaveTransaction records a standard income/expense transaction using the same
// underlying executor (DB or Tx) as this repository. When called from a txRepo
// inside WithTransaction, the write participates in the same DB transaction.
func (r *sharedExpenseRepo) SaveTransaction(ctx context.Context, t *transaction.Transaction) error {
	return (&transactionRepo{db: r.db}).Save(t)
}

func (r *sharedExpenseRepo) GetSummary(ctx context.Context, userID int64) (shared_expense.SplitSummary, error) {
	var summary shared_expense.SplitSummary

	// Total Owed To Me: I am the payer, others owe me
	owedToMeQuery := `
		SELECT COALESCE(SUM(share_amount), 0)
		FROM split_participants p
		JOIN split_bills b ON p.split_bill_id = b.id
		WHERE b.payer_id = $1 AND p.is_paid = FALSE AND (p.user_id != $1 OR p.user_id IS NULL)
	`
	err := r.db.QueryRow(owedToMeQuery, userID).Scan(&summary.TotalOwedToMe)
	if err != nil {
		return summary, err
	}

	// Total I Owe: I am a participant, I owe someone else
	iOweQuery := `
		SELECT COALESCE(SUM(p.share_amount), 0)
		FROM split_participants p
		JOIN split_bills b ON p.split_bill_id = b.id
		WHERE p.user_id = $1 AND p.is_paid = FALSE AND b.payer_id != $1
	`
	err = r.db.QueryRow(iOweQuery, userID).Scan(&summary.TotalIOwe)
	if err != nil {
		return summary, err
	}

	return summary, nil
}

func (r *sharedExpenseRepo) WithTransaction(ctx context.Context, fn func(repo shared_expense.Repository) error) error {
	db, ok := r.db.(*sql.DB)
	if !ok {
		return fmt.Errorf("repository must be initialized with *sql.DB to start a transaction")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	txRepo := &sharedExpenseRepo{db: tx}
	err = fn(txRepo)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
