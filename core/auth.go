package core

import (
	"crypto/sha256"
	
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

    "context"
	"crypto/rand"
    "encoding/base64"


	"regexp"
	"strings"
	"time"
	"net/http"
	"net"

	"zerotrusterp/apps/users/models"
)

// generateSecureToken by using the Email + sessionSecret key from .env and hashing it with sha256, this will be used for session management and should be stored in a secure cookie
func generateSessionToken(email string) string {
    // 1. Create a byte slice of 64 bytes (512 bits of entropy)
    b := make([]byte, 64)

    // 2. Fill the slice with random bytes from the OS's secure source
    _, err := rand.Read(b)
    if err != nil {
        // In a real app, handle this error properly (e.g., return an error)
        panic("Could not generate random bytes: " + err.Error())
    }

    // 3. Encode the bytes to a URL-safe string
    // This turns the raw bytes into a readable string like "u6B_...X8Q"
    sessionToken := base64.URLEncoding.EncodeToString(b)

    // 4. Register the SessionToken in the database
    RegisterSessionTokenInDB(email, sessionToken)

	//we generate a token so means the use is login so log the action
    InsertLog(GetUserByEmail(email), "authentication", "login")

    return sessionToken
}


func InsertLog(user *models.User, resource string, action string) {

	// TODO: implement log insertion to database, create a new log record with user id, resource, action and timestamp
	query := "INSERT INTO logs (user_id, username, email, resource, action) VALUES ($1, $2, $3, $4, $5)"
	_, err := DB.Exec(query, user.ID, user.Username, user.Email, resource, action)
	if err != nil {
		panic(err)
	}
	
}




func GenerateOTP() string {
    // Generate a random number between 0 and 999,999
    n, err := rand.Int(rand.Reader, big.NewInt(1000000))
    if err != nil {
        return "000000" // Fallback or handle error
    }
    return fmt.Sprintf("%06d", n.Int64())
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


func SetUserIsOnline(email string, value bool) { 
 
	email = strings.ToLower(strings.TrimSpace(email))

	fmt.Printf("Setting user %s as online\n", email)

	_, err := DB.Exec("UPDATE users SET online = $1 WHERE email = $2", value, email)
	if err != nil {
		log.Println("DB error (set user online):", err)
	}
	
}

func isEmailExistsAndActiveWithOTP(email string, otpHash string) bool {
	var valid bool

	email = strings.ToLower(strings.TrimSpace(email))

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND active = true AND otp_hash = $2 AND otp_expiry > NOW())"

	err := DB.QueryRow(query, email, otpHash).Scan(&valid)
	if err != nil {
		log.Println("DB error:", err)
		return false
	}

	return valid
}

func insertOtphashForEmail(email string, otpHash string, expiry time.Time) error {
	email = strings.ToLower(strings.TrimSpace(email))

	query := "UPDATE users SET otp_hash = $1, otp_expiry = $2 WHERE email = $3"

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


func GetUserIDByEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
    
	fmt.Printf("Getting user ID for email: %s\n", email)
	var id string

	query := "SELECT id FROM users WHERE email = $1"

	err := DB.QueryRow(query, email).Scan(&id)
	if err != nil {
		log.Println("DB error:", err)
		return "", err
	}

	return id, nil
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
	query := "SELECT COUNT(*) FROM users WHERE email = $1 AND last_Login IS NULL"

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

	queryUserRole := "INSERT INTO users_roles (user_id, role_id) VALUES ((SELECT id FROM users WHERE email = $1), $2)"

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

	query := "SELECT last_Login FROM users WHERE email = $1"

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

	_, err := DB.Exec("UPDATE users SET last_Login = NOW() WHERE email = $1", email)
	if err != nil {
		log.Println("DB error (update last login):", err)
	}
}



func RegisterSessionTokenInDB(email string, sessionToken string) {

	email = strings.ToLower(strings.TrimSpace(email))

	fmt.Printf("Registering session token for %s\n", email)

	_, err := DB.Exec("UPDATE users SET session_token = $1 WHERE email = $2", sessionToken, email)
	if err != nil {
		log.Println("DB error (register session token):", err)
	}
}


// check he has a valid session 
func IsValidSessionToken(email string, sessionToken string) bool {

	email = strings.ToLower(strings.TrimSpace(email))

	var valid bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND session_token = $2)"

	err := DB.QueryRow(query, email, sessionToken).Scan(&valid)
	if err != nil {
		log.Println("DB error:", err)
		return false
	}

	return valid
}

// some url need permission this can be check in the database permission if the resource is protected or not and 
// if the user has permission to access it or not

// Define a custom type for context keys to avoid collisions
type contextKey string
const UserKey contextKey = "currentUser"

