package transaction_test

import (
	"testing"
	"time"

	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) FindAllByUserID(userID int64, limit, offset int, filter transaction.Filter) ([]transaction.Transaction, error) {
	args := m.Called(userID, limit, offset, filter)
	return args.Get(0).([]transaction.Transaction), args.Error(1)
}

// Mock GetTotalAndSum
func (m *MockTransactionRepository) GetTotalAndSum(userID int64, filter transaction.Filter) (int64, float64, error) {
	args := m.Called(userID, filter)
	return args.Get(0).(int64), args.Get(1).(float64), args.Error(2)
}

func (m *MockTransactionRepository) GetCategorySpending(userID int64, limit, offset int, filter transaction.Filter) ([]transaction.ReportTransaction, error) {
	args := m.Called(userID, limit, offset, filter)
	return args.Get(0).([]transaction.ReportTransaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByID(id int64) (transaction.Transaction, error) {
	args := m.Called(id)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Save(t *transaction.Transaction) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTransactionRepository) Update(t *transaction.Transaction) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetAllByUserIDWithFilters(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	usecase := uc.New(mockRepo)

	userID := int64(1)
	page := 1
	limit := 10

	t.Run("FilterByCategoryID", func(t *testing.T) {
		filter := transaction.Filter{CategoryID: 5}
		mockRepo.On("FindAllByUserID", userID, limit, 0, filter).Return([]transaction.Transaction{}, nil).Once()
		mockRepo.On("GetTotalAndSum", userID, filter).Return(int64(0), float64(0), nil).Once()

		_, err := usecase.GetAllByUserID(userID, page, limit, filter)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FilterByType", func(t *testing.T) {
		filter := transaction.Filter{Type: "expense"}
		mockRepo.On("FindAllByUserID", userID, limit, 0, filter).Return([]transaction.Transaction{}, nil).Once()
		mockRepo.On("GetTotalAndSum", userID, filter).Return(int64(0), float64(0), nil).Once()

		_, err := usecase.GetAllByUserID(userID, page, limit, filter)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FilterByDateRange", func(t *testing.T) {
		start := time.Now().AddDate(0, -1, 0)
		end := time.Now()
		filter := transaction.Filter{StartDate: start, EndDate: end}
		mockRepo.On("FindAllByUserID", userID, limit, 0, filter).Return([]transaction.Transaction{}, nil).Once()
		mockRepo.On("GetTotalAndSum", userID, filter).Return(int64(0), float64(0), nil).Once()

		_, err := usecase.GetAllByUserID(userID, page, limit, filter)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FilterByAll", func(t *testing.T) {
		filter := transaction.Filter{
			Search:     "lunch",
			CategoryID: 2,
			Type:       "expense",
			StartDate:  time.Now(),
		}
		mockRepo.On("FindAllByUserID", userID, limit, 0, filter).Return([]transaction.Transaction{}, nil).Once()
		mockRepo.On("GetTotalAndSum", userID, filter).Return(int64(0), float64(0), nil).Once()

		_, err := usecase.GetAllByUserID(userID, page, limit, filter)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	usecase := uc.New(mockRepo)
	id := int64(1)

	expected := transaction.Transaction{ID: id, Note: "Test"}
	mockRepo.On("FindByID", id).Return(expected, nil).Once()

	result, err := usecase.GetByID(id)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestCreate(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	usecase := uc.New(mockRepo)

	tx := transaction.Transaction{Note: "Test"}
	mockRepo.On("Save", mock.MatchedBy(func(t *transaction.Transaction) bool {
		return t.Note == "Test"
	})).Return(nil).Once()

	err := usecase.Create(tx)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdate(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	usecase := uc.New(mockRepo)
	id := int64(1)

	existing := transaction.Transaction{ID: id, Note: "Old"}
	updateData := transaction.Transaction{Note: "New"}

	mockRepo.On("FindByID", id).Return(existing, nil).Once()
	mockRepo.On("Update", mock.MatchedBy(func(t *transaction.Transaction) bool {
		return t.Note == "New"
	})).Return(nil).Once()

	err := usecase.Update(id, updateData)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	usecase := uc.New(mockRepo)
	id := int64(1)

	mockRepo.On("Delete", id).Return(nil).Once()

	err := usecase.Delete(id)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetDashboardSummary(t *testing.T) {
	mockRepo := new(MockTransactionRepository)
	usecase := uc.New(mockRepo)
	userID := int64(1)

	txs := []transaction.Transaction{
		{Type: "income", Amount: 100},
		{Type: "expense", Amount: 40},
	}

	mockRepo.On("FindAllByUserID", userID, 1000, 0, transaction.Filter{}).Return(txs, nil).Once()

	summary, err := usecase.GetDashboardSummary(userID)

	assert.NoError(t, err)
	assert.Equal(t, 100.0, summary.TotalIncome)
	assert.Equal(t, 40.0, summary.TotalExpense)
	assert.Equal(t, 60.0, summary.Balance)
	mockRepo.AssertExpectations(t)
}
