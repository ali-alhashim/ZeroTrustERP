package core

import (
	"fmt"
	"log"
	"time"
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
		Port:            8080,
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
