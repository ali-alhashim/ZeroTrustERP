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
	})



// LOGIN
mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodPost {
        username := r.FormValue("username")

        if !isValidEmail(username) && !isValidMobile(username) {
            RenderPageNoLayout(w, "core/templates/login.html", map[string]interface{}{
                "Error": "Invalid email or mobile number",
                "Username": username,
            })
            return
        }

        // TODO: check user exists + send OTP

        RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
            "Username": username,
        })
        return
    }

    RenderPageNoLayout(w, "core/templates/login.html", nil)
})


// OTP
mux.HandleFunc("/login-otp", func(w http.ResponseWriter, r *http.Request) {

    if r.Method == http.MethodPost {
        username := r.FormValue("username")
        otp := r.FormValue("otp")

        if !isValidEmail(username) && !isValidMobile(username) {
            RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
                "Error": "Invalid email or mobile number",
                "Username": username,
            })
            return
        }

        fmt.Printf("Username: %s | OTP: %s\n", username, otp)

        // TODO: validate OTP

        RenderPageNoLayout(w, "core/templates/login-otp.html", map[string]interface{}{
            "Error": "Invalid or expired OTP",
            "Username": username,
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






