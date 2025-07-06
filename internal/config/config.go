package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Server struct {
		Host            string
		Port            string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	RateLimit struct {
		Request  int
		Duration time.Duration
	}
	CORS struct {
		AllowedOrigins []string
		AllowedMethods []string
		AllowedHeaders []string
	}
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	var cfg Config

	// Server
	cfg.Server.Host = getEnv("SERVER_HOST", "localhost")
	cfg.Server.Port = getEnv("SERVER_PORT", "4000")
	cfg.Server.ReadTimeout = time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 10)) * time.Second
	cfg.Server.WriteTimeout = time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second
	cfg.Server.IdleTimeout = time.Duration(getEnvAsInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second
	cfg.Server.ShutdownTimeout = time.Duration(getEnvAsInt("SERVER_SHUTDOWN_TIMEOUT", 5)) * time.Second

	// Database
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.Name = getEnv("DB_NAME", "")
	cfg.Database.SSLMode = getEnv("DB_SSL_MODE", "disable")
	// Rate Limit
	// CORS

	return &cfg
}

// ENV Helper
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
