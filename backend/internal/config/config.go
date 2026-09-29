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
	FrontendURL        string
	CorsAllowedOrigins string

	// BotInternalAPIKey authenticates the Telegram bot service on the
	// /internal/bot/* routes (see middleware.ServiceAuthMiddleware). Left empty
	// by default so existing deployments without a bot keep working unchanged;
	// ServiceAuthMiddleware denies every request while it's empty.
	BotInternalAPIKey string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	DB         DBConfig
	ElasticApm ElasticApmConfig
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

type ElasticApmConfig struct {
	ServerURL        string
	ServiceName      string
	Environment      string
	SecretToken      string
	VerifyServerCert bool
	ServiceVersion   string
}

func Load() *Config {
	// Load .env (ignore error in production)
	_ = godotenv.Load()

	cfg := &Config{
		AppVersion:    "1.17.0",
		AppName:       getEnv("APP_NAME", "go-app"),
		AppPort:       getEnv("APP_PORT", "8080"),
		AppEnv:        getEnv("APP_ENV", "production"),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret"),
		ClientAuthURL: getEnv("CLIENT_AUTH_URL", ""),
		FrontendURL:   getEnv("FRONTEND_URL", ""),

		BotInternalAPIKey: getEnv("BOT_INTERNAL_API_KEY", ""),

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
		ElasticApm: ElasticApmConfig{
			ServerURL:        getEnv("ELASTIC_APM_SERVER_URL", ""),
			ServiceName:      getEnv("ELASTIC_APM_SERVICE_NAME", "go-app"),
			Environment:      getEnv("ELASTIC_APM_ENVIRONMENT", "production"),
			SecretToken:      getEnv("ELASTIC_APM_SECRET_TOKEN", ""),
			VerifyServerCert: getEnv("ELASTIC_APM_VERIFY_SERVER_CERT", "true") == "true",
			ServiceVersion:   getEnv("ELASTIC_APM_SERVICE_VERSION", "1.10.0"),
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
