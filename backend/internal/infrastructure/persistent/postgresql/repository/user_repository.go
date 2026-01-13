package postgresql

import (
	"database/sql"
	"errors"

	"github.com/afandimsr/cashbook-backend/internal/domain/user"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) user.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindAll(limit, offset int) ([]user.User, error) {
	rows, err := r.db.Query("SELECT id, name, email, is_active FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.IsActive)
		// fetch roles for user
		rRows, _ := r.db.Query("SELECT r.name FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = $1", u.ID)
		var roles []string
		for rRows.Next() {
			var rn string
			rRows.Scan(&rn)
			roles = append(roles, rn)
		}
		if rRows != nil {
			rRows.Close()
		}
		u.Roles = roles
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepo) FindByID(id int64) (user.User, error) {
	var u user.User
	var googleID sql.NullString
	err := r.db.QueryRow("SELECT id, name, email, password, google_id, is_active FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &googleID, &u.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("user not found")
		}
		return u, err
	}
	if googleID.Valid {
		u.GoogleID = googleID.String
	} else {
		u.GoogleID = ""
	}
	// populate roles
	rows, _ := r.db.Query("SELECT r.name FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = $1", u.ID)
	var roles []string
	for rows.Next() {
		var rn string
		rows.Scan(&rn)
		roles = append(roles, rn)
	}
	if rows != nil {
		rows.Close()
	}
	u.Roles = roles
	return u, nil
}

func (r *userRepo) Save(u user.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	var userID int64
	err = tx.QueryRow("INSERT INTO users(name, email, password, is_active) VALUES($1, $2, $3, $4) RETURNING id", u.Name, u.Email, u.Password, u.IsActive).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// assign roles if provided
	for _, roleName := range u.Roles {
		var roleID int64
		err = tx.QueryRow("SELECT id FROM roles WHERE name = $1", roleName).Scan(&roleID)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec("INSERT INTO user_roles(user_id, role_id) VALUES($1, $2) ON CONFLICT DO NOTHING", userID, roleID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *userRepo) Update(u user.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE users SET name = $1, email = $2, password = $3, is_active = $4 WHERE id = $5", u.Name, u.Email, u.Password, u.IsActive, u.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// reset roles
	_, err = tx.Exec("DELETE FROM user_roles WHERE user_id = $1", u.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, roleName := range u.Roles {
		var roleID int64
		err = tx.QueryRow("SELECT id FROM roles WHERE name = $1", roleName).Scan(&roleID)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec("INSERT INTO user_roles(user_id, role_id) VALUES($1, $2) ON CONFLICT DO NOTHING", u.ID, roleID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *userRepo) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (r *userRepo) FindByEmail(email string) (user.User, error) {
	var u user.User
	var googleID sql.NullString
	err := r.db.QueryRow("SELECT id, name, email, password, google_id, is_active FROM users WHERE email = $1", email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &googleID, &u.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("user not found")
		}
		return u, err
	}
	if googleID.Valid {
		u.GoogleID = googleID.String
	} else {
		u.GoogleID = ""
	}
	rows, _ := r.db.Query("SELECT r.name FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = $1", u.ID)
	var roles []string
	for rows.Next() {
		var rn string
		rows.Scan(&rn)
		roles = append(roles, rn)
	}
	if rows != nil {
		rows.Close()
	}
	u.Roles = roles
	return u, nil
}

func (r *userRepo) FindByGoogleID(googleID string) (user.User, error) {
	var u user.User
	var gID sql.NullString
	err := r.db.QueryRow("SELECT id, name, email, password, google_id, is_active FROM users WHERE google_id = $1", googleID).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &gID, &u.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, errors.New("user not found")
		}
		return u, err
	}
	if gID.Valid {
		u.GoogleID = gID.String
	} else {
		u.GoogleID = ""
	}
	rows, _ := r.db.Query("SELECT r.name FROM roles r JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id = $1", u.ID)
	var roles []string
	for rows.Next() {
		var rn string
		rows.Scan(&rn)
		roles = append(roles, rn)
	}
	if rows != nil {
		rows.Close()
	}
	u.Roles = roles
	return u, nil
}
