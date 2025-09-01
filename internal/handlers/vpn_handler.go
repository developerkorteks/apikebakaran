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
// @Summary Create SSH user
// @Description Create a new SSH/WebSocket VPN user
// @Tags VPN - SSH
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateUserRequest true "SSH user creation request"
// @Success 201 {object} models.APIResponse{data=models.VPNConfig} "SSH user created successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 401 {object} models.APIResponse "Unauthorized"
// @Failure 500 {object} models.APIResponse "Failed to create SSH user"
// @Router /vpn/ssh/create [post]
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

// @Summary Get SSH users
// @Description Get a list of all SSH/WebSocket VPN users
// @Tags VPN - SSH
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.User} "SSH users retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get SSH users"
// @Router /vpn/ssh/users [get]
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

// @Summary Delete SSH user
// @Description Delete an SSH/WebSocket VPN user
// @Tags VPN - SSH
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Success 200 {object} models.APIResponse "SSH user deleted successfully"
// @Failure 400 {object} models.APIResponse "Username is required"
// @Failure 500 {object} models.APIResponse "Failed to delete SSH user"
// @Router /vpn/ssh/users/{username} [delete]
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

// @Summary Extend SSH user
// @Description Extend the expiration date of an SSH/WebSocket VPN user
// @Tags VPN - SSH
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param request body models.ExtendUserRequest true "Number of days to extend"
// @Success 200 {object} models.APIResponse "SSH user extended successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to extend SSH user"
// @Router /vpn/ssh/users/{username}/extend [put]
func (h *VPNHandler) ExtendSSHUser(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Username is required",
		})
		return
	}

	var req models.ExtendUserRequest

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

// @Summary Create VMESS user
// @Description Create a new VMESS VPN user
// @Tags VPN - VMESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateUserRequest true "VMESS user creation request"
// @Success 201 {object} models.APIResponse{data=models.VPNConfig} "VMESS user created successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to create VMESS user"
// @Router /vpn/vmess/create [post]
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

// @Summary Get VMESS users
// @Description Get a list of all VMESS VPN users
// @Tags VPN - VMESS
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.User} "VMESS users retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get VMESS users"
// @Router /vpn/vmess/users [get]
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

// @Summary Delete VMESS user
// @Description Delete a VMESS VPN user
// @Tags VPN - VMESS
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Success 200 {object} models.APIResponse "VMESS user deleted successfully"
// @Failure 400 {object} models.APIResponse "Username is required"
// @Failure 500 {object} models.APIResponse "Failed to delete VMESS user"
// @Router /vpn/vmess/users/{username} [delete]
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

// @Summary Extend VMESS user
// @Description Extend the expiration date of a VMESS VPN user
// @Tags VPN - VMESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param request body models.ExtendUserRequest true "Number of days to extend"
// @Success 200 {object} models.APIResponse "VMESS user extended successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to extend VMESS user"
// @Router /vpn/vmess/users/{username}/extend [put]
func (h *VPNHandler) ExtendVmessUser(c *gin.Context) {
	username := c.Param("username")
	var req models.ExtendUserRequest

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

// @Summary Create VLESS user
// @Description Create a new VLESS VPN user
// @Tags VPN - VLESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateUserRequest true "VLESS user creation request"
// @Success 201 {object} models.APIResponse{data=models.VPNConfig} "VLESS user created successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to create VLESS user"
// @Router /vpn/vless/create [post]
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

// @Summary Get VLESS users
// @Description Get a list of all VLESS VPN users
// @Tags VPN - VLESS
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.User} "VLESS users retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get VLESS users"
// @Router /vpn/vless/users [get]
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

// @Summary Delete VLESS user
// @Description Delete a VLESS VPN user
// @Tags VPN - VLESS
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Success 200 {object} models.APIResponse "VLESS user deleted successfully"
// @Failure 400 {object} models.APIResponse "Username is required"
// @Failure 500 {object} models.APIResponse "Failed to delete VLESS user"
// @Router /vpn/vless/users/{username} [delete]
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

