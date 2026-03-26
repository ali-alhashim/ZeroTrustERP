package core

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"strconv"
)

// ServerConfig holds all server settings
type ServerConfig struct {
	Port            int
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	MaxHeaderBytes  int
	Environment     string
}

// DefaultServerConfig returns default server settings
func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:            8000,
		Host:            "0.0.0.0",
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
		ShutdownTimeout: 30 * time.Second,
		MaxHeaderBytes:  1 << 20, // 1MB
		Environment:     "development",
	}
}

// String returns the server address
func (c *ServerConfig) String() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LogConfig logs the current server configuration
func (c *ServerConfig) LogConfig() {
	log.Printf("Server Config - Host: %s, Port: %d, Environment: %s", c.Host, c.Port, c.Environment)
}



func StartServer(port string) {
	cfg := DefaultServerConfig()
	if port != "" {
		p, err := strconv.Atoi(port)
		if err == nil && p > 0 {
			cfg.Port = p
		}
	}
	StartServerWithConfig(cfg)
}

func StartServerWithConfig(cfg *ServerConfig) {
	if cfg == nil {
		cfg = DefaultServerConfig()
	}

	
	
	

	server := &http.Server{
		Addr:           cfg.String(),
		Handler:         RegisterRoutes(),
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		IdleTimeout:    cfg.ShutdownTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	cfg.LogConfig()
	log.Printf("Starting server on %s in %s mode", cfg.String(), cfg.Environment)


	//ok server working so test the database connection
	db, err := InitDB()

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

