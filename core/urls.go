package core

import (
	"fmt"
	"net/http"
	"strings"
	"time"
	

	
)

var MainHub *Hub

func GetMainHub() *Hub {

	return MainHub
}

// RegisterRoutes sets up all HTTP request handlers
func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handleHealth)

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))


     // WebSocket endpoint
     MainHub = NewHub()
	 go MainHub.Run()

	

     mux.HandleFunc("/ws", WebSocketHandler(MainHub))




	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Zero Trust ERP Server Running... you reached the root endpoint.\n")
		//redirect to dashboard.
		
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			
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

        email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))

		if r.Method == http.MethodPost {
			
			otp := r.FormValue("otp")

            fmt.Printf("Received OTP login attempt for %s\n", email)


			if !isEmailExistsAndActive(email) {
				RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
					"Error": "Invalid login",
					"Email": email,
				})
				return
			}

            fmt.Printf("Processing OTP login for %s\n", email)

			hashedOTP := HashOTP(otp)

			// Check OTP (must include expiry check in DB)
			isValidOTP := isEmailExistsAndActiveWithOTP(email, hashedOTP)

            fmt.Printf("OTP validation result for %s: %v\n", email, isValidOTP)

			// First-time login handling
			if isFirstUserAndFirstTimeLogin(email) {
				if !isValidOTP {
                     fmt.Printf("First-time login failed OTP validation for %s\n", email)

					attempts := incrementIncorrectOTPAttempts(email)
					if attempts >= 5 {
						deleteUserByEmail(email)
                        fmt.Printf("Deleted user %s after %d failed OTP attempts\n", email, attempts)
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

                fmt.Printf("First-time login successful for %s\n", email)
			}

			// Normal OTP validation
			if !isValidOTP {
                    fmt.Printf("Normal OTP validation failed for %s\n", email)
				RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
					"Error": "Invalid or expired OTP",
					"Email": email,
				})
				return
			}

			// ✅ SUCCESS LOGIN
            fmt.Printf("OTP login successful for %s\n", email)
			// Set secure session cookie

            fmt.Printf("Setting session cookie for %s\n", email)

			http.SetCookie(w, &http.Cookie{
				Name:     "session",
				Value:    generateSessionToken(email),
				HttpOnly: true,
				Secure:   true,
				Path:     "/",
				Expires:  time.Now().Add(24 * time.Hour),
			})


			http.SetCookie(w, &http.Cookie{
				Name:    "email",
				Value:    email,
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

        fmt.Printf("Rendering OTP login page for %s\n", email)
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







