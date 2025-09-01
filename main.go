package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/apivpn/internal/config"
	"github.com/nabilulilalbab/apivpn/internal/handlers"
	"github.com/nabilulilalbab/apivpn/internal/middleware"
	"github.com/nabilulilalbab/apivpn/internal/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize services
	systemService := services.NewSystemService()
	vpnService := services.NewVPNService()
	userService := services.NewUserService()

	// Initialize handlers
	systemHandler := handlers.NewSystemHandler(systemService)
	vpnHandler := handlers.NewVPNHandler(vpnService, userService)
	userHandler := handlers.NewUserHandler(userService)

	// Setup router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())

	// API routes
	api := router.Group("/api/v1")
	{
		// Authentication
		auth := api.Group("/auth")
		{
			auth.POST("/login", userHandler.Login)
			auth.POST("/register", userHandler.Register)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// System monitoring
			system := protected.Group("/system")
			{
				system.GET("/info", systemHandler.GetSystemInfo)
				system.GET("/status", systemHandler.GetServiceStatus)
				system.GET("/bandwidth", systemHandler.GetBandwidthUsage)
				system.POST("/reboot", systemHandler.Reboot)
				system.POST("/restart", systemHandler.RestartServices)
			}

			// VPN Management
			vpn := protected.Group("/vpn")
			{
				// SSH/WebSocket
				ssh := vpn.Group("/ssh")
				{
					ssh.POST("/create", vpnHandler.CreateSSHUser)
					ssh.GET("/users", vpnHandler.GetSSHUsers)
					ssh.DELETE("/users/:username", vpnHandler.DeleteSSHUser)
					ssh.PUT("/users/:username/extend", vpnHandler.ExtendSSHUser)
				}

				// VMESS
				vmess := vpn.Group("/vmess")
				{
					vmess.POST("/create", vpnHandler.CreateVmessUser)
					vmess.GET("/users", vpnHandler.GetVmessUsers)
					vmess.DELETE("/users/:username", vpnHandler.DeleteVmessUser)
					vmess.PUT("/users/:username/extend", vpnHandler.ExtendVmessUser)
				}

				// VLESS
				vless := vpn.Group("/vless")
				{
					vless.POST("/create", vpnHandler.CreateVlessUser)
					vless.GET("/users", vpnHandler.GetVlessUsers)
					vless.DELETE("/users/:username", vpnHandler.DeleteVlessUser)
					vless.PUT("/users/:username/extend", vpnHandler.ExtendVlessUser)
				}

				// Trojan
				trojan := vpn.Group("/trojan")
				{
					trojan.POST("/create", vpnHandler.CreateTrojanUser)
					trojan.GET("/users", vpnHandler.GetTrojanUsers)
					trojan.DELETE("/users/:username", vpnHandler.DeleteTrojanUser)
					trojan.PUT("/users/:username/extend", vpnHandler.ExtendTrojanUser)
				}

				// Shadowsocks
				ss := vpn.Group("/shadowsocks")
				{
					ss.POST("/create", vpnHandler.CreateShadowsocksUser)
					ss.GET("/users", vpnHandler.GetShadowsocksUsers)
					ss.DELETE("/users/:username", vpnHandler.DeleteShadowsocksUser)
					ss.PUT("/users/:username/extend", vpnHandler.ExtendShadowsocksUser)
				}

				// General VPN operations
				vpn.GET("/users/all", vpnHandler.GetAllUsers)
				vpn.GET("/users/:username/traffic", vpnHandler.GetUserTraffic)
				vpn.POST("/users/cleanup-expired", vpnHandler.CleanupExpiredUsers)
			}

			// Domain and SSL management
			domain := protected.Group("/domain")
			{
				domain.POST("/add", systemHandler.AddDomain)
				domain.POST("/ssl/renew", systemHandler.RenewSSL)
				domain.GET("/current", systemHandler.GetCurrentDomain)
			}
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}