package recurring_transaction_test

import (
	"errors"
	"testing"
	"time"

	rt "github.com/afandimsr/cashbook-backend/internal/domain/recurring_transaction"
	"github.com/afandimsr/cashbook-backend/internal/domain/transaction"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/recurring_transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRecurringRepository struct {
	mock.Mock
}

func (m *MockRecurringRepository) FindAllByUserID(userID int64) ([]rt.RecurringTransaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]rt.RecurringTransaction), args.Error(1)
}

func (m *MockRecurringRepository) FindDue(now time.Time) ([]rt.RecurringTransaction, error) {
	args := m.Called(now)
	return args.Get(0).([]rt.RecurringTransaction), args.Error(1)
}

func (m *MockRecurringRepository) Save(r *rt.RecurringTransaction) error   { return m.Called(r).Error(0) }
func (m *MockRecurringRepository) Update(r *rt.RecurringTransaction) error { return m.Called(r).Error(0) }
func (m *MockRecurringRepository) Delete(id int64) error                   { return m.Called(id).Error(0) }
func (m *MockRecurringRepository) UpdateLastProcessed(id int64, t time.Time) error {
	return m.Called(id, t).Error(0)
}

// Minimal transaction.Repository mock (only Save is exercised here).
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

func TestRecurring_Create_SetsUserAndDefaultStart(t *testing.T) {
	repo := new(MockRecurringRepository)
	txRepo := new(MockTxRepository)
	u := uc.New(repo, txRepo)

	repo.On("Save", mock.MatchedBy(func(r *rt.RecurringTransaction) bool {
		return r.UserID == 11 && !r.StartDate.IsZero()
	})).Return(nil).Once()

	got, err := u.CreateRecurring(11, rt.RecurringTransaction{Amount: 100, Frequency: rt.Monthly})
	assert.NoError(t, err)
	assert.Equal(t, int64(11), got.UserID)
	assert.False(t, got.StartDate.IsZero())
	repo.AssertExpectations(t)
}

func TestRecurring_Create_KeepsProvidedStart(t *testing.T) {
	repo := new(MockRecurringRepository)
	txRepo := new(MockTxRepository)
	u := uc.New(repo, txRepo)

	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repo.On("Save", mock.MatchedBy(func(r *rt.RecurringTransaction) bool {
		return r.StartDate.Equal(start)
	})).Return(nil).Once()

	got, err := u.CreateRecurring(1, rt.RecurringTransaction{StartDate: start})
	assert.NoError(t, err)
	assert.Equal(t, start, got.StartDate)
	repo.AssertExpectations(t)
}

func TestRecurring_Create_SaveError(t *testing.T) {
	repo := new(MockRecurringRepository)
	txRepo := new(MockTxRepository)
	u := uc.New(repo, txRepo)

	repo.On("Save", mock.Anything).Return(errors.New("db")).Once()

	_, err := u.CreateRecurring(1, rt.RecurringTransaction{})
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestRecurring_ProcessDue_CreatesTransactions(t *testing.T) {
	repo := new(MockRecurringRepository)
	txRepo := new(MockTxRepository)
	u := uc.New(repo, txRepo)

	due := []rt.RecurringTransaction{
		{ID: 1, UserID: 2, CategoryID: 3, Amount: 50, Type: "expense", Note: "Netflix"},
		{ID: 2, UserID: 2, CategoryID: 4, Amount: 20, Type: "expense", Note: "Spotify"},
	}
	repo.On("FindDue", mock.Anything).Return(due, nil).Once()
	txRepo.On("Save", mock.MatchedBy(func(t *transaction.Transaction) bool {
		return t.Note == "Netflix (Auto-generated)"
	})).Return(nil).Once()
	txRepo.On("Save", mock.MatchedBy(func(t *transaction.Transaction) bool {
		return t.Note == "Spotify (Auto-generated)"
	})).Return(nil).Once()
	repo.On("UpdateLastProcessed", int64(1), mock.Anything).Return(nil).Once()
	repo.On("UpdateLastProcessed", int64(2), mock.Anything).Return(nil).Once()

	err := u.ProcessDueTransactions()
	assert.NoError(t, err)
	repo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestRecurring_ProcessDue_SkipsUpdateOnSaveError(t *testing.T) {
	repo := new(MockRecurringRepository)
	txRepo := new(MockTxRepository)
	u := uc.New(repo, txRepo)

	due := []rt.RecurringTransaction{{ID: 1, UserID: 2, Note: "x"}}
	repo.On("FindDue", mock.Anything).Return(due, nil).Once()
	txRepo.On("Save", mock.Anything).Return(errors.New("db")).Once()

	err := u.ProcessDueTransactions()
	assert.NoError(t, err) // ProcessDue swallows per-item errors and continues
	repo.AssertNotCalled(t, "UpdateLastProcessed", mock.Anything, mock.Anything)
	repo.AssertExpectations(t)
	txRepo.AssertExpectations(t)
}

func TestRecurring_ProcessDue_FindError(t *testing.T) {
	repo := new(MockRecurringRepository)
	txRepo := new(MockTxRepository)
	u := uc.New(repo, txRepo)

	repo.On("FindDue", mock.Anything).Return([]rt.RecurringTransaction{}, errors.New("db")).Once()

	err := u.ProcessDueTransactions()
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestRecurring_GetAndDelete(t *testing.T) {
	repo := new(MockRecurringRepository)
	txRepo := new(MockTxRepository)
	u := uc.New(repo, txRepo)

	repo.On("FindAllByUserID", int64(4)).Return([]rt.RecurringTransaction{{ID: 1}}, nil).Once()
	repo.On("Delete", int64(9)).Return(nil).Once()

	list, err := u.GetRecurring(4)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.NoError(t, u.DeleteRecurring(9))
	repo.AssertExpectations(t)
}
