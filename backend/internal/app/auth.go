package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"goauth/internal/config"
	"goauth/internal/db"
	"goauth/internal/httpServer"
	"goauth/internal/model/client"
	"goauth/internal/model/user"
	"goauth/internal/utils"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"
)

// Run initializes and starts the application
func Run() {

	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		if err != nil {
			return
		}
	}()

	// Setup error handling
	if err := setup(); err != nil {
		log.Fatalf("Failed to setup application: %v", err)
	}

	// Start HTTP server
	if err := httpServer.CreateHTTPServer(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// setup initializes all application components
func setup() error {
	// Load environment variables
	if err := loadConfig(); err != nil {
		return fmt.Errorf("config loading error: %w", err)
	}

	// Initialize utilities
	utils.Init()

	// Initialize database
	if err := initDatabase(); err != nil {
		return fmt.Errorf("database initialization error: %w", err)
	}

	// Create test data
	if err := createTestData(); err != nil {
		return fmt.Errorf("test data creation error: %w", err)
	}

	return nil
}

// loadConfig loads configuration from .env file
func loadConfig() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	config.Load()

	// Validate required configuration
	if config.Config.DbHost == "" || config.Config.DbPort == "" ||
		config.Config.DbUser == "" || config.Config.DbName == "" {
		return errors.New("missing required database configuration")
	}

	return nil
}

// initDatabase initializes database connection with retry logic
func initDatabase() error {
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		log.Printf("Attempting database connection (attempt %d/%d)...", i+1, maxRetries)

		// Create context with timeout for DB initialization
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Channel to capture async result
		done := make(chan error, 1)

		go func() {
			done <- db.DB.Connect()
		}()

		// Wait for DB init or timeout
		select {
		case err = <-done:
			if err == nil {
				log.Println("Database connection established successfully")
				return nil
			}
			log.Printf("Database connection failed: %v", err)
		case <-ctx.Done():
			err = ctx.Err()
			log.Printf("Database connection timed out: %v", err)
		}

		// Wait before retry
		if i < maxRetries-1 {
			backoffTime := time.Duration(i+1) * time.Second
			log.Printf("Waiting %v before next attempt...", backoffTime)
			time.Sleep(backoffTime)
		}
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

// createTestData creates test users and clients
func createTestData() error {
	// Only create test data in development environment
	if os.Getenv("APP_ENV") != "production" {
		// Create test user if it doesn't exist
		newUserId, err := user.Create("testuser", "testpassword")
		if err != nil {
			// Ignore error if user already exists
			if !errors.Is(err, user.ErrUserAlreadyExists) &&
				!strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return fmt.Errorf("failed to create test user: %w", err)
			}
			log.Println("Test user already exists, skipping creation")
		} else {
			log.Println("Created test user successfully")
			_, err = client.Create(newUserId)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE constraint failed") {
					log.Println("Test client already exists, skipping creation")
				} else {
					return fmt.Errorf("failed to create test client: %w", err)
				}
			} else {
				log.Println("Created test client successfully")
			}
		}
	}

	return nil
}
