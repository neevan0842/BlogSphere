package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database configuration
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_DB       string
	DATABASE_URL      string

	// Server configuration
	ADDR string

	// Database connection pool configuration
	DB_MAX_OPEN_CONNS int32
	DB_MAX_IDLE_CONNS int32
	DB_MAX_IDLE_TIME  string

	// Google OAuth configuration
	GOOGLE_CLIENT_ID     string
	GOOGLE_CLIENT_SECRET string
	GOOGLE_REDIRECT_URI  string

	// JWT Configuration
	JWT_SECRET                   string
	ACCESS_TOKEN_EXPIRE_MINUTES  int64
	REFRESH_TOKEN_EXPIRE_MINUTES int64

	// Cookie Configuration
	Secure bool

	// CORS Configuration
	CORS_ALLOWED_ORIGIN string

	// MailerSend Configuration
	MAILERSEND_API_KEY string
	FROM_EMAIL         string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		// Database configuration
		POSTGRES_USER:     getEnv("POSTGRES_USER", "postgres"),
		POSTGRES_PASSWORD: getEnv("POSTGRES_PASSWORD", "password"),
		POSTGRES_HOST:     getEnv("POSTGRES_HOST", "localhost"),
		POSTGRES_PORT:     getEnv("POSTGRES_PORT", "5433"),
		POSTGRES_DB:       getEnv("POSTGRES_DB", "blogsphere"),
		DATABASE_URL:      getEnv("DATABASE_URL", ""),

		// Server configuration
		ADDR: getEnv("ADDR", ":8080"),

		// Database connection pool configuration
		DB_MAX_OPEN_CONNS: int32(getEnvAsInt("DB_MAX_OPEN_CONNS", 30)),
		DB_MAX_IDLE_CONNS: int32(getEnvAsInt("DB_MAX_IDLE_CONNS", 30)),
		DB_MAX_IDLE_TIME:  getEnv("DB_MAX_IDLE_TIME", "15m"),

		// Google OAuth configuration
		GOOGLE_CLIENT_ID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GOOGLE_CLIENT_SECRET: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GOOGLE_REDIRECT_URI:  getEnv("GOOGLE_REDIRECT_URI", ""),

		// JWT Configuration
		JWT_SECRET:                   getEnv("JWT_SECRET", ""),
		ACCESS_TOKEN_EXPIRE_MINUTES:  getEnvAsInt("ACCESS_TOKEN_EXPIRE_MINUTES", 1440),
		REFRESH_TOKEN_EXPIRE_MINUTES: getEnvAsInt("REFRESH_TOKEN_EXPIRE_MINUTES", 10080),

		// Cookie Configuration
		Secure: getEnvAsBool("Secure", true),

		// CORS Configuration
		CORS_ALLOWED_ORIGIN: getEnv("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),

		// MailerSend Configuration
		MAILERSEND_API_KEY: getEnv("MAILERSEND_API_KEY", ""),
		FROM_EMAIL:         getEnv("FROM_EMAIL", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.Atoi(value)
		if err != nil {
			return int64(fallback)
		}
		return int64(i)
	}
	return int64(fallback)
}

func getEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return b
	}
	return fallback
}
