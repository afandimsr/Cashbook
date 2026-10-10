package recurring_transaction_test

import (
	"context"
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

func (m *MockRecurringRepository) FindDueByUserID(userID int64, now time.Time) ([]rt.RecurringTransaction, error) {
	args := m.Called(userID, now)
	return args.Get(0).([]rt.RecurringTransaction), args.Error(1)
}

func (m *MockRecurringRepository) Save(r *rt.RecurringTransaction) error { return m.Called(r).Error(0) }
func (m *MockRecurringRepository) Update(r *rt.RecurringTransaction) error {
	return m.Called(r).Error(0)
}
func (m *MockRecurringRepository) Delete(id int64) error { return m.Called(id).Error(0) }
func (m *MockRecurringRepository) ProcessOne(ctx context.Context, id int64, tx *transaction.Transaction, now time.Time) error {
	return m.Called(ctx, id, tx, now).Error(0)
}

func TestRecurring_Create_SetsUserAndDefaultStart(t *testing.T) {
	repo := new(MockRecurringRepository)
	u := uc.New(repo)

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
	u := uc.New(repo)

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
	u := uc.New(repo)

	repo.On("Save", mock.Anything).Return(errors.New("db")).Once()

	_, err := u.CreateRecurring(1, rt.RecurringTransaction{})
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestRecurring_ProcessDue_CreatesTransactions(t *testing.T) {
	repo := new(MockRecurringRepository)
	u := uc.New(repo)

	due := []rt.RecurringTransaction{
		{ID: 1, UserID: 2, CategoryID: 3, Amount: 50, Type: "expense", Note: "Netflix", Frequency: rt.Monthly},
		{ID: 2, UserID: 2, CategoryID: 4, Amount: 20, Type: "expense", Note: "Spotify", Frequency: rt.Monthly},
	}
	repo.On("FindDueByUserID", int64(2), mock.Anything).Return(due, nil).Once()
	repo.On("ProcessOne", mock.Anything, int64(1), mock.MatchedBy(func(tx *transaction.Transaction) bool {
		return tx.Note == "Netflix (Auto-generated)" && tx.UserID == 2 && tx.CategoryID == 3
	}), mock.Anything).Run(func(args mock.Arguments) {
		args.Get(2).(*transaction.Transaction).ID = 100
	}).Return(nil).Once()
	repo.On("ProcessOne", mock.Anything, int64(2), mock.MatchedBy(func(tx *transaction.Transaction) bool {
		return tx.Note == "Spotify (Auto-generated)"
	}), mock.Anything).Return(nil).Once()

	got, err := u.ProcessDueForUser(context.Background(), 2)
	assert.NoError(t, err)
	assert.Equal(t, 0, got.Failed)
	assert.Len(t, got.Processed, 2)
	assert.Equal(t, uc.ProcessedItem{
		RecurringID: 1, TransactionID: 100, CategoryID: 3, Amount: 50, Type: "expense", Note: "Netflix", Frequency: rt.Monthly,
	}, got.Processed[0])
	repo.AssertExpectations(t)
}

func TestRecurring_ProcessDue_NothingDue(t *testing.T) {
	repo := new(MockRecurringRepository)
	u := uc.New(repo)

	repo.On("FindDueByUserID", int64(2), mock.Anything).Return([]rt.RecurringTransaction{}, nil).Once()

	got, err := u.ProcessDueForUser(context.Background(), 2)
	assert.NoError(t, err)
	assert.NotNil(t, got.Processed) // serialized as [] not null
	assert.Empty(t, got.Processed)
	assert.Equal(t, 0, got.Failed)
	repo.AssertNotCalled(t, "ProcessOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestRecurring_ProcessDue_SkipsAlreadyProcessed(t *testing.T) {
	repo := new(MockRecurringRepository)
	u := uc.New(repo)

	due := []rt.RecurringTransaction{{ID: 1, UserID: 2, Note: "x"}}
	repo.On("FindDueByUserID", int64(2), mock.Anything).Return(due, nil).Once()
	repo.On("ProcessOne", mock.Anything, int64(1), mock.Anything, mock.Anything).Return(rt.ErrAlreadyProcessed).Once()

	got, err := u.ProcessDueForUser(context.Background(), 2)
	assert.NoError(t, err)
	assert.Empty(t, got.Processed)
	assert.Equal(t, 0, got.Failed) // lost a race, not a failure
	repo.AssertExpectations(t)
}

func TestRecurring_ProcessDue_CountsFailuresAndContinues(t *testing.T) {
	repo := new(MockRecurringRepository)
	u := uc.New(repo)

	due := []rt.RecurringTransaction{{ID: 1, UserID: 2, Note: "a"}, {ID: 2, UserID: 2, Note: "b"}}
	repo.On("FindDueByUserID", int64(2), mock.Anything).Return(due, nil).Once()
	repo.On("ProcessOne", mock.Anything, int64(1), mock.Anything, mock.Anything).Return(errors.New("db")).Once()
	repo.On("ProcessOne", mock.Anything, int64(2), mock.Anything, mock.Anything).Return(nil).Once()

	got, err := u.ProcessDueForUser(context.Background(), 2)
	assert.NoError(t, err)
	assert.Equal(t, 1, got.Failed)
	assert.Len(t, got.Processed, 1)
	assert.Equal(t, int64(2), got.Processed[0].RecurringID)
	repo.AssertExpectations(t)
}

func TestRecurring_ProcessDue_FindError(t *testing.T) {
	repo := new(MockRecurringRepository)
	u := uc.New(repo)

	repo.On("FindDueByUserID", int64(2), mock.Anything).Return([]rt.RecurringTransaction{}, errors.New("db")).Once()

	_, err := u.ProcessDueForUser(context.Background(), 2)
	assert.Error(t, err)
	repo.AssertExpectations(t)
}

func TestRecurring_GetAndDelete(t *testing.T) {
	repo := new(MockRecurringRepository)
	u := uc.New(repo)

	repo.On("FindAllByUserID", int64(4)).Return([]rt.RecurringTransaction{{ID: 1}}, nil).Once()
	repo.On("Delete", int64(9)).Return(nil).Once()

	list, err := u.GetRecurring(4)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.NoError(t, u.DeleteRecurring(9))
	repo.AssertExpectations(t)
}
