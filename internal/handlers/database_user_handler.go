package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/apivpn/internal/models"
	"github.com/nabilulilalbab/apivpn/internal/services"
)

type DatabaseUserHandler struct {
	userService *services.DatabaseUserService
}

func NewDatabaseUserHandler(userService *services.DatabaseUserService) *DatabaseUserHandler {
	return &DatabaseUserHandler{
		userService: userService,
	}
}

// Login authenticates admin user
// @Summary Admin login
// @Description Authenticate admin user and get JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.APIResponse{data=models.LoginResponse} "Login successful"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 401 {object} models.APIResponse "Authentication failed"
// @Router /auth/login [post]
func (h *DatabaseUserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	response, err := h.userService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Register creates a new admin user
// @Summary Register new admin user
// @Description Create a new admin user account
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Registration credentials"
// @Success 201 {object} models.APIResponse "User registered successfully"
// @Failure 400 {object} models.APIResponse "Invalid request or user already exists"
// @Router /auth/register [post]
func (h *DatabaseUserHandler) Register(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.userService.Register(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User registered successfully",
	})
}

// GetProfile returns user profile information
// @Summary Get user profile
// @Description Get authenticated user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Profile retrieved successfully"
// @Failure 401 {object} models.APIResponse "User not authenticated"
// @Failure 404 {object} models.APIResponse "User not found"
// @Router /user/profile [get]
func (h *DatabaseUserHandler) GetProfile(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	userInfo, err := h.userService.GetUserInfo(username.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    userInfo,
	})
}

// ChangePassword changes user password
// @Summary Change user password
// @Description Change authenticated user's password
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChangePasswordRequest true "Old and new passwords"
// @Success 200 {object} models.APIResponse "Password changed successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 401 {object} models.APIResponse "User not authenticated"
// @Router /user/password [put]
func (h *DatabaseUserHandler) ChangePassword(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	var req models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.userService.ChangePassword(username.(string), req.OldPassword, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// ListUsers returns all admin users
func (h *DatabaseUserHandler) ListUsers(c *gin.Context) {
	users, err := h.userService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to list users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

// UpdateUserStatus updates user active status
// @Summary Update user status
// @Description Enable or disable a user account
// @Tags User Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param request body models.UpdateUserStatusRequest true "User status"
// @Success 200 {object} models.APIResponse "User status updated successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to update user status"
// @Router /user/{username}/status [put]
func (h *DatabaseUserHandler) UpdateUserStatus(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
		return
	}

	var req models.UpdateUserStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.userService.UpdateUserStatus(username, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update user status: " + err.Error(),
		})
		return
	}

	status := "disabled"
	if req.IsActive {
		status = "enabled"
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: fmt.Sprintf("User %s successfully", status),
	})
}

// DeleteUser deletes a user
func (h *DatabaseUserHandler) DeleteUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
		return
	}

	if err := h.userService.DeleteUser(username); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}
