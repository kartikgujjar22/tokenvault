package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestAnalyzeToken_JWT(t *testing.T) {
	// A fake JWT with specific claims:
	// sub: "user_123"
	// exp: 4807490212 (Year 2122 - definitely in the future)
	// Body: {"sub":"user_123", "name":"Kartik", "exp":4807490212}
	fakeJWT := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyXzEyMyIsIm5hbWUiOiJLYXJ0aWsiLCJleHAiOjQ4MDc0OTAyMTJ9.SIGNATURE_IGNORE"

	fmt.Println("Testing JWT Inspector...")
	meta := AnalyzeToken(fakeJWT)

	//Check if identified as JWT
	if !meta.IsJWT {
		t.Error("Failed: Should be identified as JWT")
	}

	//Check Subject
	if meta.Subject != "user_123" {
		t.Errorf("Failed: Expected subject 'user_123', got '%s'", meta.Subject)
	}

	// Check Expiry (Should be the specific timestamp in the token)
	if meta.ExpiresAt != 4807490212 {
		t.Errorf("Failed: Expected exp 4807490212, got %d", meta.ExpiresAt)
	}

	fmt.Println("Success: JWT parsed correctly!")
	fmt.Printf(" -> User: %s\n   -> Expires: %d\n", meta.Subject, meta.ExpiresAt)
}

func TestAnalyzeToken_APIKey(t *testing.T) {
	apiKey := "sk_test_51Mz..."

	meta := AnalyzeToken(apiKey)

	if meta.Type != "Stripe-Key" {
		t.Errorf("Expected Type 'Stripe-Key', got '%s'", meta.Type)
	}

	// Default expiry (24h) check: should be greater than now
	if meta.ExpiresAt < time.Now().Unix() {
		t.Error("Expected default expiry to be in the future")
	}

	fmt.Println("Success: API Key identified correctly!")
}
