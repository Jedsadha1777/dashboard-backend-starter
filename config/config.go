package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBHost     string
	DBTimeZone string
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment")
	}

	DBUser = getEnv("DB_USER", "root")
	DBPassword = getEnv("DB_PASSWORD", "12345")
	DBName = getEnv("DB_NAME", "dashboard-starter")
	DBPort = getEnv("DB_PORT", "5432")
	DBHost = getEnv("DB_HOST", "localhost")
	DBTimeZone = getEnv("DB_TIMEZONE", "UTC")

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		DBHost, DBUser, DBPassword, DBName, DBPort, DBTimeZone)
}
