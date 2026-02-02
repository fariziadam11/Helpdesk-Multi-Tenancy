package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for the backend service.
type Config struct {
	// Application
	AppEnv string

	// Server
	ServerPort string

	// Database
	DBUser string
	DBPass string
	DBHost string
	DBPort string
	DBName string

	// JWT
	JWTSecret string

	// InvGate Armmada
	ArmMadaBaseURL  string
	ArmMadaUsername string
	ArmMadaPassword string
	ArmMadaPageKey  string

	ArmMadaCompanyID  int
	ArmMadaGroupID    int
	ArmMadaLocationID int

	// Logging
	LogLevel  string
	GinMode   string
	LogFormat string

	// Mailgun Email
	MailgunDomain string
	MailgunAPIKey string
	MailgunSender string

	// Frontend
	FrontendURL string
}

// Load loads configuration from environment variables (optionally via .env files).
func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env") // best-effort when running from cmd/

	cfg := &Config{
		AppEnv:            getEnv("APP_ENV", "development"),
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DBUser:            getEnv("DB_USER", "root"),
		DBPass:            getEnv("DB_PASSWORD", ""),
		DBHost:            getEnv("DB_HOST", "db"),
		DBPort:            getEnv("DB_PORT", "3306"),
		DBName:            getEnv("DB_NAME", "armmada"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		ArmMadaBaseURL:    getEnv("ARMMADA_BASE_URL", ""),
		ArmMadaUsername:   getEnv("ARMMADA_USERNAME", ""),
		ArmMadaPassword:   getEnv("ARMMADA_PASSWORD", ""),
		ArmMadaPageKey:    getEnv("ARMMADA_PAGE_KEY", ""),
		ArmMadaCompanyID:  getEnvInt("ARMMADA_COMPANY_ID", 135),
		ArmMadaGroupID:    getEnvInt("ARMMADA_GROUP_ID", 134),
		ArmMadaLocationID: getEnvInt("ARMMADA_LOCATION_ID", 136),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		GinMode:           getEnv("GIN_MODE", "debug"),
		LogFormat:         getEnv("LOG_FORMAT", "text"),
		MailgunDomain:     getEnv("MAILGUN_DOMAIN", "mg.werk.co.id"),
		MailgunAPIKey:     getEnv("MAILGUN_API_KEY", ""),
		MailgunSender:     getEnv("MAILGUN_SENDER", "Werk <no-reply@mg.werk.co.id>"),
		FrontendURL:       getEnv("FRONTEND_URL", "http://localhost:5173"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET must be provided")
	}

	if cfg.ArmMadaBaseURL == "" || cfg.ArmMadaUsername == "" || cfg.ArmMadaPassword == "" {
		return nil, fmt.Errorf("InvGate ARMMADA credentials must be provided")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return parsed
}
