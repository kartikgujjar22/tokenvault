package database

import (
	"fmt"
	"testing"
)

func TestDatabaseFlow(t *testing.T) {
	// Setup - Init DB
	if err := InitDB(); err != nil {
		t.Fatalf("Failed to init DB: %v", err)
	}
	defer DB.Close()

	// Store a Token
	project := "test-project"
	secret := "super-secret-v2-token"

	fmt.Println("Storing Token...")
	if err := StoreToken(project, "default", secret); err != nil {
		t.Fatalf("Store failed: %v", err)
	}

	// Fetch it back
	fmt.Println("Fetching Token...")
	val, meta, err := FetchToken(project, "default")
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	// Verify
	if val != secret {
		t.Errorf("Mismatch! Got %s, want %s", val, secret)
	}

	fmt.Printf("Success! Token Type detected as: %s\n", meta.Type)

}
