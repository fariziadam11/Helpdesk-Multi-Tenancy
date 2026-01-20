package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// ResetToken represents a password reset token
type ResetToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// GenerateResetToken creates a secure random token
func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// IsExpired checks if the token has expired
func (t *ResetToken) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}
