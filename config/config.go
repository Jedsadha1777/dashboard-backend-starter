package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Configuration contains all app configuration
type Configuration struct {
	Database  DatabaseConfig
	Server    ServerConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
}

// DatabaseConfig contains database related configuration
type DatabaseConfig struct {
	User     string
	Password string
	Name     string
	Port     string
	Host     string
	TimeZone string
	SSLMode  string
}

// ServerConfig contains server related configuration
type ServerConfig struct {
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	TrustedProxies []string // เพิ่มส่วนนี้
}

// JWTConfig contains JWT related configuration
type JWTConfig struct {
	Secret        string
	ExpiryMinutes int
}

// RateLimitConfig สำหรับการตั้งค่าการจำกัดอัตราการเข้าถึง
type RateLimitConfig struct {
	RequestsPerMinute int
	LimitedPaths      []string // Add paths that should be rate-limited
}

var Config Configuration

// Init initializes application configuration
func Init() error {
	// Try to load .env file, but continue if not found
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database config
	Config.Database = DatabaseConfig{
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "dashboard"),
		Port:     getEnv("DB_PORT", "5432"),
		Host:     getEnv("DB_HOST", "localhost"),
		TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// ดึงค่า trusted proxies จาก env (เพิ่มส่วนนี้)
	trustedProxiesStr := getEnv("TRUSTED_PROXIES", "")
	var trustedProxies []string
	if trustedProxiesStr != "" {
		// แยก IP address/subnet ด้วยเครื่องหมาย ','
		trustedProxies = strings.Split(trustedProxiesStr, ",")
		// ตัด whitespace
		for i := range trustedProxies {
			trustedProxies[i] = strings.TrimSpace(trustedProxies[i])
		}
	}

	// Initialize server config
	Config.Server = ServerConfig{
		Port:           getEnv("SERVER_PORT", "8080"),
		ReadTimeout:    time.Duration(getEnvAsInt("SERVER_READ_TIMEOUT", 10)) * time.Second,
		WriteTimeout:   time.Duration(getEnvAsInt("SERVER_WRITE_TIMEOUT", 10)) * time.Second,
		TrustedProxies: trustedProxies, // เพิ่มส่วนนี้
	}

	// Initialize JWT config
	Config.JWT = JWTConfig{
		Secret:        getEnv("JWT_SECRET", ""),
		ExpiryMinutes: getEnvAsInt("JWT_EXPIRY_MINUTES", 1440), // Default 24 hours
	}

	// Validate critical configuration
	if Config.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}

	if Config.Database.Password == "" {
		log.Println("WARNING: Using empty database password. Set DB_PASSWORD for production environments")
	}

	// Parse rate limit paths from env
	rateLimitPathsStr := getEnv("RATE_LIMIT_PATHS", "/api/v1/auth/login")
	rateLimitPaths := strings.Split(rateLimitPathsStr, ",")
	// Trim whitespace from each path
	for i := range rateLimitPaths {
		rateLimitPaths[i] = strings.TrimSpace(rateLimitPaths[i])
	}

	Config.RateLimit = RateLimitConfig{
		RequestsPerMinute: getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 60),
		LimitedPaths:      rateLimitPaths,
	}

	return nil
}

// getEnv retrieves environment variable with fallback value
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsInt retrieves environment variable as integer with fallback
func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}

	value := 0
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		log.Printf("WARNING: Environment variable %s is not a valid integer. Using default value %d", key, fallback)
		return fallback
	}
	return value
}

// GetDSN returns database connection string
func GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		Config.Database.Host,
		Config.Database.User,
		Config.Database.Password,
		Config.Database.Name,
		Config.Database.Port,
		Config.Database.SSLMode,
		Config.Database.TimeZone)
}
