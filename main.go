package main

import (
	"context"
	"dashboard-starter/config"
	"dashboard-starter/db"
	"dashboard-starter/routes"
	"dashboard-starter/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	// Set default timezone to UTC
	time.Local = time.UTC

	// Initialize logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	log.Println("Starting application...")

	// Load configuration
	if err := config.Init(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

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

	// Seed admin user
	if err := db.SeedAdmin(); err != nil {
		log.Fatalf("Failed to seed admin user: %v", err)
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
		log.Printf("Server starting on port %s", config.Config.Server.Port)
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
}
