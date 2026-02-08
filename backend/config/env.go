package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_DB       string
	DATABASE_URL      string
	ADDR              string
	DB_MAX_OPEN_CONNS int32
	DB_MAX_IDLE_CONNS int32
	DB_MAX_IDLE_TIME  string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		POSTGRES_USER:     getEnv("POSTGRES_USER", "postgres"),
		POSTGRES_PASSWORD: getEnv("POSTGRES_PASSWORD", "password"),
		POSTGRES_HOST:     getEnv("POSTGRES_HOST", "localhost"),
		POSTGRES_PORT:     getEnv("POSTGRES_PORT", "5433"),
		POSTGRES_DB:       getEnv("POSTGRES_DB", "blogsphere"),
		DATABASE_URL:      getEnv("DATABASE_URL", ""),
		ADDR:              getEnv("ADDR", ":8080"),
		DB_MAX_OPEN_CONNS: int32(getEnvAsInt("DB_MAX_OPEN_CONNS", 30)),
		DB_MAX_IDLE_CONNS: int32(getEnvAsInt("DB_MAX_IDLE_CONNS", 30)),
		DB_MAX_IDLE_TIME:  getEnv("DB_MAX_IDLE_TIME", "15m"),
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
