package user

type UserRepository interface {
	FindAll(limit, offset int) ([]User, error)
	FindByID(id int64) (User, error)
	FindByEmail(email string) (User, error)
	FindByGoogleID(googleID string) (User, error)
	Save(user User) error
	Update(user User) error
	Delete(id int64) error
}

type AuthService interface {
	Login(email, password string) (bool, error)
}
