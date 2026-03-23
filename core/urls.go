package core

import (
	"fmt"
	"net/http"
	
)

type AppRoute func(*http.ServeMux)

var registeredApps []AppRoute

func Register(app AppRoute) {
	registeredApps = append(registeredApps, app)
}



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







