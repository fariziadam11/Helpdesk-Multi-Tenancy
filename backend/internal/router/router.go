package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"werk-ticketing/internal/auth"
	"werk-ticketing/internal/constants"
	"werk-ticketing/internal/middleware"
	"werk-ticketing/internal/tenant"
	"werk-ticketing/internal/ticket"
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

	// Apply tenant middleware to all API routes
	// Tenant is identified from X-Tenant-ID header
	apiV1.Use(middleware.WithTenant(r.tenantRepo, middleware.TenantFromHeader))

	// Setup route groups
	r.setupAuthRoutes(apiV1)
	r.setupTicketRoutes(apiV1)
	r.setupTenantRoutes(apiV1)

	// User endpoint (proxy to InvGate user API, requires auth)
	userRoutes := apiV1.Group("/users")
	userRoutes.Use(middleware.WithAuth(r.authService))
	{
		// GET /api/users/:id - Get InvGate user detail by ID
		userRoutes.GET("/:id", r.ticketHandler.GetInvGateUser)
		// PUT /api/users/profile - Update current user's profile
		userRoutes.PUT("/profile", r.userHandler.UpdateProfile)
	}

	// Categories endpoint (public, no auth required for reference data)
	apiV1.GET("/categories", r.ticketHandler.GetCategories)
	// Ticket meta endpoint (public, no auth required for reference data)
	apiV1.GET("/ticket-meta", r.ticketHandler.GetMeta)
	// Statuses endpoint (public, no auth required for reference data)
	apiV1.GET("/statuses", r.ticketHandler.GetStatuses)
	// Articles endpoint (public, no auth required for reference data)
	// Use higher rate limit for articles since landing page fetches multiple categories
	articleRoutes := apiV1.Group("/articles")
	articleRoutes.Use(middleware.ArticleRateLimit())
	{
		articleRoutes.GET("", r.ticketHandler.GetArticlesByCategory)
	}

	// Health check endpoint (no versioning, no tenant required)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "werk-ticketing-backend",
		})
	})

	return router
}
