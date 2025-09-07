package email

import (
	"fmt"
	"net/smtp"

	"dashboard-starter/pkg/logger"

	"go.uber.org/zap"
)

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type EmailService struct {
	config Config
}

var Service *EmailService

func Init(config Config) error {
	Service = &EmailService{config: config}

	// Test connection (optional)
	if config.SMTPHost != "" {
		logger.Info("Email service initialized", zap.String("smtp_host", config.SMTPHost))
	} else {
		logger.Warn("Email service running in mock mode (no SMTP configured)")
	}

	return nil
}

func (e *EmailService) SendEmail(to, subject, body string) error {
	if e.config.SMTPHost == "" {
		logger.Info("Mock email sent",
			zap.String("to", to),
			zap.String("subject", subject),
		)
		return nil // Mock mode
	}

	auth := smtp.PlainAuth("", e.config.SMTPUser, e.config.SMTPPassword, e.config.SMTPHost)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	addr := e.config.SMTPHost + ":" + e.config.SMTPPort
	err := smtp.SendMail(addr, auth, e.config.FromEmail, []string{to}, msg)

	if err != nil {
		logger.Error("Failed to send email", zap.Error(err))
		return err
	}

	logger.Info("Email sent successfully", zap.String("to", to))
	return nil
}

func (e *EmailService) SendPasswordResetEmail(to, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
Dear User,

You have requested a password reset. Please use the following token to reset your password:

Reset Token: %s

If you did not request this, please ignore this email.

Best regards,
Dashboard Team
`, resetToken)

	return e.SendEmail(to, subject, body)
}

func (e *EmailService) SendWelcomeEmail(to, name string) error {
	subject := "Welcome to Dashboard"
	body := fmt.Sprintf(`
Dear %s,

Welcome to our Dashboard platform!

Your account has been successfully created. You can now log in and start using our services.

Best regards,
Dashboard Team
`, name)

	return e.SendEmail(to, subject, body)
}

// Helper functions
func SendEmail(to, subject, body string) error {
	if Service == nil {
		return fmt.Errorf("email service not initialized")
	}
	return Service.SendEmail(to, subject, body)
}

func SendPasswordResetEmail(to, resetToken string) error {
	if Service == nil {
		return fmt.Errorf("email service not initialized")
	}
	return Service.SendPasswordResetEmail(to, resetToken)
}

func SendWelcomeEmail(to, name string) error {
	if Service == nil {
		return fmt.Errorf("email service not initialized")
	}
	return Service.SendWelcomeEmail(to, name)
}
