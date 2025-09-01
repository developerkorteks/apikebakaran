package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// API Configuration
	Port    string
	Host    string
	JWTSecret string
	
	// Server Configuration
	Domain   string
	XrayPath string
	SSHPath  string
	
	// Database Configuration
	DBType string
	DBPath string
	
	// Security Configuration
	BCryptCost         int
	TokenExpireHours   int
	MaxLoginAttempts   int
	RateLimitRequests  int
	RateLimitWindow    int
	
	// Logging Configuration
	LogLevel string
	LogFile  string
}

func Load() *Config {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	return &Config{
		// API Configuration
		Port:      getEnv("PORT", "37849"),
		Host:      getEnv("API_HOST", "0.0.0.0"),
		JWTSecret: getEnv("JWT_SECRET", "vpn-api-default-secret-change-this"),
		
		// Server Configuration
		Domain:   getEnv("DOMAIN", ""),
		XrayPath: getEnv("XRAY_PATH", "/etc/xray"),
		SSHPath:  getEnv("SSH_PATH", "/etc/ssh"),
		
		// Database Configuration
		DBType: getEnv("DB_TYPE", "sqlite"),
		DBPath: getEnv("DB_PATH", "/etc/apivpn/vpnapi.db"),
		
		// Security Configuration
		BCryptCost:        getEnvInt("BCRYPT_COST", 12),
		TokenExpireHours:  getEnvInt("TOKEN_EXPIRE_HOURS", 24),
		MaxLoginAttempts:  getEnvInt("MAX_LOGIN_ATTEMPTS", 5),
		RateLimitRequests: getEnvInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow:   getEnvInt("RATE_LIMIT_WINDOW", 60),
		
		// Logging Configuration
		LogLevel: getEnv("LOG_LEVEL", "info"),
		LogFile:  getEnv("LOG_FILE", "/var/log/apivpn/api.log"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}