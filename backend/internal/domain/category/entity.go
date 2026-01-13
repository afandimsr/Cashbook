package category

type Category struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Type   string `json:"type"` // "income" or "expense"
	Color  string `json:"color"`
	Icon   string `json:"icon"`
}

type Repository interface {
	FindAllByUserID(userID int64) ([]Category, error)
	FindByID(id int64) (Category, error)
	Save(category *Category) error
	Update(category *Category) error
	Delete(id int64) error
}
