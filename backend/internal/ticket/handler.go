package ticket

import "github.com/gin-gonic/gin"

// Handler wires ticket service with HTTP endpoints.
type Handler struct {
	service Service
}

// NewHandler creates ticket handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// getTenantID extracts tenant ID from gin context (set by tenant middleware)
func getTenantID(c *gin.Context) (string, bool) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		return "", false
	}
	return tenantID.(string), true
}
