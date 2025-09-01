package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/apivpn/internal/models"
	"github.com/nabilulilalbab/apivpn/internal/services"
)

type VPNHandler struct {
	vpnService  *services.VPNService
	userService *services.UserService
}

func NewVPNHandler(vpnService *services.VPNService, userService *services.UserService) *VPNHandler {
	return &VPNHandler{
		vpnService:  vpnService,
		userService: userService,
	}
}

// SSH User Management
func (h *VPNHandler) CreateSSHUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	req.Protocol = "ssh"
	config, err := h.vpnService.CreateSSHUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create SSH user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "SSH user created successfully",
		Data:    config,
	})
}

func (h *VPNHandler) GetSSHUsers(c *gin.Context) {
	users, err := h.vpnService.GetSSHUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get SSH users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SSH users retrieved successfully",
		Data:    users,
	})
}

func (h *VPNHandler) DeleteSSHUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
		return
	}

	if err := h.vpnService.DeleteSSHUser(username); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete SSH user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SSH user deleted successfully",
	})
}

func (h *VPNHandler) ExtendSSHUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
		return
	}

	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.vpnService.ExtendUser("ssh", username, req.Days); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to extend SSH user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "SSH user extended successfully",
	})
}

// VMESS User Management
func (h *VPNHandler) CreateVmessUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	req.Protocol = "vmess"
	config, err := h.vpnService.CreateVmessUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create VMESS user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "VMESS user created successfully",
		Data:    config,
	})
}

func (h *VPNHandler) GetVmessUsers(c *gin.Context) {
	users, err := h.vpnService.GetVmessUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get VMESS users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "VMESS users retrieved successfully",
		Data:    users,
	})
}

func (h *VPNHandler) DeleteVmessUser(c *gin.Context) {
	username := c.Param("username")
	if err := h.vpnService.DeleteVmessUser(username); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete VMESS user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "VMESS user deleted successfully",
	})
}

func (h *VPNHandler) ExtendVmessUser(c *gin.Context) {
	username := c.Param("username")
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.vpnService.ExtendUser("vmess", username, req.Days); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to extend VMESS user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "VMESS user extended successfully",
	})
}

// VLESS User Management
func (h *VPNHandler) CreateVlessUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	req.Protocol = "vless"
	config, err := h.vpnService.CreateVlessUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create VLESS user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "VLESS user created successfully",
		Data:    config,
	})
}

func (h *VPNHandler) GetVlessUsers(c *gin.Context) {
	users, err := h.vpnService.GetVlessUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get VLESS users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "VLESS users retrieved successfully",
		Data:    users,
	})
}

func (h *VPNHandler) DeleteVlessUser(c *gin.Context) {
	username := c.Param("username")
	if err := h.vpnService.DeleteVlessUser(username); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete VLESS user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "VLESS user deleted successfully",
	})
}

func (h *VPNHandler) ExtendVlessUser(c *gin.Context) {
	username := c.Param("username")
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.vpnService.ExtendUser("vless", username, req.Days); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to extend VLESS user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "VLESS user extended successfully",
	})
}

// Trojan User Management
func (h *VPNHandler) CreateTrojanUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	req.Protocol = "trojan"
	config, err := h.vpnService.CreateTrojanUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create Trojan user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Trojan user created successfully",
		Data:    config,
	})
}

func (h *VPNHandler) GetTrojanUsers(c *gin.Context) {
	users, err := h.vpnService.GetTrojanUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get Trojan users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Trojan users retrieved successfully",
		Data:    users,
	})
}

func (h *VPNHandler) DeleteTrojanUser(c *gin.Context) {
	username := c.Param("username")
	if err := h.vpnService.DeleteTrojanUser(username); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete Trojan user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Trojan user deleted successfully",
	})
}

func (h *VPNHandler) ExtendTrojanUser(c *gin.Context) {
	username := c.Param("username")
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.vpnService.ExtendUser("trojan", username, req.Days); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to extend Trojan user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Trojan user extended successfully",
	})
}

// Shadowsocks User Management
func (h *VPNHandler) CreateShadowsocksUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	req.Protocol = "shadowsocks"
	config, err := h.vpnService.CreateShadowsocksUser(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create Shadowsocks user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Shadowsocks user created successfully",
		Data:    config,
	})
}

func (h *VPNHandler) GetShadowsocksUsers(c *gin.Context) {
	users, err := h.vpnService.GetShadowsocksUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get Shadowsocks users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Shadowsocks users retrieved successfully",
		Data:    users,
	})
}

func (h *VPNHandler) DeleteShadowsocksUser(c *gin.Context) {
	username := c.Param("username")
	if err := h.vpnService.DeleteShadowsocksUser(username); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete Shadowsocks user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Shadowsocks user deleted successfully",
	})
}

func (h *VPNHandler) ExtendShadowsocksUser(c *gin.Context) {
	username := c.Param("username")
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
		})
		return
	}

	if err := h.vpnService.ExtendUser("shadowsocks", username, req.Days); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to extend Shadowsocks user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Shadowsocks user extended successfully",
	})
}

// General VPN operations
func (h *VPNHandler) GetAllUsers(c *gin.Context) {
	allUsers := make(map[string][]models.User)

	if sshUsers, err := h.vpnService.GetSSHUsers(); err == nil {
		allUsers["ssh"] = sshUsers
	}

	if vmessUsers, err := h.vpnService.GetVmessUsers(); err == nil {
		allUsers["vmess"] = vmessUsers
	}

	if vlessUsers, err := h.vpnService.GetVlessUsers(); err == nil {
		allUsers["vless"] = vlessUsers
	}

	if trojanUsers, err := h.vpnService.GetTrojanUsers(); err == nil {
		allUsers["trojan"] = trojanUsers
	}

	if ssUsers, err := h.vpnService.GetShadowsocksUsers(); err == nil {
		allUsers["shadowsocks"] = ssUsers
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "All users retrieved successfully",
		Data:    allUsers,
	})
}

func (h *VPNHandler) GetUserTraffic(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
		return
	}

	traffic, err := h.vpnService.GetUserTraffic(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to get user traffic: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User traffic retrieved successfully",
		Data:    traffic,
	})
}

func (h *VPNHandler) CleanupExpiredUsers(c *gin.Context) {
	if err := h.vpnService.CleanupExpiredUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to cleanup expired users: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Expired users cleaned up successfully",
	})
}