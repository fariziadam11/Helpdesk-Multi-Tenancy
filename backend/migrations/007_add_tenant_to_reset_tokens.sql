-- Migration: Add tenant_id to reset_tokens table
-- This enables multi-tenant data isolation for password reset tokens

-- Step 1: Add tenant_id column
ALTER TABLE reset_tokens ADD COLUMN tenant_id VARCHAR(36) NULL AFTER id;

-- Step 2: Populate existing tokens with default tenant
UPDATE reset_tokens SET tenant_id = 'default-tenant-id' WHERE tenant_id IS NULL;

-- Step 3: Make NOT NULL
ALTER TABLE reset_tokens MODIFY COLUMN tenant_id VARCHAR(36) NOT NULL;

-- Step 4: Add index
ALTER TABLE reset_tokens ADD INDEX idx_reset_tokens_tenant_id (tenant_id);

-- Step 5: Add foreign key
ALTER TABLE reset_tokens ADD CONSTRAINT fk_reset_tokens_tenant 
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE;
