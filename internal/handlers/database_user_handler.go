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
func (h *DatabaseUserHandler) ChangePassword(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

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
func (h *DatabaseUserHandler) UpdateUserStatus(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

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