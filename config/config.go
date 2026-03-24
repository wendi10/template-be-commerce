package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Payment  PaymentConfig
}

type AppConfig struct {
	Name     string
	Env      string
	Port     string
	BaseURL  string
	Debug    bool
	LogLevel string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	Charset      string
	ParseTime    bool
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// DSN returns the MySQL data source name in the format expected by go-sql-driver/mysql.
func (d DatabaseConfig) DSN() string {
	parseTime := "false"
	if d.ParseTime {
		parseTime = "true"
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=UTC&multiStatements=true",
		d.User, d.Password, d.Host, d.Port, d.Name, d.Charset, parseTime,
	)
}

type JWTConfig struct {
	AccessTokenSecret  string
	RefreshTokenSecret string
	AccessTokenTTL     time.Duration
	RefreshTokenTTL    time.Duration
}

type PaymentConfig struct {
	DokuClientID  string
	DokuSecretKey string
	DokuBaseURL   string
	CallbackURL   string
}

// Load reads configuration from environment variables.
// It loads a .env file if present (useful for local dev).
func Load() (*Config, error) {
	// Try to load .env file; ignore error if not found
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Name:     getEnv("APP_NAME", "ecommerce-api"),
			Env:      getEnv("APP_ENV", "development"),
			Port:     getEnv("APP_PORT", "8080"),
			BaseURL:  getEnv("APP_BASE_URL", "http://localhost:8080"),
			Debug:    getEnvBool("APP_DEBUG", true),
			LogLevel: getEnv("LOG_LEVEL", "info"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "3306"),
			User:         getEnv("DB_USER", "root"),
			Password:     getEnv("DB_PASSWORD", "root"),
			Name:         getEnv("DB_NAME", "ecommerce"),
			Charset:      getEnv("DB_CHARSET", "utf8mb4"),
			ParseTime:    true,
			MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxLifetime:  getEnvDuration("DB_MAX_LIFETIME", 5*time.Minute),
		},
		JWT: JWTConfig{
			AccessTokenSecret:  getEnv("JWT_ACCESS_SECRET", "change-me-access-secret"),
			RefreshTokenSecret: getEnv("JWT_REFRESH_SECRET", "change-me-refresh-secret"),
			AccessTokenTTL:     getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL:    getEnvDuration("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
		Payment: PaymentConfig{
			DokuClientID:  getEnv("DOKU_CLIENT_ID", ""),
			DokuSecretKey: getEnv("DOKU_SECRET_KEY", ""),
			DokuBaseURL:   getEnv("DOKU_BASE_URL", "https://api-sandbox.doku.com"),
			CallbackURL:   getEnv("PAYMENT_CALLBACK_URL", "http://localhost:8080/api/v1/payments/callback"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return b
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		i, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return i
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		d, err := time.ParseDuration(value)
		if err != nil {
			return defaultValue
		}
		return d
	}
	return defaultValue
}
