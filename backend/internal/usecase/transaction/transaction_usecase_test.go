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

		_, err := usecase.GetAllByUserID(userID, page, limit, filter)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FilterByType", func(t *testing.T) {
		filter := transaction.Filter{Type: "expense"}
		mockRepo.On("FindAllByUserID", userID, limit, 0, filter).Return([]transaction.Transaction{}, nil).Once()

		_, err := usecase.GetAllByUserID(userID, page, limit, filter)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FilterByDateRange", func(t *testing.T) {
		start := time.Now().AddDate(0, -1, 0)
		end := time.Now()
		filter := transaction.Filter{StartDate: start, EndDate: end}
		mockRepo.On("FindAllByUserID", userID, limit, 0, filter).Return([]transaction.Transaction{}, nil).Once()

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

		_, err := usecase.GetAllByUserID(userID, page, limit, filter)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
