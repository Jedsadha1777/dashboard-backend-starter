package routes

import (
	"dashboard-starter/config"
	"dashboard-starter/controllers"
	"dashboard-starter/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all application routes
func SetupRouter() *gin.Engine {
	// Initialize rate limiter now that configuration is loaded
	middleware.InitRateLimiter()

	r := gin.Default()

	// ตั้งค่า trusted proxies
	trustedProxies := config.Config.Server.TrustedProxies
	if len(trustedProxies) == 0 {
		r.SetTrustedProxies(nil)
		log.Println("Warning: Not using trusted proxies. All requests will be trusted.")
	} else {
		r.SetTrustedProxies(trustedProxies)
		log.Println("Using trusted proxies:", trustedProxies)
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

	// Admin Auth routes
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

	// User Auth routes - new endpoints for user registration and login
	userAuth := v1.Group("/user/auth")
	{
		userAuth.POST("/register", controllers.UserRegister)
		userAuth.POST("/login", controllers.UserLogin)
		userAuth.POST("/refresh", controllers.RefreshToken) // Reuse the same token refresh endpoint

		// Protected routes for users
		userProtected := userAuth.Group("")
		userProtected.Use(middleware.AuthMiddleware(), middleware.UserRequired())
		{
			userProtected.POST("/logout", controllers.UserLogout)
			userProtected.GET("/profile", controllers.GetUserProfile)
			userProtected.POST("/change-password", controllers.ChangeUserPassword)
		}
	}

	// User routes - endpoints for logged in users
	user := v1.Group("/user")
	user.Use(middleware.AuthMiddleware(), middleware.UserRequired())
	{
		// User profile and dashboard
		user.GET("/dashboard", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(200, gin.H{
				"success": true,
				"data": gin.H{
					"message": "Welcome to User Dashboard",
					"id":      userID,
				},
			})
		})

		// Update own profile
		user.PUT("/profile", controllers.UpdateUserProfile)
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

		// user management routes
		users := admin.Group("/users")
		{
			users.GET("", controllers.ListUsers)
			users.POST("", controllers.CreateUser)
			users.GET("/:id", controllers.GetUser)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
			users.POST("/:id/reset-password", controllers.ResetUserPassword)
		}

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

	// Public API endpoints - accessible without authentication
	public := v1.Group("/public")
	{
		// Example: public article listing
		public.GET("/articles", func(c *gin.Context) {
			// Implement a controller for public articles
			c.JSON(200, gin.H{
				"success": true,
				"data": gin.H{
					"message": "Public articles will be listed here",
				},
			})
		})
	}

	// ใช้เพื่อการ debug ให้แสดง registerd routes ทั้งหมด
	debugRoutes(r)

	return r
}

// debugRoutes แสดง registered routes ทั้งหมดเพื่อช่วยในการ debug
func debugRoutes(r *gin.Engine) {
	routes := r.Routes()
	log.Println("Registered routes:")
	for _, route := range routes {
		log.Printf("[%s] %s", route.Method, route.Path)
	}
}
