package router

import "github.com/gin-gonic/gin"

// setupAuthRoutes configures authentication routes
func (r *Router) setupAuthRoutes(api *gin.RouterGroup) {
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", r.authHandler.Register)
		authGroup.POST("/login", r.authHandler.Login)
		authGroup.POST("/refresh", r.authHandler.RefreshToken)
		authGroup.POST("/forgot-password", r.authHandler.ForgotPassword)
		authGroup.POST("/reset-password", r.authHandler.ResetPassword)

		// Protected auth routes (require authentication)
		authGroup.POST("/revoke", r.authHandler.RevokeToken)
	}
}
