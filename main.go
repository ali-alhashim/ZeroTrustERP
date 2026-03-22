package main

import (
	"log"
	"net/http"
	"strconv"

	"zerotrusterp/core"
)

func StartServer(port string) {
	cfg := core.DefaultServerConfig()
	if port != "" {
		p, err := strconv.Atoi(port)
		if err == nil && p > 0 {
			cfg.Port = p
		}
	}
	StartServerWithConfig(cfg)
}

func StartServerWithConfig(cfg *core.ServerConfig) {
	if cfg == nil {
		cfg = core.DefaultServerConfig()
	}

	
	
	

	server := &http.Server{
		Addr:           cfg.String(),
		Handler:        core.RegisterRoutes(),
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.ShutdownTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	cfg.LogConfig()
	log.Printf("Starting server on %s in %s mode", cfg.String(), cfg.Environment)


	//ok server working so test the database connection
	db, err := core.InitDB()

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	} else {
		log.Printf("Database connection successful")
		db.Close()
	}




	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Printf("Server closed")
		} else {
			log.Fatalf("Server failed: %v", err)
		}
	}
	// Implement graceful shutdown logic here if needed

	
}


func main() {
	StartServer("8000")
}
