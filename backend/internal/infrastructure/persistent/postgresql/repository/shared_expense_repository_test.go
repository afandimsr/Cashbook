package postgresql_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	se "github.com/afandimsr/cashbook-backend/internal/domain/shared_expense"
	repo "github.com/afandimsr/cashbook-backend/internal/infrastructure/persistent/postgresql/repository"
	"github.com/stretchr/testify/assert"
)

func p64(v int64) *int64 { return &v }

func TestSharedExpenseRepo_Save(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewSharedExpenseRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO split_bills")).
		WithArgs(int64(1), int64(1), int64(5), "Lunch", 100000.0, 0.0, 0.0, 0.0, 0.0, 100000.0, now, "OPEN", "EQUAL").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow("bill-1", now, now))

	b := se.SplitBill{
		CreatorID: 1, PayerID: 1, CategoryID: 5, Title: "Lunch",
		Subtotal: 100000, TotalAmount: 100000, Date: now,
		Status: se.StatusOpen, SplitMethod: se.MethodEqual,
	}
	err := r.Save(context.Background(), &b)
	assert.NoError(t, err)
	assert.Equal(t, "bill-1", b.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSharedExpenseRepo_SaveParticipant(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewSharedExpenseRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO split_participants")).
		WithArgs("bill-1", int64(2), "", 50000.0, 50000.0, false).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow("p-1", now))

	part := se.SplitParticipant{SplitBillID: "bill-1", UserID: p64(2), BaseShare: 50000, ShareAmount: 50000}
	err := r.SaveParticipant(context.Background(), &part)
	assert.NoError(t, err)
	assert.Equal(t, "p-1", part.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSharedExpenseRepo_SaveItemAndLink(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewSharedExpenseRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO split_items")).
		WithArgs("bill-1", "Nasi", 40000.0, 1, p64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow("item-1", now))

	item := se.SplitItem{SplitBillID: "bill-1", Name: "Nasi", Price: 40000, Quantity: 1, CategoryID: p64(10)}
	assert.NoError(t, r.SaveItem(context.Background(), &item))
	assert.Equal(t, "item-1", item.ID)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO split_item_participants")).
		WithArgs("item-1", "p-1").WillReturnResult(sqlmock.NewResult(0, 1))
	assert.NoError(t, r.AddItemParticipant(context.Background(), "item-1", "p-1"))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSharedExpenseRepo_GetSummary(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewSharedExpenseRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta("b.payer_id = $1 AND p.is_paid = FALSE")).
		WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(100.0))
	mock.ExpectQuery(regexp.QuoteMeta("p.user_id = $1 AND p.is_paid = FALSE")).
		WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(20.0))

	sum, err := r.GetSummary(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, 100.0, sum.TotalOwedToMe)
	assert.Equal(t, 20.0, sum.TotalIOwe)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSharedExpenseRepo_FindByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewSharedExpenseRepo(db)

	now := time.Now()
	// 1. bill
	mock.ExpectQuery(regexp.QuoteMeta("FROM split_bills WHERE id")).
		WithArgs("bill-1").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "creator_id", "payer_id", "category_id", "title",
			"subtotal", "tax_amount", "service_charge", "other_charge", "discount_amount", "total_amount",
			"date", "status", "split_method", "created_at", "updated_at",
		}).AddRow("bill-1", int64(1), int64(1), int64(5), "Lunch",
			100000.0, 0.0, 0.0, 0.0, 0.0, 100000.0,
			now, "OPEN", "EQUAL", now, now))

	// 2. participants
	mock.ExpectQuery(regexp.QuoteMeta("FROM split_participants p")).
		WithArgs("bill-1").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "split_bill_id", "user_id", "shadow_name", "name",
			"base_share", "share_amount", "is_paid", "settled_at", "transaction_id", "created_at",
		}).AddRow("p-1", "bill-1", int64(2), "", "Bob", 50000.0, 50000.0, false, nil, nil, now))

	// 3. items (empty → link query is skipped)
	mock.ExpectQuery(regexp.QuoteMeta("FROM split_items WHERE split_bill_id")).
		WithArgs("bill-1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "split_bill_id", "name", "price", "quantity", "category_id", "created_at"}))

	bill, err := r.FindByID(context.Background(), "bill-1")
	assert.NoError(t, err)
	assert.Equal(t, "Lunch", bill.Title)
	assert.Len(t, bill.Participants, 1)
	assert.Equal(t, "Bob", bill.Participants[0].UserName)
	assert.Empty(t, bill.Items)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSharedExpenseRepo_WithTransaction_Commit(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewSharedExpenseRepo(db)

	mock.ExpectBegin()
	mock.ExpectCommit()

	err := r.WithTransaction(context.Background(), func(txRepo se.Repository) error { return nil })
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSharedExpenseRepo_WithTransaction_Rollback(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewSharedExpenseRepo(db)

	mock.ExpectBegin()
	mock.ExpectRollback()

	wantErr := errors.New("boom")
	err := r.WithTransaction(context.Background(), func(txRepo se.Repository) error { return wantErr })
	assert.ErrorIs(t, err, wantErr)
	assert.NoError(t, mock.ExpectationsWereMet())
}
