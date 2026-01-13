package budget

import (
	"github.com/afandimsr/cashbook-backend/internal/domain/budget"
)

type Usecase interface {
	GetBudgets(userID int64, month, year int) ([]budget.Budget, error)
	SetBudget(userID int64, b budget.Budget) error
	GetBudgetByCategory(userID int64, categoryID int64, month, year int) (budget.Budget, error)
}

type usecase struct {
	repo budget.Repository
}

func New(repo budget.Repository) Usecase {
	return &usecase{repo: repo}
}

func (u *usecase) GetBudgets(userID int64, month, year int) ([]budget.Budget, error) {
	return u.repo.FindAllByUserID(userID, month, year)
}

func (u *usecase) SetBudget(userID int64, b budget.Budget) error {
	b.UserID = userID

	existing, err := u.repo.FindByCategory(userID, b.CategoryID, b.Month, b.Year)
	if err == nil {
		// Update existing
		existing.Amount = b.Amount
		return u.repo.Update(&existing)
	}

	// Create new
	return u.repo.Save(&b)
}

func (u *usecase) GetBudgetByCategory(userID int64, categoryID int64, month, year int) (budget.Budget, error) {
	return u.repo.FindByCategory(userID, categoryID, month, year)
}
