package budget_test

import (
	"errors"
	"testing"

	"github.com/afandimsr/cashbook-backend/internal/domain/budget"
	uc "github.com/afandimsr/cashbook-backend/internal/usecase/budget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBudgetRepository struct {
	mock.Mock
}

func (m *MockBudgetRepository) FindAllByUserID(userID int64, month, year int) ([]budget.Budget, error) {
	args := m.Called(userID, month, year)
	return args.Get(0).([]budget.Budget), args.Error(1)
}

func (m *MockBudgetRepository) FindByCategory(userID int64, categoryID int64, month, year int) (budget.Budget, error) {
	args := m.Called(userID, categoryID, month, year)
	return args.Get(0).(budget.Budget), args.Error(1)
}

func (m *MockBudgetRepository) Save(b *budget.Budget) error   { return m.Called(b).Error(0) }
func (m *MockBudgetRepository) Update(b *budget.Budget) error { return m.Called(b).Error(0) }
func (m *MockBudgetRepository) Delete(id int64) error         { return m.Called(id).Error(0) }

func TestBudget_GetBudgets(t *testing.T) {
	repo := new(MockBudgetRepository)
	u := uc.New(repo)

	repo.On("FindAllByUserID", int64(1), 5, 2026).Return([]budget.Budget{{ID: 1}}, nil).Once()

	got, err := u.GetBudgets(1, 5, 2026)
	assert.NoError(t, err)
	assert.Len(t, got, 1)
	repo.AssertExpectations(t)
}

func TestBudget_SetBudget_CreatesWhenMissing(t *testing.T) {
	repo := new(MockBudgetRepository)
	u := uc.New(repo)

	// FindByCategory returns an error → treated as "not found" → Save.
	repo.On("FindByCategory", int64(1), int64(7), 5, 2026).Return(budget.Budget{}, errors.New("no rows")).Once()
	repo.On("Save", mock.MatchedBy(func(b *budget.Budget) bool {
		return b.UserID == 1 && b.CategoryID == 7 && b.Amount == 500
	})).Return(nil).Once()

	err := u.SetBudget(1, budget.Budget{CategoryID: 7, Amount: 500, Month: 5, Year: 2026})
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestBudget_SetBudget_UpdatesWhenExisting(t *testing.T) {
	repo := new(MockBudgetRepository)
	u := uc.New(repo)

	existing := budget.Budget{ID: 3, UserID: 1, CategoryID: 7, Amount: 100, Month: 5, Year: 2026}
	repo.On("FindByCategory", int64(1), int64(7), 5, 2026).Return(existing, nil).Once()
	repo.On("Update", mock.MatchedBy(func(b *budget.Budget) bool {
		return b.ID == 3 && b.Amount == 900 // amount overwritten
	})).Return(nil).Once()

	err := u.SetBudget(1, budget.Budget{CategoryID: 7, Amount: 900, Month: 5, Year: 2026})
	assert.NoError(t, err)
	repo.AssertNotCalled(t, "Save", mock.Anything)
	repo.AssertExpectations(t)
}

func TestBudget_GetBudgetByCategory(t *testing.T) {
	repo := new(MockBudgetRepository)
	u := uc.New(repo)

	repo.On("FindByCategory", int64(1), int64(2), 1, 2026).Return(budget.Budget{ID: 8}, nil).Once()

	got, err := u.GetBudgetByCategory(1, 2, 1, 2026)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), got.ID)
	repo.AssertExpectations(t)
}
