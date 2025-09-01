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
func (h *SystemHandler) AddDomain(c *gin.Context) {
	var req struct {
		Domain string `json:"domain" binding:"required"`
	}

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