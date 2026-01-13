package seeder

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func SeedRoles(db *sql.DB) {

	roles := []string{"ADMIN", "USER"}

	for _, role := range roles {
		var roleID int64
		err := db.QueryRow("SELECT id FROM roles WHERE name = $1", role).Scan(&roleID)
		if err != nil {
			if err == sql.ErrNoRows {
				// Insert role
				_, err = db.Exec("INSERT INTO roles(name) VALUES($1)", role)
				if err != nil {
					log.Fatalf("failed to insert role %s: %v", role, err)
				}
			}
		} else {
			log.Printf("Role %s already exists with id %d", role, roleID)
		}
	}
}

// SeedAdminUser seeds a default admin user into the database.
func SeedAdminUser(db *sql.DB) {
	name := "Admin"
	email := "admin@example.com"
	password := "admin123"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	var userID int64
	// Check if user exists
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Insert user
			err = db.QueryRow("INSERT INTO users(name, email, password, is_active) VALUES($1, $2, $3, $4) RETURNING id", name, email, string(hashed), true).Scan(&userID)
			if err != nil {
				log.Fatalf("failed to insert user: %v", err)
			}
		} else {
			log.Fatalf("failed to query user: %v", err)
		}
	}

	// Find ADMIN role id
	var roleID int64
	err = db.QueryRow("SELECT id FROM roles WHERE name = $1", "ADMIN").Scan(&roleID)
	if err != nil {
		log.Fatalf("failed to find ADMIN role: %v", err)
	}

	// Insert into user_roles
	_, err = db.Exec("INSERT INTO user_roles(user_id, role_id) VALUES($1, $2) ON CONFLICT DO NOTHING", userID, roleID)
	if err != nil {
		log.Fatalf("failed to insert user_roles: %v", err)
	}

	log.Printf("Admin user seeded: id=%d, email=%s", userID, email)
}
