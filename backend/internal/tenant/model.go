package tenant

import "time"

// Tenant represents a tenant/organization in the multi-tenant system.
// Each tenant has its own InvGate credentials and branding configuration.
type Tenant struct {
	ID   string `gorm:"type:char(36);primaryKey" json:"id"`
	Name string `gorm:"size:255;not null" json:"name"`
	Slug string `gorm:"size:100;not null;uniqueIndex" json:"slug"`

	// InvGate Configuration - each tenant has their own credentials
	InvGateCompanyID  int    `gorm:"column:invgate_company_id;not null" json:"invgate_company_id"`
	InvGateGroupID    int    `gorm:"column:invgate_group_id;not null" json:"invgate_group_id"`
	InvGateLocationID int    `gorm:"column:invgate_location_id;not null" json:"invgate_location_id"`
	InvGateBaseURL    string `gorm:"column:invgate_base_url;size:255;not null" json:"invgate_base_url"`
	InvGateUsername   string `gorm:"column:invgate_username;size:255;not null" json:"invgate_username"`
	InvGatePassword   string `gorm:"column:invgate_password;size:255;not null" json:"-"` // Never expose in JSON

	// Email Configuration
	EmailDomain string `gorm:"column:email_domain;size:255" json:"email_domain,omitempty"`
	EmailSender string `gorm:"column:email_sender;size:255" json:"email_sender,omitempty"`

	// Branding
	LogoURL      string `gorm:"column:logo_url;size:255" json:"logo_url,omitempty"`
	PrimaryColor string `gorm:"column:primary_color;size:7;default:#1976D2" json:"primary_color"`

	// Status
	IsActive  bool      `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (Tenant) TableName() string {
	return "tenants"
}

// CreateTenantRequest is the DTO for creating a new tenant
type CreateTenantRequest struct {
	Name              string `json:"name" binding:"required"`
	Slug              string `json:"slug" binding:"required"`
	InvGateCompanyID  int    `json:"invgate_company_id" binding:"required"`
	InvGateGroupID    int    `json:"invgate_group_id" binding:"required"`
	InvGateLocationID int    `json:"invgate_location_id" binding:"required"`
	InvGateBaseURL    string `json:"invgate_base_url" binding:"required"`
	InvGateUsername   string `json:"invgate_username" binding:"required"`
	InvGatePassword   string `json:"invgate_password" binding:"required"`
	EmailDomain       string `json:"email_domain,omitempty"`
	EmailSender       string `json:"email_sender,omitempty"`
	LogoURL           string `json:"logo_url,omitempty"`
	PrimaryColor      string `json:"primary_color,omitempty"`
}

// UpdateTenantRequest is the DTO for updating a tenant
type UpdateTenantRequest struct {
	Name              string `json:"name,omitempty"`
	InvGateCompanyID  *int   `json:"invgate_company_id,omitempty"`
	InvGateGroupID    *int   `json:"invgate_group_id,omitempty"`
	InvGateLocationID *int   `json:"invgate_location_id,omitempty"`
	InvGateBaseURL    string `json:"invgate_base_url,omitempty"`
	InvGateUsername   string `json:"invgate_username,omitempty"`
	InvGatePassword   string `json:"invgate_password,omitempty"`
	EmailDomain       string `json:"email_domain,omitempty"`
	EmailSender       string `json:"email_sender,omitempty"`
	LogoURL           string `json:"logo_url,omitempty"`
	PrimaryColor      string `json:"primary_color,omitempty"`
	IsActive          *bool  `json:"is_active,omitempty"`
}

// TenantPublicInfo is the public-facing tenant info (for frontend branding)
type TenantPublicInfo struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	LogoURL      string `json:"logo_url,omitempty"`
	PrimaryColor string `json:"primary_color"`
}

// ToPublicInfo converts a Tenant to TenantPublicInfo
func (t *Tenant) ToPublicInfo() TenantPublicInfo {
	return TenantPublicInfo{
		ID:           t.ID,
		Name:         t.Name,
		Slug:         t.Slug,
		LogoURL:      t.LogoURL,
		PrimaryColor: t.PrimaryColor,
	}
}
