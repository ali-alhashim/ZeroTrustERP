package core

import (
	
	"encoding/hex"
	"net/http"
	"time"
	"crypto/sha256"
	"regexp"
	"fmt"
	"math/rand"
)

// generateSecureToken creates a random secure token
func generateSecureToken(email string) string {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return ""
	}
	return hex.EncodeToString(token)
}

// CreateSession generates a secure, HttpOnly cookie
func CreateSession(w http.ResponseWriter, email string) {
	cookie := &http.Cookie{
		Name:     "erp_session",
		Value:    generateSecureToken(email), // Encrypted string
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true, // Prevents JS from stealing the session (XSS Protection)
		Secure:   true, // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)
}


func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func HashOTP(otp string) string {
	hash := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(hash[:])
}




func isValidEmail(email string) bool {
    // Simple regex for email validation
    re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
    return re.MatchString(email)
}

func isValidMobile(mobile string) bool {
    // Simple regex for mobile number validation (Saudi Arabia)
    re := regexp.MustCompile(`^966[5-9]\d{8}$`)
    return re.MatchString(mobile)
}



