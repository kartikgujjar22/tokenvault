package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kartikgujjar22/tokenvault/internal/crypto"
	"github.com/kartikgujjar22/tokenvault/internal/utils"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB opens the connection and runs migrations
func InitDB() error {
	home, _ := os.UserHomeDir()
	dbPath := filepath.Join(home, ".tokenvault", "token.db")

	// Ensure folder exists
	_ = os.MkdirAll(filepath.Dir(dbPath), 0700)

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Ping to check connection
	if err = DB.Ping(); err != nil {
		return err
	}

	// Run our Auto-Migration
	return RunMigrations(DB)
}

// StoreToken encrypts and saves a token
func StoreToken(project, tag, rawToken string) error {
	// Analyze
	meta := utils.AnalyzeToken(rawToken)

	//  Encrypt
	encrypted, err := crypto.Encrypt(rawToken)
	if err != nil {
		return err
	}

	// Upsert (Insert or Update if exists)
	query := `
	INSERT INTO tokens (project, tag, encrypted_value, token_type, expires_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(project, tag) DO UPDATE SET
		encrypted_value = excluded.encrypted_value,
		token_type = excluded.token_type,
		expires_at = excluded.expires_at,
		created_at = CURRENT_TIMESTAMP;
	`
	_, err = DB.Exec(query, project, tag, encrypted, meta.Type, meta.ExpiresAt)
	return err
}

// FetchToken retrieves and decrypts a token
func FetchToken(project, tag string) (string, *utils.TokenMetadata, error) {
	var encrypted string
	var meta utils.TokenMetadata

	query := `SELECT encrypted_value, token_type, expires_at FROM tokens WHERE project = ? AND tag = ?`

	err := DB.QueryRow(query, project, tag).Scan(&encrypted, &meta.Type, &meta.ExpiresAt)
	if err == sql.ErrNoRows {
		return "", nil, fmt.Errorf("token not found")
	} else if err != nil {
		return "", nil, err
	}

	// Decrypt
	decrypted, err := crypto.Decrypt(encrypted)
	if err != nil {
		return "", nil, fmt.Errorf("decryption failed: %v", err)
	}

	return decrypted, &meta, nil
}
