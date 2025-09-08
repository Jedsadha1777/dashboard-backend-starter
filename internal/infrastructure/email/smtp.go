package email

import (
	"fmt"
	"log"
	"net/smtp"
)

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

type SMTPEmailService struct {
	config Config
}

func NewSMTPEmailService(config Config) *SMTPEmailService {
	service := &SMTPEmailService{config: config}

	if config.SMTPHost != "" {
		log.Printf("Email service initialized with SMTP host: %s", config.SMTPHost)
	} else {
		log.Println("Email service running in mock mode (no SMTP configured)")
	}

	return service
}

func (e *SMTPEmailService) SendEmail(to, subject, body string) error {
	if e.config.SMTPHost == "" {
		log.Printf("Mock email sent to: %s, subject: %s", to, subject)
		return nil
	}

	auth := smtp.PlainAuth("", e.config.SMTPUser, e.config.SMTPPassword, e.config.SMTPHost)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	addr := e.config.SMTPHost + ":" + e.config.SMTPPort
	err := smtp.SendMail(addr, auth, e.config.FromEmail, []string{to}, msg)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	log.Printf("Email sent successfully to: %s", to)
	return nil
}

func (e *SMTPEmailService) SendPasswordResetEmail(to, resetToken string) error {
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

func (e *SMTPEmailService) SendWelcomeEmail(to, name string) error {
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
