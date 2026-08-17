package postgresql_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/afandimsr/cashbook-backend/internal/domain/category"
	repo "github.com/afandimsr/cashbook-backend/internal/infrastructure/persistent/postgresql/repository"
	"github.com/stretchr/testify/assert"
)

func TestCategoryRepo_FindAllByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	r := repo.NewCategoryRepo(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "type", "color", "icon"}).
		AddRow(int64(1), int64(9), "Food", "expense", "#a", "x").
		AddRow(int64(2), int64(9), "Salary", "income", "#b", "y")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, name, type, color, icon FROM categories WHERE user_id")).
		WithArgs(int64(9)).WillReturnRows(rows)

	got, err := r.FindAllByUserID(9)
	assert.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, "Food", got[0].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_FindByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewCategoryRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta("FROM categories WHERE id")).
		WithArgs(int64(3)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "name", "type", "color", "icon"}).
			AddRow(int64(3), int64(1), "Bills", "expense", "#c", "z"))

	got, err := r.FindByID(3)
	assert.NoError(t, err)
	assert.Equal(t, "Bills", got.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_Save(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewCategoryRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO categories")).
		WithArgs(int64(1), "Transport", "expense", "#d", "bus").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(42)))

	c := category.Category{UserID: 1, Name: "Transport", Type: "expense", Color: "#d", Icon: "bus"}
	err := r.Save(&c)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), c.ID) // RETURNING id populated
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewCategoryRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE categories SET")).
		WithArgs("New", "income", "#e", "i", int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := r.Update(&category.Category{ID: 5, Name: "New", Type: "income", Color: "#e", Icon: "i"})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepo_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	r := repo.NewCategoryRepo(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM categories WHERE id")).
		WithArgs(int64(7)).WillReturnResult(sqlmock.NewResult(0, 1))

	assert.NoError(t, r.Delete(7))
	assert.NoError(t, mock.ExpectationsWereMet())
}
