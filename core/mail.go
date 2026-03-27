package core

import (
	"fmt"
	"net/smtp"
	"os"
	"time"
	
)






// SendEmail sends a high-deliverability transactional email
func SendEmail(to, subject, body string) error {
    
    
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT")
    smtpUser := os.Getenv("SMTP_USERNAME")
    smtpPass := os.Getenv("SMTP_PASSWORD")
    from     := os.Getenv("SMTP_FROM") // e.g., "Zero Trust ERP <zeroerp8@gmail.com>"

    if smtpHost == "" || smtpPort == "" || smtpUser == "" || smtpPass == "" {
        return fmt.Errorf("SMTP config is missing")
    }

    auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

    // Critical headers for Inbox placement
    date := time.Now().Format(time.RFC1123Z)
    messageId := fmt.Sprintf("<%d.%s@%s>", time.Now().UnixNano(), "system", smtpHost)

    // Constructing the message with headers that satisfy Spam Filters
    msg := fmt.Sprintf("From: %s\r\n"+
        "To: %s\r\n"+
        "Subject: %s\r\n"+
        "Date: %s\r\n"+
        "Message-ID: %s\r\n"+
        "MIME-Version: 1.0\r\n"+
        "Content-Type: text/html; charset=\"UTF-8\"\r\n"+
        "Content-Transfer-Encoding: 7bit\r\n"+
        "Precedence: bulk\r\n"+                  // Identifies as automated/transactional
        "Auto-Submitted: auto-generated\r\n"+    // Helps bypass "Reply-To" filters
        "X-Priority: 3\r\n"+                     // Normal Priority
        "X-Mailer: ZeroTrustERP-Go\r\n"+
        "\r\n"+
        "<!DOCTYPE html><html><body>%s</body></html>", 
        from, to, subject, date, messageId, body)

    addr := smtpHost + ":" + smtpPort

    // Ensure the 'from' email in SendMail matches the one in the headers
    // If 'from' in .env is "Name <email@gmail.com>", we must extract just the email
    fromEmail := smtpUser 
    
    err := smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(msg))
    if err != nil {
        return err
    }

    return nil
}