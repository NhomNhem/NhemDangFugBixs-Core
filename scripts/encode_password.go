package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run encode_password.go <your-password>")
		fmt.Println("Example: go run encode_password.go MyP@ss#123")
		os.Exit(1)
	}

	password := os.Args[1]
	encoded := url.QueryEscape(password)

	fmt.Println("Original password:", password)
	fmt.Println("Encoded password:", encoded)
	fmt.Println("\nUse the ENCODED password in your connection string!")
	fmt.Println("\nExample:")
	fmt.Printf("postgresql://postgres.xxxxx:%s@db.xxxxx.supabase.co:5432/postgres\n", encoded)
}
