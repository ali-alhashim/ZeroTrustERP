package core

import "net/http"

type App interface {
    Name() string
    Init() error                         // Runs table creations/migrations
    Routes(mux *http.ServeMux)           // Registers endpoints
    IsActive() bool                      // Allows deactivating (like the Setup app)
}