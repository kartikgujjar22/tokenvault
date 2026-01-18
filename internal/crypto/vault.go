package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

func Encrypt(plaintext string) (string, error) {
	key, err := GetKey()
	if err != nil {
		return "", err
	}

	// Create a new AES Cipher Block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Wrap it in GCM (Galois/Counter Mode) - handles encryption & authentication
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a Nonce (Number used ONCE).
	// Standard size for GCM is 12 bytes.
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	// Seal(dst, nonce, plaintext, additionalData)
	// We prepend the nonce to the ciphertext so we can use it for decryption later
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return as Hex string so we can save it in SQLite text column
	return hex.EncodeToString(ciphertext), nil
}

// Decrypt takes the hex-encoded string and returns the original text
func Decrypt(encryptedHex string) (string, error) {
	key, err := GetKey()
	if err != nil {
		return "", err
	}

	data, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Split the nonce and the actual encrypted msg
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.New("decryption failed: invalid key or corrupted data")
	}

	return string(plaintext), nil
}
