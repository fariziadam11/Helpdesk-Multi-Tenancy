package user

import (
	"time"
)

// ResetToken represents a password reset token in the database
type ResetToken struct {
	ID        string    `gorm:"type:char(36);primaryKey;default:(UUID())"`
	UserID    string    `gorm:"type:char(36);not null"`
	Token     string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// TableName specifies the table name for ResetToken
func (ResetToken) TableName() string {
	return "password_reset_tokens"
}
