package database

import "time"

type Token struct {
	ID             int64     `json:"id"`
	Project        string    `json:"project"`
	Tag            string    `json:"tag"`        // e.g. "admin", "default"
	EncryptedValue string    `json:"-"`          // Never return this in JSON
	TokenType      string    `json:"token_type"` // "Bearer", "API-Key"
	ExpiresAt      int64     `json:"expires_at"` // Unix Timestamp
	MetaJSON       string    `json:"meta_json"`  // Extra metadata
	CreatedAt      time.Time `json:"created_at"`
}
