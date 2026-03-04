package main

import (
	"context"
	"log"
	"os"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("SUPABASE_DATABASE_URL")
	if dbURL == "" {
		log.Fatal("SUPABASE_DATABASE_URL not set")
	}

	log.Println("🔄 Connecting to database...")
	log.Printf("URL: %s", maskPassword(dbURL))

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	// Example query to test connection
	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("❌ Query failed: %v", err)
	}

	log.Println("✅ Connected successfully!")
	log.Println("PostgreSQL version:", version)
}

func maskPassword(connStr string) string {
	// Simple masking for logging
	if len(connStr) > 50 {
		return connStr[:30] + "***" + connStr[len(connStr)-20:]
	}
	return "***"
}
