package crypto

import (
	"fmt"
	"testing"
)

func TestEncryption(t *testing.T) {
	secret := "MySuperSecretToken_123"
	fmt.Println("Testing with secret:", secret)

	// Try to Encrypt
	encrypted, err := Encrypt(secret)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}
	fmt.Println("Encrypted Hex:", encrypted)

	// Try to Decrypt
	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	//Verify they match
	if decrypted != secret {
		t.Errorf("Expected %s but got %s", secret, decrypted)
	} else {
		fmt.Println("Success: Decrypted value matches original!")
	}
}
