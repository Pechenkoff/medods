package utils

import (
	"fmt"
	"medods/internal/entities"

	"gopkg.in/gomail.v2"
)

type emailUtils struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
}

type EmailUtils interface {
	SendEmail(emailReq entities.EmailRequest) error
}

// NewEmailUtils - create a new copy od EmailUtils
func NewEmailUtils(smtpHost, smtpUsername, smtpPassword string, smtpPort int) EmailUtils {
	return &emailUtils{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
	}
}

// SendEmail - send email to user with body
func (u *emailUtils) SendEmail(emailReq entities.EmailRequest) error {
	message := gomail.NewMessage()
	message.SetHeader("From", emailReq.From)
	message.SetHeader("To", emailReq.To)
	message.SetHeader("Subject", emailReq.Subject)
	message.SetBody("text/plain", emailReq.Body)

	dialer := gomail.NewDialer(u.SMTPHost, u.SMTPPort, u.SMTPUsername, u.SMTPPassword)

	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("could not send email: %v", err)
	}

	return nil
}
