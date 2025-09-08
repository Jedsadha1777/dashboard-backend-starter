package ports

import "time"

// EmailService defines email sending capabilities
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendWelcomeEmail(to, name string) error
	SendPasswordResetEmail(to, resetToken string) error
}

// CacheService defines caching capabilities
type CacheService interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string, dest interface{}) error
	Delete(key string) error
	Exists(key string) bool
}

// LoggerService defines logging capabilities
type LoggerService interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
}
