package main

import (
	"log"

	"werk-ticketing/internal/config"
	"werk-ticketing/internal/database"
	"werk-ticketing/internal/tenant"

	"github.com/google/uuid"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	// Create default tenant
	defaultTenant := &tenant.Tenant{
		ID:                uuid.New().String(),
		Name:              "Default Company",
		Slug:              "default",
		InvGateCompanyID:  135,
		InvGateGroupID:    134,
		InvGateLocationID: 136,
		InvGateBaseURL:    "https://support.armmada.id/api/v1/",
		InvGateUsername:   "armmadaweb",
		InvGatePassword:   "j8f2yDzuVhYI4eG67hbsbck0",
		PrimaryColor:      "#1976D2",
		IsActive:          true,
	}

	// Check if tenant already exists
	var existingTenant tenant.Tenant
	result := db.Where("slug = ?", "default").First(&existingTenant)

	if result.Error != nil {
		// Tenant doesn't exist, create it
		if err := db.Create(defaultTenant).Error; err != nil {
			log.Fatalf("failed to create default tenant: %v", err)
		}
		log.Printf("✅ Default tenant created with ID: %s", defaultTenant.ID)
	} else {
		log.Printf("ℹ️  Default tenant already exists with ID: %s", existingTenant.ID)
		defaultTenant.ID = existingTenant.ID
	}

	// Update existing users to use default tenant
	result = db.Exec("UPDATE users SET tenant_id = ? WHERE tenant_id IS NULL OR tenant_id = ''", defaultTenant.ID)
	if result.Error != nil {
		log.Printf("⚠️  Warning: failed to update users: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("✅ Updated %d users with default tenant", result.RowsAffected)
	}

	// Update existing reset tokens to use default tenant
	result = db.Exec("UPDATE password_reset_tokens SET tenant_id = ? WHERE tenant_id IS NULL OR tenant_id = ''", defaultTenant.ID)
	if result.Error != nil {
		log.Printf("⚠️  Warning: failed to update reset tokens: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("✅ Updated %d reset tokens with default tenant", result.RowsAffected)
	}

	log.Println("✅ Seeding completed successfully!")
}
