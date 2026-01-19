package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kartikgujjar22/tokenvault/internal/crypto"
	"github.com/kartikgujjar22/tokenvault/internal/utils"
)

// RunMigrations checks the DB state and upgrades if necessary
func RunMigrations(db *sql.DB) error {
	// Check if the v2 table exists
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='tokens';`
	row := db.QueryRow(query)
	var name string
	err := row.Scan(&name)

	if err == sql.ErrNoRows {

		fmt.Println("Initializing V2 Database Schema...")
		return createV2Table(db)
	}
	_, err = db.Query("SELECT encrypted_value FROM tokens LIMIT 1;")
	if err != nil {
		fmt.Println("Detected V1 Database. Starting Migration to V2 (Encryption)...")
		return migrateV1ToV2(db)
	}

	return nil
}

func createV2Table(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project TEXT NOT NULL,
		tag TEXT DEFAULT 'default',
		encrypted_value TEXT NOT NULL,
		token_type TEXT,
		expires_at INTEGER,
		meta_json TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(project, tag)
	);`
	_, err := db.Exec(query)
	return err
}

func migrateV1ToV2(db *sql.DB) error {
	// Rename old table to backup
	_, err := db.Exec("ALTER TABLE tokens RENAME TO tokens_v1_backup;")
	if err != nil {
		return fmt.Errorf("failed to rename old table: %v", err)
	}

	// Create the NEW V2 table
	if err := createV2Table(db); err != nil {
		return fmt.Errorf("failed to create v2 table: %v", err)
	}

	//Read all old data
	rows, err := db.Query("SELECT project_name, token_value FROM tokens_v1_backup;")
	if err != nil {
		return err
	}
	defer rows.Close()

	//Process and Migrate each row
	count := 0
	for rows.Next() {
		var project, rawToken string
		if err := rows.Scan(&project, &rawToken); err != nil {
			log.Printf("Skipping corrupt row: %v", err)
			continue
		}

		// Inspect (Get Metadata)
		meta := utils.AnalyzeToken(rawToken)

		//Encrypt
		encrypted, err := crypto.Encrypt(rawToken)
		if err != nil {
			log.Printf("Failed to encrypt token for %s: %v", project, err)
			continue
		}

		// Insert into New Table
		_, err = db.Exec(`
			INSERT INTO tokens (project, tag, encrypted_value, token_type, expires_at, meta_json)
			VALUES (?, ?, ?, ?, ?, ?)`,
			project, "default", encrypted, meta.Type, meta.ExpiresAt, "{}",
		)
		if err != nil {
			log.Printf("Failed to insert v2 row for %s: %v", project, err)
		} else {
			count++
		}
	}

	fmt.Printf("Migration Complete: Upgraded %d tokens to V2 Encryption.\n", count)

	return nil
}
