package services

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/nabilulilalbab/apivpn/internal/models"
)

type UserService struct {
	usersFile string
}

func NewUserService() *UserService {
	return &UserService{
		usersFile: "/etc/apivpn/users.json",
	}
}

type AdminUser struct {
	Username     string    `json:"username"`
	PasswordHash string    `json:"password_hash"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
	LastLogin    time.Time `json:"last_login"`
	IsActive     bool      `json:"is_active"`
}

// Login authenticates admin user
func (u *UserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	users, err := u.loadUsers()
	if err != nil {
		return nil, err
	}

	user, exists := users[req.Username]
	if !exists {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is disabled")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, expiresAt, err := u.generateToken(req.Username)
	if err != nil {
		return nil, err
	}

	// Update last login
	user.LastLogin = time.Now()
	users[req.Username] = user
	u.saveUsers(users)

	return &models.LoginResponse{
		Token:     token,
		Username:  req.Username,
		ExpiresAt: expiresAt,
	}, nil
}

// Register creates a new admin user
func (u *UserService) Register(req *models.LoginRequest) error {
	users, err := u.loadUsers()
	if err != nil {
		users = make(map[string]AdminUser)
	}

	// Check if user already exists
	if _, exists := users[req.Username]; exists {
		return fmt.Errorf("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Create user
	user := AdminUser{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		IsActive:     true,
	}

	users[req.Username] = user
	return u.saveUsers(users)
}

// ChangePassword changes user password
func (u *UserService) ChangePassword(username, oldPassword, newPassword string) error {
	users, err := u.loadUsers()
	if err != nil {
		return err
	}

	user, exists := users[username]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return fmt.Errorf("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	users[username] = user

	return u.saveUsers(users)
}

// GetUserInfo returns user information
func (u *UserService) GetUserInfo(username string) (*AdminUser, error) {
	users, err := u.loadUsers()
	if err != nil {
		return nil, err
	}

	user, exists := users[username]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Don't return password hash
	user.PasswordHash = ""
	return &user, nil
}

// CreateDefaultAdmin creates default admin user if no users exist
func (u *UserService) CreateDefaultAdmin() error {
	users, err := u.loadUsers()
	if err != nil || len(users) == 0 {
		// Create default admin
		defaultPassword := u.generateRandomPassword()
		
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		users = map[string]AdminUser{
			"admin": {
				Username:     "admin",
				PasswordHash: string(hashedPassword),
				CreatedAt:    time.Now(),
				IsActive:     true,
			},
		}

		if err := u.saveUsers(users); err != nil {
			return err
		}

		// Write default credentials to file
		credFile := "/etc/apivpn/default_credentials.txt"
		os.MkdirAll("/etc/apivpn", 0755)
		content := fmt.Sprintf("Default Admin Credentials:\nUsername: admin\nPassword: %s\n\nPlease change this password after first login!\n", defaultPassword)
		os.WriteFile(credFile, []byte(content), 0600)

		fmt.Printf("Default admin created. Credentials saved to %s\n", credFile)
	}

	return nil
}

// Helper methods
func (u *UserService) loadUsers() (map[string]AdminUser, error) {
	users := make(map[string]AdminUser)

	if _, err := os.Stat(u.usersFile); os.IsNotExist(err) {
		return users, nil
	}

	data, err := os.ReadFile(u.usersFile)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return users, nil
	}

	err = json.Unmarshal(data, &users)
	return users, err
}

func (u *UserService) saveUsers(users map[string]AdminUser) error {
	// Ensure directory exists
	os.MkdirAll("/etc/apivpn", 0755)

	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(u.usersFile, data, 0600)
}

func (u *UserService) generateToken(username string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"username": username,
		"user_id":  username,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your-secret-key-change-this"))
	
	return tokenString, expiresAt, err
}

func (u *UserService) generateRandomPassword() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:12]
}