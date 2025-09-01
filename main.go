package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nabilulilalbab/apivpn/internal/config"
	"github.com/nabilulilalbab/apivpn/internal/database"
	"github.com/nabilulilalbab/apivpn/internal/handlers"
	"github.com/nabilulilalbab/apivpn/internal/middleware"
	"github.com/nabilulilalbab/apivpn/internal/services"

	// Swagger imports
	_ "github.com/nabilulilalbab/apivpn/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title VPN API
// @version 1.0
// @description API for VPN management system supporting SSH, VMESS, VLESS, Trojan, and Shadowsocks protocols
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logging
	setupLogging(cfg.LogFile)

	// Initialize database
	if err := database.InitDB(cfg.DBPath); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Initialize services
	systemService := services.NewSystemService()
	vpnService := services.NewVPNService()

	// Use database user service instead of file-based
	userService := services.NewDatabaseUserService(
		cfg.JWTSecret,
		cfg.BCryptCost,
		cfg.TokenExpireHours,
		cfg.MaxLoginAttempts,
	)

	// Create default admin user if none exists
	if err := userService.CreateDefaultAdmin(); err != nil {
		log.Printf("Warning: Failed to create default admin: %v", err)
	}

	// Initialize handlers
	systemHandler := handlers.NewSystemHandler(systemService)
	vpnHandler := handlers.NewVPNHandler(vpnService, nil) // VPN handler doesn't need user service
	userHandler := handlers.NewDatabaseUserHandler(userService)

	// Setup Gin mode
	if cfg.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.RateLimit(cfg.RateLimitRequests, cfg.RateLimitWindow))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		dbStats := database.GetDBStats()
		c.JSON(200, gin.H{
			"status":   "ok",
			"database": dbStats,
			"version":  "1.0.0",
		})
	})

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		protected.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			// User management
			user := protected.Group("/user")
			{
				user.GET("/profile", userHandler.GetProfile)
				user.PUT("/password", userHandler.ChangePassword)
				user.GET("/list", userHandler.ListUsers)
				user.PUT("/:username/status", userHandler.UpdateUserStatus)
				user.DELETE("/:username", userHandler.DeleteUser)
			}

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

	// Log startup information
	log.Printf("🚀 VPN API Server starting...")
	log.Printf("📡 Host: %s", cfg.Host)
	log.Printf("🔌 Port: %s", cfg.Port)
	log.Printf("💾 Database: %s", cfg.DBPath)
	log.Printf("🔐 JWT Secret: %s", maskSecret(cfg.JWTSecret))
	log.Printf("📊 Log Level: %s", cfg.LogLevel)

	// Log database event
	database.LogEvent("info", fmt.Sprintf("VPN API Server started on %s:%s", cfg.Host, cfg.Port), "system", "", "")

	// Start server
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if err := router.Run(address); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupLogging(logFile string) {
	// Create log directory if not exists
	if logFile != "" {
		dir := filepath.Dir(logFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Failed to create log directory: %v", err)
		}
	}
}

func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}
