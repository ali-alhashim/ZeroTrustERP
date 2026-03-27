package core

import (
	"crypto/sha256"
	"crypto/hmac"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

// generateSecureToken by using the Email + sessionSecret key from .env and hashing it with sha256, this will be used for session management and should be stored in a secure cookie
func generateSessionToken(email string) string {
   

    sessionSecret := os.Getenv("sessionSecret")

	// Create a new HMAC instance using SHA256 and your secret key
	h := hmac.New(sha256.New, []byte(sessionSecret))

	// Write the email to the HMAC instance
     h.Write([]byte(email))

	return hex.EncodeToString(h.Sum(nil))
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

func isEmailExists(email string) bool {

	var count int

	query := "SELECT COUNT(*) FROM users WHERE email = ?"

	err := DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Println("DB error:", err)
		return false
	}

	return count > 0
}

func isEmailExistsAndActive(email string) bool {
	var exists bool

	email = strings.ToLower(strings.TrimSpace(email))

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND active = true)"

	err := DB.QueryRow(query, email).Scan(&exists)
	if err != nil {
		log.Println("DB error:", err)
		return false
	}

	return exists
}

func isEmailExistsAndActiveWithOTP(email string, otpHash string) bool {
	var valid bool

	email = strings.ToLower(strings.TrimSpace(email))

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND active = true AND otphash = $2 AND otpexpiry > NOW())"

	err := DB.QueryRow(query, email, otpHash).Scan(&valid)
	if err != nil {
		log.Println("DB error:", err)
		return false
	}

	return valid
}

func insertOtphashForEmail(email string, otpHash string, expiry time.Time) error {
	email = strings.ToLower(strings.TrimSpace(email))

	query := "UPDATE users SET otphash = $1, otpexpiry = $2 WHERE email = $3"

	_, err := DB.Exec(query, otpHash, expiry, email)
	return err
}

// ok this first time we run the app now users yet 0 so the first user to register will be admin by default, we can change this later if we want to add more roles
func isFirstUser() bool {
	var count int

	query := "SELECT COUNT(*) FROM users"

	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		log.Println("DB error:", err)
		return false
	}

	return count == 0
}

func isFirstUserAndFirstTimeLogin(email string) bool {
	email = strings.ToLower(strings.TrimSpace(email))

	var totalUsers int
	err := DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers)
	if err != nil {
		log.Println("DB error (count users):", err)
		return false
	}

	// Must be the ONLY user
	if totalUsers != 1 {
		return false
	}

	var count int
	query := "SELECT COUNT(*) FROM users WHERE email = $1 AND lastLogin IS NULL"

	err = DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Println("DB error (user check):", err)
		return false
	}

	return count == 1
}

// register the first user as admin by default
func registerFirstUserAsAdmin(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))

	query := "INSERT INTO users (email,  active, username) VALUES ($1, true,'system admin')"

	_, err := DB.Exec(query, email)
	if err != nil {
		log.Println("DB error (insert user):", err)
		return err
	}

	// we need to seed admin role and assign it to this user as well, we can do this in the same transaction
	// the role will be in roles table with name "admin" and description "System administrator with full permissions"
	// and to assign the role to the user we need to insert a record in the user_roles table with user_id and role_id

	queryRole := "INSERT INTO roles (name, description) VALUES ('admin', 'System administrator with full permissions') RETURNING id"

	var roleID int
	err2 := DB.QueryRow(queryRole).Scan(&roleID)
	if err2 != nil {
		log.Println("DB error (insert role):", err2)
		return err2
	}

	queryUserRole := "INSERT INTO user_roles (user_id, role_id) VALUES ((SELECT id FROM users WHERE email = $1), $2)"

	_, err = DB.Exec(queryUserRole, email, roleID)
	if err != nil {
		log.Println("DB error (assign role):", err)
		return err
	}

	return nil
}

// first time user login
func isFirstTimeLogin(email string) bool {

	var lastLogin sql.NullTime

	query := "SELECT lastLogin FROM users WHERE email = $1"

	err := DB.QueryRow(query, email).Scan(&lastLogin)
	if err != nil {
		log.Println("DB error:", err)
		return false
	}

	return !lastLogin.Valid // If last_login is NULL, it's the first time login
}



func incrementIncorrectOTPAttempts(email string) int {
	email = strings.ToLower(strings.TrimSpace(email))

	var attempts int

	// Get current attempts
	err := DB.QueryRow("SELECT incorrect_otp_attempts FROM users WHERE email = $1", email).Scan(&attempts)
	if err != nil {
		log.Println("DB error (get attempts):", err)
		return 0
	}

	attempts++

	fmt.Printf("Incorrect OTP attempt %d for %s\n", attempts, email)

	// Update attempts in DB
	_, err = DB.Exec("UPDATE users SET incorrectotpattempts = $1 WHERE email = $2", attempts, email)
	if err != nil {
		log.Println("DB error (update attempts):", err)
	}

	return attempts
}



func deleteUserByEmail(email string) {
	email = strings.ToLower(strings.TrimSpace(email))

	_, err := DB.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		log.Println("DB error (delete user):", err)
	}
}


func updateLastLogin(email string) {

	email = strings.ToLower(strings.TrimSpace(email))

	fmt.Printf("Updating last login for %s\n", email)

	_, err := DB.Exec("UPDATE users SET lastLogin = NOW() WHERE email = $1", email)
	if err != nil {
		log.Println("DB error (update last login):", err)
	}
}
