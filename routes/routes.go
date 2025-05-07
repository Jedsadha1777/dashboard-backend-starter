package routes

import (
	"dashboard-starter/config"
	"dashboard-starter/controllers"
	"dashboard-starter/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all application routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	trustedProxies := config.Config.Server.TrustedProxies
	if len(trustedProxies) == 0 {
		r.SetTrustedProxies(nil)
	} else {
		r.SetTrustedProxies(trustedProxies)
	}

	// Apply global middlewares
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimitMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API versioning
	v1 := r.Group("/api/v1")

	// User routes
	users := v1.Group("/users")
	{
		users.GET("", controllers.ListUsers)
		users.POST("", controllers.CreateUser)
		users.GET("/:id", controllers.GetUser)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
	}

	// Auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh", controllers.RefreshToken)
		auth.POST("/device", controllers.DeviceAuth)

		// Protected routes
		protected := auth.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", controllers.Logout)
			protected.GET("/profile", controllers.GetProfile)
		}
	}

	// Admin dashboard routes
	// First use AuthMiddleware to verify token, then AdminRequired to ensure user is admin
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminRequired())
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			adminID, _ := c.Get("admin_id")
			c.JSON(200, gin.H{
				"success": true,
				"data": gin.H{
					"message": "Welcome to Admin Dashboard",
					"id":      adminID,
				},
			})
		})

		// Device management
		devices := admin.Group("/devices")
		{
			devices.POST("", controllers.CreateDevice)
			devices.GET("", controllers.ListDevices)
			devices.GET("/:id", controllers.GetDevice)
			devices.PUT("/:id", controllers.UpdateDevice)
			devices.DELETE("/:id", controllers.DeleteDevice)
			devices.POST("/:id/reset-key", controllers.ResetDeviceApiKey)
		}

		// Article management routes
		articles := admin.Group("/articles")
		{
			articles.POST("", controllers.CreateArticle)
			articles.GET("", controllers.ListArticles)
			articles.GET("/:id", controllers.GetArticle)
			articles.PUT("/:id", controllers.UpdateArticle)
			articles.DELETE("/:id", controllers.DeleteArticle)
			articles.POST("/:id/publish", controllers.PublishArticle)
		}
	}

	return r
}
