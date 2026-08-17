package report_test

import (
	"errors"
	"testing"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTxRepository struct{ mock.Mock }

func (m *MockTxRepository) FindAllByUserID(userID int64, limit, offset int, f transaction.Filter) ([]transaction.Transaction, error) {
	args := m.Called(userID, limit, offset, f)
	return args.Get(0).([]transaction.Transaction), args.Error(1)
}
func (m *MockTxRepository) GetTotalAndSum(userID int64, f transaction.Filter) (int64, float64, error) {
	args := m.Called(userID, f)
	return args.Get(0).(int64), args.Get(1).(float64), args.Error(2)
}
func (m *MockTxRepository) GetCategorySpending(userID int64, limit, offset int, f transaction.Filter) ([]transaction.ReportTransaction, error) {
	args := m.Called(userID, limit, offset, f)
	return args.Get(0).([]transaction.ReportTransaction), args.Error(1)
}
func (m *MockTxRepository) FindByID(id int64) (transaction.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}
func (m *MockTxRepository) Save(t *transaction.Transaction) error   { return m.Called(t).Error(0) }
func (m *MockTxRepository) Update(t *transaction.Transaction) error { return m.Called(t).Error(0) }
func (m *MockTxRepository) Delete(id int64) error                   { return m.Called(id).Error(0) }

func TestReport_GetCategorySpending_AggregatesExpensesInPeriod(t *testing.T) {
	repo := new(MockTxRepository)
	u := uc.New(repo)

	may := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	june := time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC)

	rows := []transaction.ReportTransaction{
		{CategoryID: 1, CategoryName: "Food", Color: "#a", Amount: 100, Type: "expense", Date: may},
		{CategoryID: 1, CategoryName: "Food", Color: "#a", Amount: 50, Type: "expense", Date: may},
		{CategoryID: 2, CategoryName: "Salary", Amount: 999, Type: "income", Date: may}, // income excluded
		{CategoryID: 3, CategoryName: "Toys", Amount: 30, Type: "expense", Date: june},  // other month excluded
	}
	repo.On("GetCategorySpending", int64(1), 1000, 0, transaction.Filter{}).Return(rows, nil).Once()

	res, err := u.GetCategorySpending(1, 5, 2026)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, int64(1), res[0].CategoryID)
	assert.Equal(t, 150.0, res[0].TotalAmount)
	repo.AssertExpectations(t)
}

func TestReport_GetCategorySpending_RepoError(t *testing.T) {
	repo := new(MockTxRepository)
	u := uc.New(repo)

	repo.On("GetCategorySpending", int64(1), 1000, 0, transaction.Filter{}).
		Return([]transaction.ReportTransaction{}, errors.New("db")).Once()

	_, err := u.GetCategorySpending(1, 5, 2026)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}
