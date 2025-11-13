package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Config holds all application configuration
type Config struct {
	// Application
	AppEnv  string
	AppPort string
	AppName string

	// Database
	DatabaseURL         string
	DBMaxOpenConns      int
	DBMaxIdleConns      int
	DBConnMaxLifetime   time.Duration

	// Redis
	RedisURL      string
	RedisPassword string
	RedisDB       int

	// JWT
	JWTSecret           string
	JWTExpiry           time.Duration
	RefreshTokenExpiry  time.Duration

	// Supabase
	SupabaseURL        string
	SupabaseServiceKey string
	SupabaseJWTSecret  string

	// CORS
	CORSAllowedOrigins []string
	CORSAllowedMethods []string
	CORSAllowedHeaders []string

	// Logging
	LogLevel  string
	LogFormat string
	LogOutput string

	// Rate Limiting
	RateLimitRequests int
	RateLimitWindow   time.Duration
	RateLimitEnabled  bool

	// Feature Flags
	EnableAnalytics bool
	EnableDebug     bool
	EnableProfiling bool
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		// Application
		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
		AppName: getEnv("APP_NAME", "lumen-backend"),

		// Database
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		DBMaxOpenConns:     getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:     getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime:  getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),

		// Redis
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// JWT
		JWTSecret:          getEnv("JWT_SECRET", ""),
		JWTExpiry:          getEnvAsDuration("JWT_EXPIRY", 24*time.Hour),
		RefreshTokenExpiry: getEnvAsDuration("REFRESH_TOKEN_EXPIRY", 168*time.Hour),

		// Supabase
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_KEY", ""),
		SupabaseJWTSecret:  getEnv("SUPABASE_JWT_SECRET", ""),

		// CORS
		CORSAllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		CORSAllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		CORSAllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Authorization"}),

		// Logging
		LogLevel:  getEnv("LOG_LEVEL", "debug"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
		LogOutput: getEnv("LOG_OUTPUT", "stdout"),

		// Rate Limiting
		RateLimitRequests: getEnvAsInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvAsDuration("RATE_LIMIT_WINDOW", 60*time.Second),
		RateLimitEnabled:  getEnvAsBool("RATE_LIMIT_ENABLED", true),

		// Feature Flags
		EnableAnalytics: getEnvAsBool("ENABLE_ANALYTICS", false),
		EnableDebug:     getEnvAsBool("ENABLE_DEBUG", false),
		EnableProfiling: getEnvAsBool("ENABLE_PROFILING", false),
	}
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	var value int
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return valueStr == "true" || valueStr == "1"
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
