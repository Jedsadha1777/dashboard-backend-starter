package db

import (
	"context"
	"dashboard-starter/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

// Init initializes the database connection
func Init() error {
	var err error

	// Custom logger configuration for GORM
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	// Database connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Open connection
	DB, err = gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	// Check connection
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// Test connection with context timeout
	if err := sqlDB.PingContext(ctx); err != nil {
		return err
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(config.Config.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Config.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Config.Database.MaxLifetime) * time.Minute)

	// Run migrations
	if err := runMigrations(); err != nil {
		return err
	}

	log.Println("Database connection established successfully")
	return nil
}

// runMigrations runs database migrations
func runMigrations() error {
	log.Println("Running database migrations...")

	// ลงทะเบียนโมเดลทั้งหมด
	RegisterAllModels()

	// ทำ migration
	if err := DB.AutoMigrate(GetAllModels()...); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Transaction runs a function within a database transaction
func Transaction(fn func(tx *gorm.DB) error) error {
	tx := DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
