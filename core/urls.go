package core

import (
	"fmt"
	"net/http"
	"zerotrusterp/apps/users"
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

    // Register app-specific routes [don't forget to import the app packages at the top example: "zerotrusterp/apps/users"]
	users.RegisterRoutes(mux)

	return mux
}



// handleHealth returns API health status
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}







