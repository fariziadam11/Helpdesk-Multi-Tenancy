package user

import "time"

// ResetToken represents a password reset token
type ResetToken struct {
	ID        string    `gorm:"type:char(36);primaryKey;default:(UUID())"`
	TenantID  string    `gorm:"type:char(36);not null;index:idx_reset_tokens_tenant_id"`
	UserID    string    `gorm:"type:char(36);not null"`
	Token     string    `gorm:"size:255;not null;uniqueIndex"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Foreign key relationship
	Tenant Tenant `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for GORM
func (ResetToken) TableName() string {
	return "password_reset_tokens"
}
