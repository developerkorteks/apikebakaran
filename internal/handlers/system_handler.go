package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/apivpn/internal/models"
	"github.com/nabilulilalbab/apivpn/internal/services"
)

type SystemHandler struct {
	systemService *services.SystemService
}

func NewSystemHandler(systemService *services.SystemService) *SystemHandler {
	return &SystemHandler{
		systemService: systemService,
	}
}

// GetSystemInfo returns system information
// @Summary Get system information
// @Description Get detailed system information including CPU, RAM, disk usage, and network info
// @Tags System
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.SystemInfo} "System information retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get system information"
// @Router /system/info [get]
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	info, err := h.systemService.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get system information: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "System information retrieved successfully",
		Data:    info,
	})
}

// GetServiceStatus returns VPN service status
// @Summary Get VPN service status
// @Description Get the status of all managed VPN services
// @Tags System
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=map[string]string} "Service status retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get service status"
// @Router /system/status [get]
func (h *SystemHandler) GetServiceStatus(c *gin.Context) {
	status, err := h.systemService.GetServiceStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get service status: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Service status retrieved successfully",
		Data:    status,
	})
}

// GetBandwidthUsage returns bandwidth usage information
// @Summary Get bandwidth usage
// @Description Get daily and monthly bandwidth usage
// @Tags System
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=map[string]string} "Bandwidth usage retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get bandwidth usage"
// @Router /system/bandwidth [get]
func (h *SystemHandler) GetBandwidthUsage(c *gin.Context) {
	info, err := h.systemService.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get bandwidth usage: " + err.Error(),
		})
		return
	}

	bandwidthData := map[string]string{
		"daily":   info.DailyBandwidth,
		"monthly": info.MonthlyBandwidth,
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Bandwidth usage retrieved successfully",
		Data:    bandwidthData,
	})
}

// AddDomain adds a new domain to the system
// @Summary Add a new domain
// @Description Add a new domain to the system for VPN services
// @Tags Domain
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AddDomainRequest true "Domain name"
// @Success 200 {object} models.APIResponse "Domain added successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to add domain"
// @Router /domain/add [post]
func (h *SystemHandler) AddDomain(c *gin.Context) {
	var req models.AddDomainRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.systemService.AddDomain(req.Domain); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to add domain: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Domain added successfully. Don't forget to renew SSL certificate.",
		Data: map[string]string{
			"domain": req.Domain,
		},
	})
}

// GetCurrentDomain returns the current domain
// @Summary Get current domain
// @Description Get the currently configured domain and server IP
// @Tags Domain
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=map[string]string} "Current domain retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get current domain"
// @Router /domain/current [get]
func (h *SystemHandler) GetCurrentDomain(c *gin.Context) {
	info, err := h.systemService.GetSystemInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get current domain: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Current domain retrieved successfully",
		Data: map[string]string{
			"domain": info.Domain,
			"ip":     info.IP,
		},
	})
}

// RenewSSL renews SSL certificate
// @Summary Renew SSL certificate
// @Description Renew the SSL certificate for the current domain
// @Tags Domain
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "SSL certificate renewed successfully"
// @Failure 500 {object} models.APIResponse "Failed to renew SSL certificate"
// @Router /domain/ssl/renew [post]
func (h *SystemHandler) RenewSSL(c *gin.Context) {
	if err := h.systemService.RenewSSL(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to renew SSL certificate: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SSL certificate renewed successfully",
	})
}

// Reboot system
// @Summary Reboot system
// @Description Reboot the entire system
// @Tags System
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "System reboot initiated"
// @Failure 500 {object} models.APIResponse "Failed to reboot system"
// @Router /system/reboot [post]
func (h *SystemHandler) Reboot(c *gin.Context) {
	if err := h.systemService.Reboot(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to reboot system: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "System reboot initiated",
	})
}

// RestartServices restarts VPN services
// @Summary Restart VPN services
// @Description Restart all managed VPN services
// @Tags System
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "VPN services restarted successfully"
// @Failure 500 {object} models.APIResponse "Failed to restart services"
// @Router /system/restart [post]
func (h *SystemHandler) RestartServices(c *gin.Context) {
	if err := h.systemService.RestartServices(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to restart services: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "VPN services restarted successfully",
	})
}