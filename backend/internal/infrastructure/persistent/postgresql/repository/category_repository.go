package postgresql

import (
	"database/sql"

	"github.com/afandimsr/cashbook-backend/internal/domain/category"
)

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) category.Repository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) FindAllByUserID(userID int64) ([]category.Category, error) {
	rows, err := r.db.Query("SELECT id, user_id, name, type, color, icon FROM categories WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []category.Category
	for rows.Next() {
		var c category.Category
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.Color, &c.Icon)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *categoryRepo) FindByID(id int64) (category.Category, error) {
	var c category.Category
	err := r.db.QueryRow("SELECT id, user_id, name, type, color, icon FROM categories WHERE id = $1", id).Scan(&c.ID, &c.UserID, &c.Name, &c.Type, &c.Color, &c.Icon)
	return c, err
}

func (r *categoryRepo) Save(c *category.Category) error {
	return r.db.QueryRow(
		"INSERT INTO categories(user_id, name, type, color, icon) VALUES($1, $2, $3, $4, $5) RETURNING id",
		c.UserID, c.Name, c.Type, c.Color, c.Icon,
	).Scan(&c.ID)
}

func (r *categoryRepo) Update(c *category.Category) error {
	_, err := r.db.Exec(
		"UPDATE categories SET name = $1, type = $2, color = $3, icon = $4 WHERE id = $5",
		c.Name, c.Type, c.Color, c.Icon, c.ID,
	)
	return err
}

func (r *categoryRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}
