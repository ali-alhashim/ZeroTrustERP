package core

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"
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



