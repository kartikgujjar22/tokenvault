package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func GetKey() ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("Could not find user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".tokenvault")
	keyPath := filepath.Join(configDir, "master.key")

	//creating the directory if it does not exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0700); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	//check if key exist or not
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return generateAndSaveKey(keyPath)
	}

	// reading the existing key in the directpry
	keyHex, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key: %v", err)
	}

	//decoding the hex back to byes
	keyBytes, err := hex.DecodeString(string(keyHex))
	if err != nil {
		return nil, fmt.Errorf("corrupt key file: %v", err)
	}

	return keyBytes, nil

}

func generateAndSaveKey(path string) ([]byte, error) {
	fmt.Println("Generating new master key")

	//generate 32 bytes (256 bits) of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, filepath.ErrBadPattern

	}

	// convert to hex string for safe storage in a text file
	keyString := hex.EncodeToString(bytes)

	//save to file
	if err := os.WriteFile(path, []byte(keyString), 0600); err != nil {
		return nil, err
	}

	return bytes, nil
}
