package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// User model for database
type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"uniqueIndex;not null" json:"username"`
	Password    string    `gorm:"not null" json:"-"`
	Email       string    `json:"email"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	IsAdmin     bool      `gorm:"default:false" json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastLogin   *time.Time `json:"last_login"`
	LoginAttempts int     `gorm:"default:0" json:"-"`
}

// VPNUser model for VPN users
type VPNUser struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"uniqueIndex;not null" json:"username"`
	Protocol    string    `gorm:"not null" json:"protocol"` // ssh, vmess, vless, trojan, shadowsocks
	Password    string    `json:"password,omitempty"`
	UUID        string    `json:"uuid,omitempty"`
	Port        int       `json:"port,omitempty"`
	ExpiryDate  time.Time `json:"expiry_date"`
	CreatedDate time.Time `json:"created_date"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedBy   string    `json:"created_by"`
	Config      string    `gorm:"type:text" json:"config"` // JSON string for additional config
	
	// Traffic tracking
	UploadBytes   int64 `gorm:"default:0" json:"upload_bytes"`
	DownloadBytes int64 `gorm:"default:0" json:"download_bytes"`
	TotalBytes    int64 `gorm:"default:0" json:"total_bytes"`
	
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TrafficLog model for tracking bandwidth usage
type TrafficLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"index;not null" json:"username"`
	Protocol    string    `json:"protocol"`
	Upload      int64     `json:"upload"`
	Download    int64     `json:"download"`
	Total       int64     `json:"total"`
	Date        time.Time `gorm:"index" json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

// SystemLog model for system events
type SystemLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Level     string    `json:"level"` // info, warning, error
	Message   string    `json:"message"`
	Component string    `json:"component"` // api, vpn, system
	UserID    string    `json:"user_id,omitempty"`
	IPAddress string    `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Initialize database connection
func InitDB(dbPath string) error {
	// Create directory if not exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %v", err)
	}

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent, // Change to logger.Info for debug
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Configure SQLite for better performance
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Enable WAL mode for better concurrency
	DB.Exec("PRAGMA journal_mode=WAL;")
	DB.Exec("PRAGMA synchronous=NORMAL;")
	DB.Exec("PRAGMA cache_size=1000;")
	DB.Exec("PRAGMA foreign_keys=ON;")

	// Auto migrate schemas
	if err := autoMigrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	log.Printf("Database initialized successfully at: %s", dbPath)
	return nil
}

// Auto migrate all models
func autoMigrate() error {
	return DB.AutoMigrate(
		&User{},
		&VPNUser{},
		&TrafficLog{},
		&SystemLog{},
	)
}

// Close database connection
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health check for database
func HealthCheck() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Get database statistics
func GetDBStats() map[string]interface{} {
	sqlDB, err := DB.DB()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	stats := sqlDB.Stats()
	
	var userCount, vpnUserCount, trafficLogCount, systemLogCount int64
	DB.Model(&User{}).Count(&userCount)
	DB.Model(&VPNUser{}).Count(&vpnUserCount)
	DB.Model(&TrafficLog{}).Count(&trafficLogCount)
	DB.Model(&SystemLog{}).Count(&systemLogCount)

	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":              stats.InUse,
		"idle":                stats.Idle,
		"users_count":         userCount,
		"vpn_users_count":     vpnUserCount,
		"traffic_logs_count":  trafficLogCount,
		"system_logs_count":   systemLogCount,
	}
}

// Log system events
func LogEvent(level, message, component, userID, ipAddress string) {
	log := SystemLog{
		Level:     level,
		Message:   message,
		Component: component,
		UserID:    userID,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
	}
	
	DB.Create(&log)
}

// Clean old logs (keep last 30 days)
func CleanOldLogs() error {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	
	// Clean old traffic logs
	if err := DB.Where("created_at < ?", thirtyDaysAgo).Delete(&TrafficLog{}).Error; err != nil {
		return err
	}
	
	// Clean old system logs (keep only errors and warnings for longer)
	if err := DB.Where("created_at < ? AND level = ?", thirtyDaysAgo, "info").Delete(&SystemLog{}).Error; err != nil {
		return err
	}
	
	return nil
}