// @Summary Extend VLESS user
// @Description Extend the expiration date of a VLESS VPN user
// @Tags VPN - VLESS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param request body models.ExtendUserRequest true "Number of days to extend"
// @Success 200 {object} models.APIResponse "VLESS user extended successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to extend VLESS user"
// @Router /vpn/vless/users/{username}/extend [put]
func (h *VPNHandler) ExtendVlessUser(c *gin.Context) {
	username := c.Param("username")
	var req models.ExtendUserRequest

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

// @Summary Create Trojan user
// @Description Create a new Trojan VPN user
// @Tags VPN - Trojan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateUserRequest true "Trojan user creation request"
// @Success 201 {object} models.APIResponse{data=models.VPNConfig} "Trojan user created successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to create Trojan user"
// @Router /vpn/trojan/create [post]
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

// @Summary Get Trojan users
// @Description Get a list of all Trojan VPN users
// @Tags VPN - Trojan
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.User} "Trojan users retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get Trojan users"
// @Router /vpn/trojan/users [get]
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

// @Summary Delete Trojan user
// @Description Delete a Trojan VPN user
// @Tags VPN - Trojan
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Success 200 {object} models.APIResponse "Trojan user deleted successfully"
// @Failure 400 {object} models.APIResponse "Username is required"
// @Failure 500 {object} models.APIResponse "Failed to delete Trojan user"
// @Router /vpn/trojan/users/{username} [delete]
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

// @Summary Extend Trojan user
// @Description Extend the expiration date of a Trojan VPN user
// @Tags VPN - Trojan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param request body models.ExtendUserRequest true "Number of days to extend"
// @Success 200 {object} models.APIResponse "Trojan user extended successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to extend Trojan user"
// @Router /vpn/trojan/users/{username}/extend [put]
func (h *VPNHandler) ExtendTrojanUser(c *gin.Context) {
	username := c.Param("username")
	var req models.ExtendUserRequest

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

// @Summary Create Shadowsocks user
// @Description Create a new Shadowsocks VPN user
// @Tags VPN - Shadowsocks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateUserRequest true "Shadowsocks user creation request"
// @Success 201 {object} models.APIResponse{data=models.VPNConfig} "Shadowsocks user created successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to create Shadowsocks user"
// @Router /vpn/shadowsocks/create [post]
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

// @Summary Get Shadowsocks users
// @Description Get a list of all Shadowsocks VPN users
// @Tags VPN - Shadowsocks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.User} "Shadowsocks users retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to get Shadowsocks users"
// @Router /vpn/shadowsocks/users [get]
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

// @Summary Delete Shadowsocks user
// @Description Delete a Shadowsocks VPN user
// @Tags VPN - Shadowsocks
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Success 200 {object} models.APIResponse "Shadowsocks user deleted successfully"
// @Failure 400 {object} models.APIResponse "Username is required"
// @Failure 500 {object} models.APIResponse "Failed to delete Shadowsocks user"
// @Router /vpn/shadowsocks/users/{username} [delete]
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

// @Summary Extend Shadowsocks user
// @Description Extend the expiration date of a Shadowsocks VPN user
// @Tags VPN - Shadowsocks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Param request body models.ExtendUserRequest true "Number of days to extend"
// @Success 200 {object} models.APIResponse "Shadowsocks user extended successfully"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 500 {object} models.APIResponse "Failed to extend Shadowsocks user"
// @Router /vpn/shadowsocks/users/{username}/extend [put]
func (h *VPNHandler) ExtendShadowsocksUser(c *gin.Context) {
	username := c.Param("username")
	var req models.ExtendUserRequest

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

// @Summary Get all VPN users
// @Description Get a list of all users for all VPN protocols
// @Tags VPN - General
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=map[string][]models.User} "All users retrieved successfully"
// @Router /vpn/users/all [get]
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

// @Summary Get user traffic
// @Description Get traffic usage for a specific user
// @Tags VPN - General
// @Produce json
// @Security BearerAuth
// @Param username path string true "Username"
// @Success 200 {object} models.APIResponse{data=map[string]interface{}} "User traffic retrieved successfully"
// @Failure 400 {object} models.APIResponse "Username is required"
// @Failure 500 {object} models.APIResponse "Failed to get user traffic"
// @Router /vpn/users/{username}/traffic [get]
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

// @Summary Cleanup expired users
// @Description Remove all expired VPN users from the system
// @Tags VPN - General
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Expired users cleaned up successfully"
// @Failure 500 {object} models.APIResponse "Failed to cleanup expired users"
// @Router /vpn/users/cleanup-expired [post]
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
