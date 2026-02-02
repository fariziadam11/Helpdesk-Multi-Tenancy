package tenant

import "errors"

var (
	// ErrTenantNotFound is returned when a tenant is not found
	ErrTenantNotFound = errors.New("tenant not found")

	// ErrTenantInactive is returned when a tenant is inactive
	ErrTenantInactive = errors.New("tenant is inactive")

	// ErrTenantSlugExists is returned when a tenant slug already exists
	ErrTenantSlugExists = errors.New("tenant slug already exists")
)
