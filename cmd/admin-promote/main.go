package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect: %v\n", err)
	}
	defer pool.Close()

	log.Println("✅ Connected to database")

	// Get PlayFab ID from args or use first user
	var playfabID string
	if len(os.Args) > 1 {
		playfabID = os.Args[1]
		log.Printf("🔍 Looking for user with PlayFab ID: %s\n", playfabID)

		// Promote specific user
		result, err := pool.Exec(context.Background(),
			"UPDATE users SET is_admin = true WHERE playfab_id = $1",
			playfabID,
		)
		if err != nil {
			log.Fatalf("Failed to update user: %v\n", err)
		}

		if result.RowsAffected() == 0 {
			log.Printf("❌ No user found with PlayFab ID: %s\n", playfabID)
			log.Println("\n💡 Listing all users:")
			listUsers(pool)
			return
		}

		log.Printf("✅ User promoted to admin: %s\n", playfabID)

	} else {
		// No argument provided - list users and promote first one
		log.Print("🔍 Listing all users in database:\n\n")

		rows, err := pool.Query(context.Background(), `
			SELECT id, playfab_id, username, is_admin, created_at
			FROM users
			ORDER BY created_at DESC
			LIMIT 10
		`)
		if err != nil {
			log.Fatalf("Failed to query users: %v\n", err)
		}
		defer rows.Close()

		var count int
		var firstPlayFabID string

		for rows.Next() {
			var id, pfID, username string
			var isAdmin bool
			var createdAt string

			rows.Scan(&id, &pfID, &username, &isAdmin, &createdAt)

			if count == 0 {
				firstPlayFabID = pfID
			}

			adminStatus := ""
			if isAdmin {
				adminStatus = " [ADMIN]"
			}

			log.Printf("%d. PlayFab ID: %s | Username: %s%s\n", count+1, pfID, username, adminStatus)
			count++
		}

		if count == 0 {
			log.Println("❌ No users found in database")
			log.Println("\n💡 Create a user first by logging in from Unity")
			return
		}

		// Promote first user
		log.Printf("\n🎯 Promoting first user to admin: %s\n", firstPlayFabID)
		_, err = pool.Exec(context.Background(),
			"UPDATE users SET is_admin = true WHERE playfab_id = $1",
			firstPlayFabID,
		)
		if err != nil {
			log.Fatalf("Failed to promote user: %v\n", err)
		}

		log.Println("✅ First user promoted to admin!")
	}

	// Verify
	log.Println("\n📊 Current admin users:")
	rows, err := pool.Query(context.Background(), `
		SELECT playfab_id, username
		FROM users
		WHERE is_admin = true
	`)
	if err != nil {
		log.Fatalf("Failed to query admins: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pfID, username string
		rows.Scan(&pfID, &username)
		log.Printf("   👑 %s (%s)\n", username, pfID)
	}
}

func listUsers(pool *pgxpool.Pool) {
	rows, err := pool.Query(context.Background(), `
		SELECT playfab_id, username, is_admin
		FROM users
		ORDER BY created_at DESC
		LIMIT 10
	`)
	if err != nil {
		log.Printf("Failed to list users: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n📋 Available users:")
	var count int
	for rows.Next() {
		var pfID, username string
		var isAdmin bool
		rows.Scan(&pfID, &username, &isAdmin)

		adminStatus := ""
		if isAdmin {
			adminStatus = " [ADMIN]"
		}

		fmt.Printf("   %d. %s (%s)%s\n", count+1, username, pfID, adminStatus)
		count++
	}

	if count == 0 {
		fmt.Println("   No users found")
	}
}
