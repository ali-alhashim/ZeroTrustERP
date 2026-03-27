package core

import (
	"fmt"
	"net/http"
    
   
)

// RegisterRoutes sets up all HTTP request handlers
func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handleHealth)

	

	// Static file server for assets
	fs := http.FileServer(http.Dir("./static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    // Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Zero Trust ERP Server Running...")
        // root -> go to dashboard if logged in, otherwise show login page
        
	})



// LOGIN
mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodPost {
        email := r.FormValue("email")

        if !isValidEmail(email)  {
            RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
                "Error": "Invalid email",
                "Email": email,
            })
            return
        }

        // TODO: check email exists & is active + send OTP
        if !isEmailExistsAndActive(email) {
            RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
                "Error": "Email not found or inactive",
                "Email": email,
            })
            return
        }

         // ok email exists and is active, generate OTP and send email

        otp := GenerateOTP()

        mailErr := SendEmail(email, "Your OTP Code", fmt.Sprintf("Your OTP code is: %s", otp))
        if mailErr != nil {
            fmt.Printf("Failed to send OTP email: %v\n", mailErr)
            RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
                "Error": "Failed to send OTP email",
                "Email": email,
            })
            return
        }

        // ok OTP sent successfully, store OTP in  database with expiration time (e.g. 3 minutes)

        fmt.Printf("Generated OTP for %s: %s\n", email, otp)


          


        RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
            "Email": email,
        })
        return
    }

    RenderPageNoLayout(w, "core/templates/login.html", nil)
})


// OTP
mux.HandleFunc("/login-otp", func(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodPost {
        email := r.FormValue("email")
        otp := r.FormValue("otp")

        if !isValidEmail(email)  {
            RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
                "Error": "Invalid email",
                "Email": email,
            })
            return
        }

        fmt.Printf("Email: %s | OTP: %s\n", email, otp)

        // TODO: validate OTP

        RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
            "Error": "Invalid or expired OTP",
            "Email": email,
        })
        return
    }

    RenderPageNoLayout(w, "core/templates/login-otp.html", nil)
})




   for _, app := range registeredApps {
	app(mux)
   }

	return mux
}



// handleHealth returns API health status
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}






