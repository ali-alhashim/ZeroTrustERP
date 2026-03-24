package core

import (
	"fmt"
	"net/http"
    "regexp"
	
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




	 // Login page
    mux.HandleFunc("/login-otp", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            username := r.FormValue("username")
           
            // check if username is email or mobile number
            if !isValidEmail(username) && !isValidMobile(username) {
                fmt.Printf("Invalid login attempt with username: %s\n", username)
                data := map[string]interface{}{"Error": "Invalid email or mobile number"}
                RenderPage(w, "core/templates/login.html", data)
                return
            }

            // is valid email or mobile number, check if user exists in database
            // if user does not exist, return error message
            // if user exists, generate OTP and send to user via email or SMS
            // for now, we will just print the OTP to the console
            fmt.Printf("login-otp: %s\n", username)
            
            
            // Simple error message
            data := map[string]interface{}{"Error": "Invalid credentials"}
            RenderPage(w, "core/templates/login.html", data)
            return
        }
        RenderPage(w, "core/templates/login.html", nil)
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



