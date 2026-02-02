package router

import "github.com/gin-gonic/gin"

// setupPublicTenantRoutes configures public tenant routes (no auth or tenant middleware required)
func (r *Router) setupPublicTenantRoutes(api *gin.RouterGroup) {
	// Public endpoint for getting tenant branding info (no auth required)
	// Used by frontend to display tenant-specific branding
	api.GET("/tenants/:slug/info", r.tenantHandler.GetPublicInfo)
}

// setupAdminTenantRoutes configures admin tenant management routes (requires tenant context and auth)
func (r *Router) setupAdminTenantRoutes(api *gin.RouterGroup) {
	// Admin routes for tenant management (require authentication)
	adminRoutes := api.Group("/admin/tenants")
	// TODO: Add admin role check middleware here
	// adminRoutes.Use(middleware.WithAdminRole(r.authService))
	{
		// POST /admin/tenants - Create a new tenant
		adminRoutes.POST("", r.tenantHandler.Create)

		// GET /admin/tenants - List all tenants
		adminRoutes.GET("", r.tenantHandler.List)

		// GET /admin/tenants/:id - Get tenant by ID
		adminRoutes.GET("/:id", r.tenantHandler.GetByID)

		// PUT /admin/tenants/:id - Update tenant
		adminRoutes.PUT("/:id", r.tenantHandler.Update)

		// DELETE /admin/tenants/:id - Delete (soft) tenant
		adminRoutes.DELETE("/:id", r.tenantHandler.Delete)

		// PATCH /admin/tenants/:id/status - Activate/Deactivate tenant
		adminRoutes.PATCH("/:id/status", r.tenantHandler.UpdateStatus)
	}
}