func AuthMiddleware(next http.Handler, resource...string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. Automatically read the cookies from the request
        emailCookie, errE := r.Cookie("email")
        sessionCookie, errS := r.Cookie("session")

        // 2. If cookies are missing, they aren't logged in
        if errE != nil || errS != nil {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        // 3. Use  function to check the database
        if !IsValidSessionToken(emailCookie.Value, sessionCookie.Value) {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            return
        }

        // 4. Success! Let them through to the next page
		// but we need to check the permissions
		// what is the resource they are trying to access and does the user have permission to access it or not so pass the resource name in the header and check it in the database if the user has permission to access it or not
		if(len(resource) > 0){
			fmt.Printf(" Checking permissions for %s on resource %s\n", emailCookie.Value, resource[0])
            
			parts := strings.Split(resource[0], ":")

			resourceName := parts[0] 
            actionType   := parts[1] 
			fmt.Print("\n ok the resource is: ", resourceName, " and the action type is: ", actionType,"\n")
			
			if !isAuthorized(GetUserByEmail(emailCookie.Value),resourceName, actionType) {
                fmt.Print("You Are not Authorized for the Selected Resource ", resource[0])
				
				//HTTP 403 Forbidden 
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("403 - You are not authorized for this resource"))
				return // Stop execution here!
			}

			fmt.Print("\n ***Ok the User is Authorized*** \n")
		}

        
    
        ctx := context.WithValue(r.Context(), UserKey, GetUserByEmail(emailCookie.Value))


        next.ServeHTTP(w, r.WithContext(ctx))
    })
}


func GetUserByEmail(email string) *models.User {
    
	fmt.Printf("Getting user by email: %s\n", email)

	email = strings.ToLower(strings.TrimSpace(email))

	var user models.User

	query := "SELECT id, username, email, active FROM users WHERE email = $1"

	err := DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Active)
	if err != nil {
		log.Println("DB error:", err)
		return nil
	}

	return &user
}


func GetRealIP(r *http.Request) string {
	var ip string

	// 1. Priority: Cloudflare-specific header (Trusted if using CF)
	if cfIP := r.Header.Get("CF-Connecting-IP"); cfIP != "" {
		ip = cfIP
	}

	// 2. Secondary: Standard Proxy header
	if ip == "" {
		xForwardedFor := r.Header.Get("X-Forwarded-For")
		if xForwardedFor != "" {
			// X-Forwarded-For can be a comma-separated list. 
			// The first one is the original client.
			parts := strings.Split(xForwardedFor, ",")
			ip = strings.TrimSpace(parts[0])
		}
	}

	// 3. Tertiary: Other common proxy headers
	if ip == "" {
		if xRealIP := r.Header.Get("X-Real-IP"); xRealIP != "" {
			ip = xRealIP
		}
	}

	// 4. Final Fallback: The direct connection IP
	if ip == "" {
		// net.SplitHostPort correctly handles IPv4 and IPv6 [brackets]
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If it fails (no port), just use the raw string
			ip = r.RemoteAddr
		} else {
			ip = host
		}
	}

	// Clean up any remaining whitespace or weird formatting
	return strings.TrimSpace(ip)
}


func GetAllResources() []string{

	fmt.Println("Fetching all resources for permissions...")

	var resources []string
    
	//resoures is the names of tables in the database but not relations tables like users_roles or roles_permissions 
	// just the main tables that we want to protect with permissions 
	// like users, products, orders etc.. so we can check the permissions based on the resource name and action in the database
	query := `
        SELECT table_name 
        FROM information_schema.columns 
        WHERE table_schema = 'public' 
          AND column_name = 'id'
          AND table_name NOT IN ('%_%')
        ORDER BY table_name;`

	
	rows, err := DB.Query(query)
    if err != nil {
        fmt.Println("Error fetching resources:", err)
        return resources
    }
    defer rows.Close()

    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            fmt.Println("Error scanning resource name:", err)
            continue
        }
        resources = append(resources, name)
    }

	fmt.Printf("Resources found: %v\n", resources)


	return resources
}


func isAuthorized(theUser *models.User, resource string, action string) bool {
    // This query checks two things:
    // 1. Is the user an 'admin'?
    // 2. Does the user have a role linked to the specific resource and action?
	fmt.Print(" \n check if the user with ID: ", theUser.ID, " has permisstion as ", action , " on ", resource ,"\n")

    query := `
        SELECT EXISTS (
            SELECT 1 
            FROM users_roles ur
            JOIN roles r ON ur.role_id = r.id
            LEFT JOIN roles_permissions rp ON r.id = rp.role_id
            LEFT JOIN permissions p ON rp.permission_id = p.id
            WHERE ur.user_id = $1 
            AND (
                r.name = 'admin' 
                OR (p.resource = $2 AND p.action = $3)
            )
        )`

    var authorized bool
    err := DB.QueryRow(query, theUser.ID, resource, action).Scan(&authorized)
    if err != nil {
        fmt.Println("Error checking authorization:", err)
        return false
    }

    return authorized
}