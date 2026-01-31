package tenant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"werk-ticketing/internal/errors"
	"werk-ticketing/internal/response"
)

// Handler handles HTTP requests for tenant management
type Handler struct {
	repo Repository
}

// NewHandler creates a new tenant handler
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// Create handles POST /admin/tenants
func (h *Handler) Create(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid request body")
		return
	}

	// Check if slug already exists
	existing, err := h.repo.FindBySlug(c.Request.Context(), req.Slug)
	if err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to check slug")
		return
	}
	if existing != nil {
		response.ErrorWithCode(c, http.StatusConflict, errors.ErrCodeConflict, "tenant with this slug already exists")
		return
	}

	tenant := &Tenant{
		ID:                uuid.New().String(),
		Name:              req.Name,
		Slug:              req.Slug,
		InvGateCompanyID:  req.InvGateCompanyID,
		InvGateGroupID:    req.InvGateGroupID,
		InvGateLocationID: req.InvGateLocationID,
		InvGateBaseURL:    req.InvGateBaseURL,
		InvGateUsername:   req.InvGateUsername,
		InvGatePassword:   req.InvGatePassword,
		EmailDomain:       req.EmailDomain,
		EmailSender:       req.EmailSender,
		LogoURL:           req.LogoURL,
		PrimaryColor:      req.PrimaryColor,
		IsActive:          true,
	}

	if tenant.PrimaryColor == "" {
		tenant.PrimaryColor = "#1976D2"
	}

	if err := h.repo.Create(c.Request.Context(), tenant); err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to create tenant")
		return
	}

	response.Success(c, http.StatusCreated, tenant.ToPublicInfo())
}

// List handles GET /admin/tenants
func (h *Handler) List(c *gin.Context) {
	tenants, err := h.repo.FindAllIncludingInactive(c.Request.Context())
	if err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to list tenants")
		return
	}

	response.Success(c, http.StatusOK, tenants)
}

// GetByID handles GET /admin/tenants/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant id is required")
		return
	}

	tenant, err := h.repo.FindByIDIncludingInactive(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to get tenant")
		return
	}
	if tenant == nil {
		response.ErrorWithCode(c, http.StatusNotFound, errors.ErrCodeNotFound, "tenant not found")
		return
	}

	response.Success(c, http.StatusOK, tenant)
}

// Update handles PUT /admin/tenants/:id
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant id is required")
		return
	}

	var req UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid request body")
		return
	}

	tenant, err := h.repo.FindByIDIncludingInactive(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to get tenant")
		return
	}
	if tenant == nil {
		response.ErrorWithCode(c, http.StatusNotFound, errors.ErrCodeNotFound, "tenant not found")
		return
	}

	// Update fields if provided
	if req.Name != "" {
		tenant.Name = req.Name
	}
	if req.InvGateCompanyID != nil {
		tenant.InvGateCompanyID = *req.InvGateCompanyID
	}
	if req.InvGateGroupID != nil {
		tenant.InvGateGroupID = *req.InvGateGroupID
	}
	if req.InvGateLocationID != nil {
		tenant.InvGateLocationID = *req.InvGateLocationID
	}
	if req.InvGateBaseURL != "" {
		tenant.InvGateBaseURL = req.InvGateBaseURL
	}
	if req.InvGateUsername != "" {
		tenant.InvGateUsername = req.InvGateUsername
	}
	if req.InvGatePassword != "" {
		tenant.InvGatePassword = req.InvGatePassword
	}
	if req.EmailDomain != "" {
		tenant.EmailDomain = req.EmailDomain
	}
	if req.EmailSender != "" {
		tenant.EmailSender = req.EmailSender
	}
	if req.LogoURL != "" {
		tenant.LogoURL = req.LogoURL
	}
	if req.PrimaryColor != "" {
		tenant.PrimaryColor = req.PrimaryColor
	}
	if req.IsActive != nil {
		tenant.IsActive = *req.IsActive
	}

	if err := h.repo.Update(c.Request.Context(), tenant); err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to update tenant")
		return
	}

	response.Success(c, http.StatusOK, tenant)
}

// Delete handles DELETE /admin/tenants/:id
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant id is required")
		return
	}

	tenant, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to get tenant")
		return
	}
	if tenant == nil {
		response.ErrorWithCode(c, http.StatusNotFound, errors.ErrCodeNotFound, "tenant not found")
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to delete tenant")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "tenant deleted successfully"})
}

// UpdateStatus handles PATCH /admin/tenants/:id/status
func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant id is required")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid request body")
		return
	}

	tenant, err := h.repo.FindByID(c.Request.Context(), id)
	if err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to get tenant")
		return
	}
	if tenant == nil {
		response.ErrorWithCode(c, http.StatusNotFound, errors.ErrCodeNotFound, "tenant not found")
		return
	}

	tenant.IsActive = req.IsActive
	if err := h.repo.Update(c.Request.Context(), tenant); err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to update tenant status")
		return
	}

	response.Success(c, http.StatusOK, tenant.ToPublicInfo())
}

// GetPublicInfo handles GET /tenants/:slug/info (public endpoint for branding)
func (h *Handler) GetPublicInfo(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant slug is required")
		return
	}

	tenant, err := h.repo.FindBySlug(c.Request.Context(), slug)
	if err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to get tenant")
		return
	}
	if tenant == nil {
		response.ErrorWithCode(c, http.StatusNotFound, errors.ErrCodeNotFound, "tenant not found")
		return
	}

	response.Success(c, http.StatusOK, tenant.ToPublicInfo())
}
