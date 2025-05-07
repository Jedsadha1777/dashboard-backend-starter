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
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
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

		// Additional admin routes can be added here
		devices := admin.Group("/devices")
		{
			devices.POST("", controllers.CreateDevice)                    // สร้างอุปกรณ์ใหม่
			devices.GET("", controllers.ListDevices)                      // ดูรายการอุปกรณ์
			devices.GET("/:id", controllers.GetDevice)                    // ดูข้อมูลอุปกรณ์
			devices.PUT("/:id", controllers.UpdateDevice)                 // อัปเดตข้อมูลอุปกรณ์
			devices.DELETE("/:id", controllers.DeleteDevice)              // ลบอุปกรณ์
			devices.POST("/:id/reset-key", controllers.ResetDeviceApiKey) // รีเซ็ต API key
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
