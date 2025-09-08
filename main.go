package main

import (
	"context"
	"dashboard-starter/config"
	"dashboard-starter/db"
	"dashboard-starter/internal/infrastructure/logger"
	"dashboard-starter/routes"
	"dashboard-starter/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dashboard-starter/internal/infrastructure/cache"
	"dashboard-starter/internal/infrastructure/email"

	"github.com/gin-gonic/gin"
)

func init() {
	time.Local = time.UTC
}

// @title Dashboard API
// @version 1.0
// @description Admin Dashboard Backend API with JWT Authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize logger first
	_, err := logger.NewZapLogger("development")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	log.Println("Starting application...")

	// Load configuration
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)

	}

	// Set Gin mode based on environment
	if config.Config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Redis cache (optional)
	cacheConfig := cache.Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	redisCache, err := cache.NewRedisCache(cacheConfig)
	if err != nil {
		log.Printf("Starting without Redis cache: %v", err)
	} else {
		log.Println("Redis cache initialized successfully")
	}

	// Initialize Email service
	emailConfig := email.Config{
		SMTPHost:     "", // Empty = mock mode
		SMTPPort:     "587",
		SMTPUser:     "",
		SMTPPassword: "",
		FromEmail:    "noreply@yourdomain.com",
		FromName:     "Dashboard Team",
	}

	emailService := email.NewSMTPEmailService(emailConfig)
	log.Println("Email service initialized")

	// Initialize validation
	if err := utils.InitValidator(); err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
	}

	utils.InitPasswordConfig(config.Config.Security.MinPasswordLength)

	// Initialize JWT
	if err := utils.InitJWT(); err != nil {
		log.Fatalf("Failed to initialize JWT: %v", err)
	}

	// Initialize database
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Ensure database is closed when application exits
	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()

	// Seed admin user only in development or when explicitly enabled
	shouldSeed := false
	if config.Config.Environment == "development" {
		shouldSeed = true
	} else if envSeed := os.Getenv("AUTO_SEED"); envSeed == "true" {
		shouldSeed = true
	}

	if shouldSeed {
		if err := db.SeedAdmin(); err != nil {
			log.Printf("Warning: Failed to seed admin user: %v", err)
		}
	}

	// Setup HTTP router with new architecture
	router := routes.SetupRouter()

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.Config.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Config.Server.ReadTimeout,
		WriteTimeout: config.Config.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port: %s", config.Config.Server.Port)
		log.Printf("Swagger UI available at: http://localhost:%s/swagger/index.html", config.Config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")

	// Log usage information for unused services (for future expansion)
	_ = redisCache   // Redis cache is initialized but not used in handlers yet
	_ = emailService // Email service is initialized but not used in handlers yet
}
