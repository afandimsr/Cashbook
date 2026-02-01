package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppVersion         string
	AppName            string
	AppPort            string
	AppEnv             string
	JWTSecret          string
	ClientAuthURL      string
	CorsAllowedOrigins string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	DB DBConfig
}

type DBConfig struct {
	Driver  string
	Host    string
	Port    string
	User    string
	Pass    string
	Name    string
	SSLMode string
}

func Load() *Config {
	// Load .env (ignore error in production)
	_ = godotenv.Load()

	cfg := &Config{
		AppVersion:    "1.8.0",
		AppName:       getEnv("APP_NAME", "go-app"),
		AppPort:       getEnv("APP_PORT", "8080"),
		AppEnv:        getEnv("APP_ENV", "production"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret"),
		ClientAuthURL: getEnv("CLIENT_AUTH_URL", ""),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),

		DB: DBConfig{
			Driver:  getEnv("DB_DRIVER", "postgres"),
			Host:    getEnv("DB_HOST", "localhost"),
			Port:    getEnv("DB_PORT", "5432"),
			User:    getEnv("DB_USER", "root"),
			Pass:    getEnv("DB_PASS", ""),
			Name:    getEnv("DB_NAME", ""),
			SSLMode: getEnv("DB_SSLMODE", "disable"),
		},
		CorsAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}

	validate(cfg)
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func validate(cfg *Config) {
	if cfg.DB.Name == "" {
		log.Fatal("DB_NAME is required")
	}
}
