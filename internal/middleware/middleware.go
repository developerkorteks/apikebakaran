package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nabilulilalbab/apivpn/internal/database"
	"github.com/nabilulilalbab/apivpn/internal/models"
)

// Rate limiting
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, windowSeconds int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   time.Duration(windowSeconds) * time.Second,
	}
}

func (rl *RateLimiter) Allow(clientIP string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean old requests
	if requests, exists := rl.requests[clientIP]; exists {
		var validRequests []time.Time
		for _, reqTime := range requests {
			if reqTime.After(cutoff) {
				validRequests = append(validRequests, reqTime)
			}
		}
		rl.requests[clientIP] = validRequests
	}

	// Check if limit exceeded
	if len(rl.requests[clientIP]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[clientIP] = append(rl.requests[clientIP], now)
	return true
}

var globalRateLimiter *RateLimiter

// CORS middleware
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// Logger middleware
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Log to database for important events
		if param.StatusCode >= 400 {
			level := "warning"
			if param.StatusCode >= 500 {
				level = "error"
			}
			
			message := fmt.Sprintf("%s %s - Status: %d", param.Method, param.Path, param.StatusCode)
			database.LogEvent(level, message, "api", "", param.ClientIP)
		}

		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// RateLimit middleware
func RateLimit(limit int, windowSeconds int) gin.HandlerFunc {
	globalRateLimiter = NewRateLimiter(limit, windowSeconds)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		if !globalRateLimiter.Allow(clientIP) {
			database.LogEvent("warning", fmt.Sprintf("Rate limit exceeded for IP: %s", clientIP), "api", "", clientIP)
			
			c.JSON(http.StatusTooManyRequests, models.APIResponse{
				Success: false,
				Error:   "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// AuthRequired middleware
func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Authorization header required",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			database.LogEvent("warning", fmt.Sprintf("Invalid token attempt from IP: %s", c.ClientIP()), "auth", "", c.ClientIP())
			
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Error:   "Invalid token",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Check token expiration
			if exp, ok := claims["exp"].(float64); ok {
				if time.Now().Unix() > int64(exp) {
					c.JSON(http.StatusUnauthorized, models.APIResponse{
						Success: false,
						Error:   "Token expired",
					})
					c.Abort()
					return
				}
			}

			c.Set("username", claims["username"])
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}