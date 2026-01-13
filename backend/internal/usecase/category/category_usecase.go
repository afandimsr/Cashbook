package category

import (
	"github.com/afandimsr/cashbook-backend/internal/domain/category"
)

type Usecase interface {
	GetAllByUserID(userID int64) ([]category.Category, error)
	GetByID(id int64) (category.Category, error)
	Create(c category.Category) error
	Update(id int64, c category.Category) error
	Delete(id int64) error
}

type usecase struct {
	repo category.Repository
}

func New(repo category.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) GetAllByUserID(userID int64) ([]category.Category, error) {
	return u.repo.FindAllByUserID(userID)
}

func (u *usecase) GetByID(id int64) (category.Category, error) {
	return u.repo.FindByID(id)
}

func (u *usecase) Create(c category.Category) error {
	return u.repo.Save(&c)
}

func (u *usecase) Update(id int64, c category.Category) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}

	existing.Name = c.Name
	existing.Type = c.Type
	existing.Color = c.Color
	existing.Icon = c.Icon

	return u.repo.Update(&existing)
}

func (u *usecase) Delete(id int64) error {
	return u.repo.Delete(id)
}
