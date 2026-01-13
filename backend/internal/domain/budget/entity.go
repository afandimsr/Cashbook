package budget

type Budget struct {
	ID         int64   `json:"id"`
	UserID     int64   `json:"user_id"`
	CategoryID int64   `json:"category_id"`
	Amount     float64 `json:"amount"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
}

type Repository interface {
	FindAllByUserID(userID int64, month, year int) ([]Budget, error)
	FindByCategory(userID int64, categoryID int64, month, year int) (Budget, error)
	Save(budget *Budget) error
	Update(budget *Budget) error
	Delete(id int64) error
}
