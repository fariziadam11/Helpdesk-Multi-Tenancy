package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"werk-ticketing/internal/auth"
	"werk-ticketing/internal/constants"
	"werk-ticketing/internal/middleware"
	"werk-ticketing/internal/tenant"
	"werk-ticketing/internal/ticket"
	"werk-ticketing/internal/upload"
	"werk-ticketing/internal/user"
)

// Router holds all route dependencies
type Router struct {
	authHandler   *auth.Handler
	ticketHandler *ticket.Handler
	userHandler   *user.Handler
	tenantHandler *tenant.Handler
	authService   auth.Service
	tenantRepo    tenant.Repository
	logger        *logrus.Logger
}

// NewRouter creates a new router instance
func NewRouter(
	authHandler *auth.Handler,
	ticketHandler *ticket.Handler,
	userHandler *user.Handler,
	tenantHandler *tenant.Handler,
	authService auth.Service,
	tenantRepo tenant.Repository,
	logger *logrus.Logger,
) *Router {
	return &Router{
		authHandler:   authHandler,
		ticketHandler: ticketHandler,
		userHandler:   userHandler,
		tenantHandler: tenantHandler,
		authService:   authService,
		tenantRepo:    tenantRepo,
		logger:        logger,
	}
}

// SetupRoutes configures all application routes
func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.New()

	// Global middleware (order matters!)
	router.Use(
		gin.Logger(),
		middleware.Recover(r.logger),
		middleware.CORS(),
		middleware.SecurityHeaders(),
		middleware.RateLimit(),
	)

	// Set max request size
	router.MaxMultipartMemory = constants.MaxRequestSize

	// API versioning: /api/v1
	apiV1 := router.Group("/api/v1")

	// Public routes (no tenant or auth required)
	// These are accessible without any authentication or tenant context
	r.setupAuthRoutes(apiV1)         // Login, register, forgot-password, etc.
	r.setupPublicTenantRoutes(apiV1) // Tenant public info endpoint

	// Public reference data endpoints (no auth, but may need tenant context in future)
	apiV1.GET("/categories", r.ticketHandler.GetCategories)
	apiV1.GET("/ticket-meta", r.ticketHandler.GetMeta)
	apiV1.GET("/statuses", r.ticketHandler.GetStatuses)

	// Articles endpoint (public, no auth required)
	articleRoutes := apiV1.Group("/articles")
	articleRoutes.Use(middleware.ArticleRateLimit())
	{
		articleRoutes.GET("", r.ticketHandler.GetArticlesByCategory)
	}

	// Protected routes (require tenant context and authentication)
	// Apply tenant middleware to these routes
	protectedRoutes := apiV1.Group("")
	protectedRoutes.Use(middleware.WithTenant(r.tenantRepo, middleware.TenantFromHeader))
	{
		r.setupTicketRoutes(protectedRoutes)
		r.setupAdminTenantRoutes(protectedRoutes) // Admin tenant CRUD routes

		// User endpoint (proxy to InvGate user API, requires auth)
		userRoutes := protectedRoutes.Group("/users")
		userRoutes.Use(middleware.WithAuth(r.authService))
		{
			// GET /api/users/:id - Get InvGate user detail by ID
			userRoutes.GET("/:id", r.ticketHandler.GetInvGateUser)
			// PUT /api/users/profile - Update current user's profile
			userRoutes.PUT("/profile", r.userHandler.UpdateProfile)
		}
	}

	// Upload Endpoint (Protected)
	protectedRoutes.POST("/upload", func(c *gin.Context) {
		upload.NewHandler().UploadFile(c)
	})

	// Health check endpoint (no versioning, no tenant required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "werk-ticketing-backend",
		})
	})

	// Serve uploaded files (Static)
	router.Static("/uploads", "./uploads")

	return router
}
