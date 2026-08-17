package postgresql_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	repo "github.com/afandimsr/cashbook-backend/internal/infrastructure/persistent/postgresql/repository"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepo_Save(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTransactionRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO transactions")).
		WithArgs(int64(1), int64(2), 500.0, "Lunch", now, "expense").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(99)))

	tx := &transaction.Transaction{UserID: 1, CategoryID: 2, Amount: 500, Note: "Lunch", Date: now, Type: "expense"}
	err := r.Save(tx)
	assert.NoError(t, err)
	assert.Equal(t, int64(99), tx.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepo_FindByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTransactionRepo(db)

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("FROM transactions WHERE id")).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "category_id", "amount", "note", "date", "type"}).
			AddRow(int64(99), int64(1), int64(2), 500.0, "Lunch", now, "expense"))

	got, err := r.FindByID(99)
	assert.NoError(t, err)
	assert.Equal(t, "Lunch", got.Note)
	assert.Equal(t, 500.0, got.Amount)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTransactionRepo_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewTransactionRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM transactions WHERE id")).
		WithArgs(int64(5)).WillReturnResult(sqlmock.NewResult(0, 1))

	assert.NoError(t, r.Delete(5))
	assert.NoError(t, mock.ExpectationsWereMet())
}
