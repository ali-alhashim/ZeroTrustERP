package core

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	_ "github.com/lib/pq"   //go get github.com/lib/pq the download location will be $HOME/go/pkg/mod
)

// LoadEnv reads .env file and injects into system environment
func LoadEnv(path string) error {
    file, err := os.Open(path)
    if err != nil {
		log.Printf("No .env file found at %s, skipping environment variable loading", path)
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        // Skip comments and empty lines
        if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
            continue
        }
        // Split by the first "=" only
        parts := strings.SplitN(line, "=", 2)
        if len(parts) == 2 {
            key := strings.TrimSpace(parts[0])
            val := strings.TrimSpace(parts[1])
            os.Setenv(key, val)
        }
    }
    return scanner.Err()
}

func InitDB() (*sql.DB, error) {
    log.Printf("Initializing database connection...")
    LoadEnv(".env")

    host := os.Getenv("DB_HOST")
    port := os.Getenv("DB_PORT")
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

    // 1. Open the connection (Lazy)
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

    // 2. Ping is where the "Database does not exist" error actually shows up
    err = db.Ping()
    if err != nil {
        // Check for the error code 3D000 in the Ping error
        if strings.Contains(err.Error(), "3D000") || strings.Contains(err.Error(), "does not exist") {
            log.Printf("Database %s not found. Attempting to create it...", dbname)
            
            // Connect to default 'postgres' db to run CREATE command
            adminDsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", host, port, user, password)
            adminDb, adminErr := sql.Open("postgres", adminDsn)
            if adminErr != nil {
                return nil, fmt.Errorf("could not connect to admin db: %w", adminErr)
            }
            defer adminDb.Close()

            _, createErr := adminDb.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
            if createErr != nil {
                return nil, fmt.Errorf("failed to create database: %w", createErr)
            }
            log.Printf("✅ Database %s created successfully!", dbname)
            
            // Close the old "failed" db connection and open a fresh one
            db.Close()
            return sql.Open("postgres", dsn)
        }
        return nil, fmt.Errorf("connection failed: %w", err)
    }

    log.Printf("✅ Database connected successfully to %s:%s/%s", host, port, dbname)
    return db, nil
}



