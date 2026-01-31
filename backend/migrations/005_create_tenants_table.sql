-- Migration: Create tenants table
-- This table stores all tenant configurations including InvGate credentials

CREATE TABLE IF NOT EXISTS tenants (
    id CHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    
    -- InvGate credentials per tenant
    invgate_company_id INT NOT NULL,
    invgate_group_id INT NOT NULL,
    invgate_location_id INT NOT NULL,
    invgate_base_url VARCHAR(255) NOT NULL,
    invgate_username VARCHAR(255) NOT NULL,
    invgate_password VARCHAR(255) NOT NULL,
    
    -- Email configuration
    email_domain VARCHAR(255),
    email_sender VARCHAR(255),
    
    -- Branding
    logo_url VARCHAR(255),
    primary_color VARCHAR(7) DEFAULT '#1976D2',
    
    -- Status
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_slug (slug),
    INDEX idx_active (is_active)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert default tenant for existing data migration
INSERT INTO tenants (
    id, name, slug, 
    invgate_company_id, invgate_group_id, invgate_location_id,
    invgate_base_url, invgate_username, invgate_password,
    is_active
) VALUES (
    'default-tenant-id', 
    'Default Company', 
    'default',
    135, 134, 136,
    'https://support.armmada.id/api/v1/',
    'armmadaweb',
    'j8f2yDzuVhYI4eG67hbsbck0',
    TRUE
);
