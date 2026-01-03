package http

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/johnroshan2255/auth-service/internal/middleware"
	"github.com/johnroshan2255/auth-service/internal/service"
)

func SetupRouter(authService *service.AuthService) *gin.Engine {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Configure trusted proxies - empty slice means no trusted proxies (safer default)
	// Set specific IPs if behind a reverse proxy (e.g., []string{"127.0.0.1", "::1"})
	router.SetTrustedProxies([]string{})

	// Add CORS middleware for frontend
	router.Use(middleware.CORSMiddleware())

	authHandler := NewAuthHandler(authService)

	api := router.Group("/api/v1")
	{
		api.GET("/health", authHandler.HealthCheck)
		
		auth := api.Group("/auth")
		{
			auth.POST("/signup", authHandler.Signup)
			auth.POST("/login", authHandler.Login)
			auth.POST("/validate", authHandler.ValidateToken)
			auth.GET("/me", middleware.AuthMiddleware(), authHandler.GetCurrentUser)
		}
	}

	return router
}

