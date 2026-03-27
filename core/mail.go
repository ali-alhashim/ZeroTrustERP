package core

import (
	"fmt"
	"net/smtp"
	"os"
	
)






// SendEmail sends a basic email
func SendEmail(to, subject, body string) error {

    LoadEnv(".env")
	
    smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")
	from     := os.Getenv("SMTP_FROM")

	// read SMTP config from environment variables
	fmt.Printf("SMTP Config - Host: %s, Port: %s, User: %s\n", smtpHost, smtpPort, smtpUser)




	if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("SMTP config is missing")
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		body

	addr := smtpHost + ":" + smtpPort

	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return err
	}

	return nil
}