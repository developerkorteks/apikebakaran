package config

import (
	"os"
)

type Config struct {
	Port      string
	JWTSecret string
	Domain    string
	XrayPath  string
	SSHPath   string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8080"),
		JWTSecret: getEnv("JWT_SECRET", "your-secret-key-change-this"),
		Domain:    getEnv("DOMAIN", ""),
		XrayPath:  getEnv("XRAY_PATH", "/etc/xray"),
		SSHPath:   getEnv("SSH_PATH", "/etc/ssh"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}