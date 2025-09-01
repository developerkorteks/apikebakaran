package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nabilulilalbab/apivpn/internal/database"
	"github.com/nabilulilalbab/apivpn/internal/models"
)

type DatabaseUserService struct {
	jwtSecret        string
	bcryptCost       int
	tokenExpireHours int
	maxLoginAttempts int
}

func NewDatabaseUserService(jwtSecret string, bcryptCost, tokenExpireHours, maxLoginAttempts int) *DatabaseUserService {
	return &DatabaseUserService{
		jwtSecret:        jwtSecret,
		bcryptCost:       bcryptCost,
		tokenExpireHours: tokenExpireHours,
		maxLoginAttempts: maxLoginAttempts,
	}
}

// Login authenticates admin user
func (u *DatabaseUserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	var user database.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		database.LogEvent("warning", fmt.Sprintf("Login attempt for non-existent user: %s", req.Username), "auth", "", "")
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if account is locked due to too many failed attempts
	if user.LoginAttempts >= u.maxLoginAttempts {
		database.LogEvent("warning", fmt.Sprintf("Account locked due to too many failed attempts: %s", req.Username), "auth", req.Username, "")
		return nil, fmt.Errorf("account locked due to too many failed login attempts")
	}

	if !user.IsActive {
		database.LogEvent("warning", fmt.Sprintf("Login attempt for disabled account: %s", req.Username), "auth", req.Username, "")
		return nil, fmt.Errorf("account is disabled")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// Increment failed login attempts
		database.DB.Model(&user).Update("login_attempts", user.LoginAttempts+1)
		database.LogEvent("warning", fmt.Sprintf("Failed login attempt for user: %s", req.Username), "auth", req.Username, "")
		return nil, fmt.Errorf("invalid credentials")
	}

	// Reset login attempts and update last login
	now := time.Now()
	database.DB.Model(&user).Updates(map[string]interface{}{
		"login_attempts": 0,
		"last_login":     &now,
	})

	// Generate JWT token
	token, expiresAt, err := u.generateToken(req.Username)
	if err != nil {
		return nil, err
	}

	database.LogEvent("info", fmt.Sprintf("Successful login for user: %s", req.Username), "auth", req.Username, "")

	return &models.LoginResponse{
		Token:     token,
		Username:  req.Username,
		ExpiresAt: expiresAt,
	}, nil
}

// Register creates a new admin user
func (u *DatabaseUserService) Register(req *models.LoginRequest) error {
	// Check if user already exists
	var existingUser database.User
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), u.bcryptCost)
	if err != nil {
		return err
	}

	// Create user
	user := database.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		IsActive:  true,
		IsAdmin:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}

	database.LogEvent("info", fmt.Sprintf("New admin user registered: %s", req.Username), "auth", req.Username, "")
	return nil
}

// ChangePassword changes user password
func (u *DatabaseUserService) ChangePassword(username, oldPassword, newPassword string) error {
	var user database.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), u.bcryptCost)
	if err != nil {
		return err
	}

	// Update password
	if err := database.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		return err
	}

	database.LogEvent("info", fmt.Sprintf("Password changed for user: %s", username), "auth", username, "")
	return nil
}

// GetUserInfo returns user information
func (u *DatabaseUserService) GetUserInfo(username string) (*database.User, error) {
	var user database.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Don't return password hash
	user.Password = ""
	return &user, nil
}

// CreateDefaultAdmin creates default admin user if no users exist
func (u *DatabaseUserService) CreateDefaultAdmin() error {
	var count int64
	database.DB.Model(&database.User{}).Count(&count)
	
	if count == 0 {
		// Generate random password
		defaultPassword := u.generateRandomPassword()
		
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), u.bcryptCost)
		if err != nil {
			return err
		}

		user := database.User{
			Username:  "admin",
			Password:  string(hashedPassword),
			IsActive:  true,
			IsAdmin:   true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := database.DB.Create(&user).Error; err != nil {
			return err
		}

		// Write default credentials to file
		credFile := "/etc/apivpn/default_credentials.txt"
		content := fmt.Sprintf("Default Admin Credentials:\nUsername: admin\nPassword: %s\n\nPlease change this password after first login!\n", defaultPassword)
		
		if err := writeToFile(credFile, content); err != nil {
			return err
		}

		database.LogEvent("info", "Default admin user created", "system", "admin", "")
		fmt.Printf("Default admin created. Credentials saved to %s\n", credFile)
	}

	return nil
}

// ListUsers returns all admin users
func (u *DatabaseUserService) ListUsers() ([]database.User, error) {
	var users []database.User
	if err := database.DB.Select("id, username, email, is_active, is_admin, created_at, updated_at, last_login").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUserStatus updates user active status
func (u *DatabaseUserService) UpdateUserStatus(username string, isActive bool) error {
	if err := database.DB.Model(&database.User{}).Where("username = ?", username).Update("is_active", isActive).Error; err != nil {
		return err
	}
	
	status := "disabled"
	if isActive {
		status = "enabled"
	}
	database.LogEvent("info", fmt.Sprintf("User %s %s", username, status), "admin", username, "")
	return nil
}

// DeleteUser deletes a user (except admin)
func (u *DatabaseUserService) DeleteUser(username string) error {
	if username == "admin" {
		return fmt.Errorf("cannot delete admin user")
	}
	
	if err := database.DB.Where("username = ?", username).Delete(&database.User{}).Error; err != nil {
		return err
	}
	
	database.LogEvent("info", fmt.Sprintf("User deleted: %s", username), "admin", username, "")
	return nil
}

// Helper methods
func (u *DatabaseUserService) generateToken(username string) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(u.tokenExpireHours) * time.Hour)

	claims := jwt.MapClaims{
		"username": username,
		"user_id":  username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.jwtSecret))
	
	return tokenString, expiresAt, err
}

func (u *DatabaseUserService) generateRandomPassword() string {
	return generateRandomString(12)
}

// Helper functions
func generateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

func writeToFile(filename, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}