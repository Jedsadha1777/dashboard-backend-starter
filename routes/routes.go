package routes

import (
	"dashboard-starter/config"
	"dashboard-starter/internal/interfaces/http"
	"dashboard-starter/internal/interfaces/http/middleware"
	"log"

	_ "dashboard-starter/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"dashboard-starter/pkg/metrics"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Initialize dependency injection container
	container := http.NewContainer()

	// Metrics
	r.Use(metrics.PrometheusMiddleware())
	r.GET("/metrics", metrics.Handler())

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

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API versioning
	v1 := r.Group("/api/v1")

	// Admin Auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/login", container.AuthHandler.Login)
		auth.POST("/refresh", container.AuthHandler.RefreshToken)
		auth.POST("/device", container.DeviceHandler.AuthenticateDevice)

		// Protected routes
		protected := auth.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.POST("/logout", container.AuthHandler.Logout)
			protected.GET("/profile", container.AuthHandler.GetProfile)

		}
	}

	// User Auth routes
	userAuth := v1.Group("/user/auth")
	{
		userAuth.POST("/register", container.UserHandler.Register)
		userAuth.POST("/login", container.UserHandler.Login)
		userAuth.POST("/refresh", container.AuthHandler.RefreshToken) // Reuse the same token refresh endpoint

		// Protected routes for users
		userProtected := userAuth.Group("")
		userProtected.Use(middleware.AuthMiddleware(), middleware.UserRequired())
		{
			userProtected.POST("/logout", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				c.JSON(200, gin.H{
					"success": true,
					"data": gin.H{
						"message": "User logged out successfully",
						"user_id": userID,
					},
				})
			})
			userProtected.GET("/profile", container.UserHandler.GetProfile)
			userProtected.POST("/change-password", container.UserHandler.ChangePassword)
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
		user.PUT("/profile", container.UserHandler.UpdateProfile)
	}

	// Admin dashboard routes
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

		// User management routes
		users := admin.Group("/users")
		{
			users.GET("", container.UserHandler.ListUsers)
			users.POST("", container.UserHandler.CreateUser)
			users.GET("/:id", container.UserHandler.GetUser)
			users.PUT("/:id", container.UserHandler.UpdateUser)
			users.DELETE("/:id", container.UserHandler.DeleteUser)
			users.POST("/:id/reset-password", container.UserHandler.ResetPassword)
		}

		// Device management
		devices := admin.Group("/devices")
		{
			devices.POST("", container.DeviceHandler.CreateDevice)
			devices.GET("", container.DeviceHandler.ListDevices)
			devices.GET("/:id", container.DeviceHandler.GetDevice)
			devices.PUT("/:id", container.DeviceHandler.UpdateDevice)
			devices.DELETE("/:id", container.DeviceHandler.DeleteDevice)
			devices.POST("/:id/reset-key", container.DeviceHandler.ResetAPIKey)
		}

		// Article management routes
		articles := admin.Group("/articles")
		{
			articles.POST("", container.ArticleHandler.CreateArticle)
			articles.GET("", container.ArticleHandler.ListArticles)
			articles.GET("/:id", container.ArticleHandler.GetArticle)
			articles.PUT("/:id", container.ArticleHandler.UpdateArticle)
			articles.DELETE("/:id", container.ArticleHandler.DeleteArticle)
			articles.POST("/:id/publish", container.ArticleHandler.PublishArticle)
		}
	}

	// Public API endpoints - accessible without authentication
	public := v1.Group("/public")
	{
		public.GET("/articles", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"data": gin.H{
					"message": "Public articles will be listed here",
				},
			})
		})
	}

	// Debug routes
	debugRoutes(r)

	return r
}

func debugRoutes(r *gin.Engine) {
	routes := r.Routes()
	log.Println("Registered routes:")
	for _, route := range routes {
		log.Printf("[%s] %s", route.Method, route.Path)
	}
}
