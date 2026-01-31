package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"werk-ticketing/internal/errors"
	"werk-ticketing/internal/response"
	"werk-ticketing/internal/tenant"
)

const (
	tenantIDKey  = "tenant_id"
	tenantObjKey = "tenant"
	cacheExpiry  = 5 * time.Minute
)

// TenantIdentificationMethod defines how tenant is identified from request
type TenantIdentificationMethod int

const (
	// TenantFromHeader identifies tenant from X-Tenant-ID header
	TenantFromHeader TenantIdentificationMethod = iota
	// TenantFromSubdomain identifies tenant from subdomain
	TenantFromSubdomain
	// TenantFromQuery identifies tenant from query parameter
	TenantFromQuery
	// TenantAutoDetect tries header first, then subdomain, then query
	TenantAutoDetect
)

// tenantCache stores recently looked up tenants
var tenantCache = &sync.Map{}

type cachedTenant struct {
	tenant    *tenant.Tenant
	expiresAt time.Time
}

// WithTenant ensures request has a valid tenant identified.
// It validates tenant exists and is active before proceeding.
func WithTenant(tenantRepo tenant.Repository, method TenantIdentificationMethod) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenantID string
		var slug string
		useSlug := false

		switch method {
		case TenantFromHeader:
			tenantID = c.GetHeader("X-Tenant-ID")
		case TenantFromSubdomain:
			slug = extractTenantFromSubdomain(c.Request.Host)
			useSlug = true
		case TenantFromQuery:
			tenantID = c.Query("tenant_id")
		case TenantAutoDetect:
			// Try header first
			tenantID = c.GetHeader("X-Tenant-ID")
			if tenantID == "" {
				// Try subdomain
				slug = extractTenantFromSubdomain(c.Request.Host)
				if slug != "" {
					useSlug = true
				} else {
					// Try query
					tenantID = c.Query("tenant_id")
				}
			}
		}

		if tenantID == "" && slug == "" {
			response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant not identified")
			c.Abort()
			return
		}

		// Try cache first
		var t *tenant.Tenant
		var err error

		cacheKey := tenantID
		if useSlug {
			cacheKey = "slug:" + slug
		}

		if cached, ok := tenantCache.Load(cacheKey); ok {
			ct := cached.(*cachedTenant)
			if time.Now().Before(ct.expiresAt) {
				t = ct.tenant
			} else {
				tenantCache.Delete(cacheKey)
			}
		}

		// Fetch from DB if not cached
		if t == nil {
			if useSlug {
				t, err = tenantRepo.FindBySlug(c.Request.Context(), slug)
			} else {
				t, err = tenantRepo.FindByID(c.Request.Context(), tenantID)
			}

			if err != nil {
				response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to validate tenant")
				c.Abort()
				return
			}

			if t != nil {
				// Cache the result
				tenantCache.Store(cacheKey, &cachedTenant{
					tenant:    t,
					expiresAt: time.Now().Add(cacheExpiry),
				})
			}
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

		// Set tenant ID and object in context for use by handlers
		c.Set(tenantIDKey, t.ID)
		c.Set(tenantObjKey, t)
		c.Next()
	}
}

// WithTenantBySlug identifies tenant by slug (useful for subdomain or path-based routing)
func WithTenantBySlug(tenantRepo tenant.Repository) gin.HandlerFunc {
	return WithTenant(tenantRepo, TenantFromSubdomain)
}

// extractTenantFromSubdomain extracts tenant slug from subdomain
// e.g., "tenant1.app.example.com" -> "tenant1"
// Handles various formats:
// - tenant1.app.example.com -> tenant1
// - tenant1.localhost -> tenant1
// - localhost -> ""
func extractTenantFromSubdomain(host string) string {
	// Remove port if present
	if idx := strings.Index(host, ":"); idx != -1 {
		host = host[:idx]
	}

	// Skip localhost without subdomain
	if host == "localhost" {
		return ""
	}

	parts := strings.Split(host, ".")

	// Handle localhost with subdomain (e.g., tenant1.localhost)
	if len(parts) == 2 && parts[1] == "localhost" {
		return parts[0]
	}

	// Handle regular domain (e.g., tenant1.app.example.com)
	if len(parts) >= 3 {
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

// GetTenant extracts the full tenant object from the request context.
func GetTenant(c *gin.Context) *tenant.Tenant {
	if t, ok := c.Get(tenantObjKey); ok {
		if tenantObj, ok := t.(*tenant.Tenant); ok {
			return tenantObj
		}
	}
	return nil
}

// InvalidateTenantCache removes a tenant from cache (call after update)
func InvalidateTenantCache(tenantID string, slug string) {
	if tenantID != "" {
		tenantCache.Delete(tenantID)
	}
	if slug != "" {
		tenantCache.Delete("slug:" + slug)
	}
}
