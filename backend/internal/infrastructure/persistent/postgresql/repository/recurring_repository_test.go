package postgresql_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	rt "github.com/afandimsr/cashbook-backend/internal/domain/recurring_transaction"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	repo "github.com/afandimsr/cashbook-backend/internal/infrastructure/persistent/postgresql/repository"
	"github.com/stretchr/testify/assert"
)

func TestRecurringRepo_FindDueByUserID_ScopesToUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewRecurringRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("WHERE user_id = $2 AND")).
		WithArgs(now, int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "category_id", "amount", "type", "note", "frequency", "start_date", "last_processed"}).
			AddRow(int64(1), int64(7), int64(3), 50.0, "expense", "Netflix", "monthly", now, nil))

	got, err := r.FindDueByUserID(7, now)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	assert.True(t, got[0].LastProcessed.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecurringRepo_ProcessOne_ClaimsAndInserts(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewRecurringRepo(db)

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE recurring_transactions SET last_processed = $1 WHERE id = $2 AND")).
		WithArgs(now, int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO transactions")).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(99)))
	mock.ExpectCommit()

	tx := &transaction.Transaction{UserID: 7, CategoryID: 3, Amount: 50, Date: now, Type: "expense"}
	err := r.ProcessOne(context.Background(), 1, tx, now)
	assert.NoError(t, err)
	assert.Equal(t, int64(99), tx.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecurringRepo_ProcessOne_AlreadyProcessedRollsBack(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewRecurringRepo(db)

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE recurring_transactions")).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	err := r.ProcessOne(context.Background(), 1, &transaction.Transaction{}, now)
	assert.ErrorIs(t, err, rt.ErrAlreadyProcessed)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecurringRepo_ProcessOne_InsertFailureRollsBack(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewRecurringRepo(db)

	now := time.Now()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE recurring_transactions")).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO transactions")).
		WillReturnError(errors.New("insert failed"))
	mock.ExpectRollback()

	err := r.ProcessOne(context.Background(), 1, &transaction.Transaction{}, now)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
