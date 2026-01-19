package utils

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

type TokenMetadata struct {
	Type      string `json:"type"`       // "Bearer", "API-Key", "Basic"
	ExpiresAt int64  `json:"expires_at"` // Unix Timestamp
	IsJWT     bool   `json:"is_jwt"`
	Subject   string `json:"sub,omitempty"` // User ID or Email
}

// AnalyzeToken inspects a token string to extract useful metadata.
// It does NOT verify the signature, it only reads the payload.
func AnalyzeToken(rawToken string) TokenMetadata {
	meta := TokenMetadata{
		Type: "Unknown",
		// Default to 24 hours if we can't find an expiry date
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		IsJWT:     false,
	}

	//Check if it looks like a JWT (Header.Payload.Signature)
	parts := strings.Split(rawToken, ".")
	if len(parts) == 3 {
		meta.Type = "Bearer"
		meta.IsJWT = true

		// 2. Decode the Payload
		// JWTs use Base64 Raw URL Encoding (no padding, url-safe)
		payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
		if err == nil {
			var claims map[string]interface{}
			if err := json.Unmarshal(payloadBytes, &claims); err == nil {
				//Extract Expiry
				if exp, ok := claims["exp"].(float64); ok {
					meta.ExpiresAt = int64(exp)
				}

				// Extract Subject (sub) - usually the user ID
				if sub, ok := claims["sub"].(string); ok {
					meta.Subject = sub
				}
			}
		}
	} else if strings.HasPrefix(rawToken, "sk_") {
		// Stripe/API Key detection logic
		meta.Type = "Stripe-Key"
	} else if strings.HasPrefix(rawToken, "ghp_") {
		meta.Type = "GitHub-Token"
	}

	return meta
}
