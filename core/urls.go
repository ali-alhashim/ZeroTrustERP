package core

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// RegisterRoutes sets up all HTTP request handlers
func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handleHealth)

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Zero Trust ERP Server Running...")
	})

	// =========================
	// LOGIN
	// =========================
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))

			// Validate email
			if !isValidEmail(email) {
				RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
					"Error": "Invalid email",
					"Email": email,
				})
				return
			}

			// ⚠️ First user bootstrap (be careful in production)
			if isFirstUser() {
				err := registerFirstUserAsAdmin(email)
				if err != nil {
					fmt.Printf("Failed to register first user: %v\n", err)
					RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
						"Error": "System initialization failed",
						"Email": email,
					})
					return
				}
			}

			// Generic message (avoid user enumeration)
			if !isEmailExistsAndActive(email) {
				RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
					"Error": "Invalid login",
					"Email": email,
				})
				return
			}

			// Generate OTP
			otp := GenerateOTP()

			// Send email
			err := SendEmail(
				email,
				"Your OTP Code",
				fmt.Sprintf("Your verification code is: %s (expires in 3 minutes)", otp),
			)
			if err != nil {
				fmt.Printf("Failed to send OTP email: %v\n", err)
				RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
					"Error": "Failed to send OTP",
					"Email": email,
				})
				return
			}

			// Store OTP (hashed + expiry)
			err = insertOtphashForEmail(email, HashOTP(otp), time.Now().Add(3*time.Minute))
			if err != nil {
				fmt.Printf("Failed to store OTP: %v\n", err)
				RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
					"Error": "System error",
					"Email": email,
				})
				return
			}

			// Do NOT log OTP in production
			fmt.Printf("OTP generated for %s\n", email)

			RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
				"Email": email,
			})
			return
		}

		RenderPageNoLayout(w, "core/templates/login.html", nil)
	})

	// =========================
	// OTP LOGIN
	// =========================
	mux.HandleFunc("/login-otp", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
			otp := r.FormValue("otp")

			if !isEmailExistsAndActive(email) {
				RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
					"Error": "Invalid login",
					"Email": email,
				})
				return
			}

			hashedOTP := HashOTP(otp)

			// Check OTP (must include expiry check in DB)
			isValidOTP := isEmailExistsAndActiveWithOTP(email, hashedOTP)

			// First-time login handling
			if isFirstUserAndFirstTimeLogin(email) {
				if !isValidOTP {

					attempts := incrementIncorrectOTPAttempts(email)
					if attempts >= 5 {
						deleteUserByEmail(email)
						RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
							"Error": "Too many failed attempts. Please register again.",
							"Email": email,
						})
						return
					}

					RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
						"Error": "Check your email/spam. Ensure correct email.",
						"Email": email,
					})
					return
				}
			}

			// Normal OTP validation
			if !isValidOTP {
				RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
					"Error": "Invalid or expired OTP",
					"Email": email,
				})
				return
			}

			// ✅ SUCCESS LOGIN

			// Set secure session cookie
			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    generateSessionToken(email),
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
			})

			// Update last login
			updateLastLogin(email)

			// Redirect to dashboard
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}

		RenderPageNoLayout(w, "core/templates/login-otp.html", nil)
	})

	// External apps
	for _, app := range registeredApps {
		app(mux)
	}

	return mux
}

// =========================
// HEALTH
// =========================
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
