package models

import (
	"time"
)

// User represents a VPN user
type User struct {
	Username    string    `json:"username"`
	Password    string    `json:"password,omitempty"`
	Email       string    `json:"email,omitempty"`
	ExpiryDate  time.Time `json:"expiry_date"`
	CreatedDate time.Time `json:"created_date"`
	IsActive    bool      `json:"is_active"`
	Protocol    string    `json:"protocol"` // ssh, vmess, vless, trojan, shadowsocks
	UUID        string    `json:"uuid,omitempty"`
	Port        int       `json:"port,omitempty"`
}

// CreateUserRequest represents request to create a new user
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Days     int    `json:"days" binding:"required,min=1"`
	Protocol string `json:"protocol,omitempty"` // Set by handler, not required in request
}

// ChangePasswordRequest represents request to change password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdateUserStatusRequest represents request to update user status
type UpdateUserStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// ExtendUserRequest represents request to extend user expiration
type ExtendUserRequest struct {
	Days int `json:"days" binding:"required,min=1"`
}

// AddDomainRequest represents request to add a new domain
type AddDomainRequest struct {
	Domain string `json:"domain" binding:"required"`
}

// SystemInfo represents system information
type SystemInfo struct {
	OS              string  `json:"os"`
	Kernel          string  `json:"kernel"`
	CPUName         string  `json:"cpu_name"`
	CPUCores        int     `json:"cpu_cores"`
	CPUUsage        string  `json:"cpu_usage"`
	RAMUsed         int     `json:"ram_used_mb"`
	RAMTotal        int     `json:"ram_total_mb"`
	RAMUsage        string  `json:"ram_usage_percent"`
	Uptime          string  `json:"uptime"`
	Domain          string  `json:"domain"`
	IP              string  `json:"ip"`
	DailyBandwidth  string  `json:"daily_bandwidth"`
	MonthlyBandwidth string `json:"monthly_bandwidth"`
}

// ServiceStatus represents status of VPN services
type ServiceStatus struct {
	SSH       bool `json:"ssh"`
	Nginx     bool `json:"nginx"`
	Xray      bool `json:"xray"`
	Dropbear  bool `json:"dropbear"`
	Stunnel   bool `json:"stunnel"`
	SSHWebSocket bool `json:"ssh_websocket"`
}

// UserTraffic represents user bandwidth usage
type UserTraffic struct {
	Username string `json:"username"`
	Upload   string `json:"upload"`
	Download string `json:"download"`
	Total    string `json:"total"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
}

// VPNConfig represents VPN configuration for client
type VPNConfig struct {
	Protocol string            `json:"protocol"`
	Server   string            `json:"server"`
	Port     int               `json:"port"`
	Username string            `json:"username"`
	Password string            `json:"password,omitempty"`
	UUID     string            `json:"uuid,omitempty"`
	Config   map[string]string `json:"config"`
}

// APIResponse represents standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}