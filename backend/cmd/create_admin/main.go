package main

import (
	"database/sql"
	"log"

	"github.com/afandimsr/cashbook-backend/internal/config"
	"github.com/afandimsr/cashbook-backend/internal/seeder"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	dbUser := cfg.DB.User
	dbPassword := cfg.DB.Pass
	dbName := cfg.DB.Name
	dbHost := cfg.DB.Host
	dbPort := cfg.DB.Port
	dbSSLMode := cfg.DB.SSLMode

	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + dbSSLMode
	log.Printf("Connecting to DB at %s", dbURL)

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// Ensure connection
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	// Seed roles
	seeder.SeedRoles(db)
	// Seed the default admin user
	seeder.SeedAdminUser(db)
}
