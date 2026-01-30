-- Migration: Add tenant_id to users table
-- This enables multi-tenant data isolation for users

-- Step 1: Add tenant_id column (nullable first for migration)
ALTER TABLE users ADD COLUMN tenant_id VARCHAR(36) NULL AFTER id;

-- Step 2: Populate existing users with default tenant
UPDATE users SET tenant_id = 'default-tenant-id' WHERE tenant_id IS NULL;

-- Step 3: Make NOT NULL after population
ALTER TABLE users MODIFY COLUMN tenant_id VARCHAR(36) NOT NULL;

-- Step 4: Add index for performance
ALTER TABLE users ADD INDEX idx_users_tenant_id (tenant_id);

-- Step 5: Add foreign key constraint
ALTER TABLE users ADD CONSTRAINT fk_users_tenant 
    FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE RESTRICT;

-- Step 6: Update unique constraint for email to be per-tenant
-- First drop existing unique index on email if exists
ALTER TABLE users DROP INDEX IF EXISTS email;
ALTER TABLE users DROP INDEX IF EXISTS idx_email;

-- Add composite unique index for tenant + email
ALTER TABLE users ADD UNIQUE INDEX idx_users_tenant_email (tenant_id, email);
