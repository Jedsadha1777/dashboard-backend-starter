package main

import (
	"context"
	"dashboard-starter/config"
	"dashboard-starter/db"
	"dashboard-starter/pkg/logger"
	"dashboard-starter/routes"
	"dashboard-starter/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"dashboard-starter/pkg/cache"
	"dashboard-starter/pkg/email"
)

func init() {
	// Set default timezone to UTC
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
	if err := logger.Init("development"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Logger.Sync()

	logger.Info("Starting application...")

	// Load configuration
	if err := config.Init(); err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize Redis cache (optional)
	cacheConfig := cache.Config{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DB:       0,
	}

	if err := cache.Init(cacheConfig); err != nil {
		logger.Warn("Starting without Redis cache", zap.Error(err))
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

	if err := email.Init(emailConfig); err != nil {
		logger.Warn("Email service initialization failed", zap.Error(err))
	}

	// Initialize validation
	if err := utils.InitValidator(); err != nil {
		logger.Fatal("Failed to initialize validator", zap.Error(err))
	}

	utils.InitPasswordConfig(config.Config.Security.MinPasswordLength)

	// Initialize JWT
	if err := utils.InitJWT(); err != nil {
		logger.Fatal("Failed to initialize JWT", zap.Error(err))
	}

	// Initialize database
	if err := db.Init(); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Ensure database is closed when application exits
	sqlDB, _ := db.DB.DB()
	defer sqlDB.Close()

	// Seed admin user
	if err := db.SeedAdmin(); err != nil {
		logger.Fatal("Failed to seed admin user", zap.Error(err))
	}

	// Setup HTTP router
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
		logger.Info("Server starting", zap.String("port", config.Config.Server.Port))
		logger.Info("Swagger UI available", zap.String("url", fmt.Sprintf("http://localhost:%s/swagger/index.html", config.Config.Server.Port)))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited properly")
}
