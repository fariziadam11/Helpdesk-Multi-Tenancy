package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"werk-ticketing/internal/errors"
	"werk-ticketing/internal/response"
	"werk-ticketing/internal/tenant"
)

const tenantIDKey = "tenant_id"

// TenantIdentificationMethod defines how tenant is identified from request
type TenantIdentificationMethod int

const (
	// TenantFromHeader identifies tenant from X-Tenant-ID header
	TenantFromHeader TenantIdentificationMethod = iota
	// TenantFromSubdomain identifies tenant from subdomain
	TenantFromSubdomain
	// TenantFromQuery identifies tenant from query parameter
	TenantFromQuery
)

// WithTenant ensures request has a valid tenant identified.
// It validates tenant exists and is active before proceeding.
func WithTenant(tenantRepo tenant.Repository, method TenantIdentificationMethod) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenantID string

		switch method {
		case TenantFromHeader:
			tenantID = c.GetHeader("X-Tenant-ID")
		case TenantFromSubdomain:
			tenantID = extractTenantFromSubdomain(c.Request.Host)
		case TenantFromQuery:
			tenantID = c.Query("tenant_id")
		}

		if tenantID == "" {
			response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant not identified")
			c.Abort()
			return
		}

		// Validate tenant exists and is active
		t, err := tenantRepo.FindByID(c.Request.Context(), tenantID)
		if err != nil {
			response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to validate tenant")
			c.Abort()
			return
		}
		if t == nil {
			response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant not found")
			c.Abort()
			return
		}
		if !t.IsActive {
			response.ErrorWithCode(c, http.StatusForbidden, errors.ErrCodeForbidden, "tenant is not active")
			c.Abort()
			return
		}

		// Set tenant ID in context for use by handlers
		c.Set(tenantIDKey, tenantID)
		c.Next()
	}
}

// WithTenantBySlug identifies tenant by slug (useful for subdomain or path-based routing)
func WithTenantBySlug(tenantRepo tenant.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := extractTenantFromSubdomain(c.Request.Host)
		if slug == "" {
			slug = c.GetHeader("X-Tenant-Slug")
		}

		if slug == "" {
			response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant not identified")
			c.Abort()
			return
		}

		// Validate tenant exists and is active
		t, err := tenantRepo.FindBySlug(c.Request.Context(), slug)
		if err != nil {
			response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to validate tenant")
			c.Abort()
			return
		}
		if t == nil {
			response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant not found")
			c.Abort()
			return
		}
		if !t.IsActive {
			response.ErrorWithCode(c, http.StatusForbidden, errors.ErrCodeForbidden, "tenant is not active")
			c.Abort()
			return
		}

		// Set tenant ID in context for use by handlers
		c.Set(tenantIDKey, t.ID)
		c.Next()
	}
}

// extractTenantFromSubdomain extracts tenant slug from subdomain
// e.g., "tenant1.app.example.com" -> "tenant1"
func extractTenantFromSubdomain(host string) string {
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	parts := strings.Split(host, ".")
	if len(parts) >= 3 {
		// Assumes format: subdomain.domain.tld
		return parts[0]
	}
	return ""
}

// GetTenantID extracts the tenant ID from the request context.
func GetTenantID(c *gin.Context) string {
	if tenantID, ok := c.Get(tenantIDKey); ok {
		if s, ok := tenantID.(string); ok {
			return s
		}
	}
	return ""
}
