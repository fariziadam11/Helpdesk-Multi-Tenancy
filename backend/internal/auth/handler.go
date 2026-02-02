package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"werk-ticketing/internal/errors"
	"werk-ticketing/internal/response"
)

// Handler exposes HTTP handlers for auth routes.
type Handler struct {
	service Service
}

// NewHandler wires auth service into http handler.
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

// Register handles POST /auth/register
// Register handles POST /auth/register
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid JSON body")
		return
	}

	// For registration, tenant ID must be provided in the body
	if req.TenantID == "" {
		// Try to get from context as fallback (though unlikely for public route)
		if tid, ok := getTenantID(c); ok {
			req.TenantID = tid
		} else {
			response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant_id is required")
			return
		}
	}

	resp, err := h.service.Register(c.Request.Context(), req.TenantID, req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.AppError(c, appErr)
		} else {
			response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, err.Error())
		}
		return
	}

	response.Write(c, http.StatusCreated, resp)
}

// Login handles POST /auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid JSON body")
		return
	}

	// Login without requiring tenant ID in request
	// The service will find the user and return their tenant information
	resp, err := h.service.Login(c.Request.Context(), "", req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.AppError(c, appErr)
		} else {
			response.ErrorWithCode(c, http.StatusUnauthorized, errors.ErrCodeUnauthorized, err.Error())
		}
		return
	}

	response.Write(c, http.StatusOK, resp)
}

// RefreshToken handles POST /auth/refresh
func (h *Handler) RefreshToken(c *gin.Context) {
	tenantID, ok := getTenantID(c)
	if !ok {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant not identified")
		return
	}

	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid JSON body")
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), tenantID, req.RefreshToken)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.AppError(c, appErr)
		} else {
			response.ErrorWithCode(c, http.StatusUnauthorized, errors.ErrCodeUnauthorized, err.Error())
		}
		return
	}

	response.Write(c, http.StatusOK, resp)
}

// RevokeToken handles POST /auth/revoke
func (h *Handler) RevokeToken(c *gin.Context) {
	// Get token from Authorization header
	header := c.GetHeader("Authorization")
	if header == "" {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "missing authorization header")
		return
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid authorization format")
		return
	}

	token := parts[1]
	if err := h.service.RevokeToken(c.Request.Context(), token); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.AppError(c, appErr)
		} else {
			response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, err.Error())
		}
		return
	}

	response.Write(c, http.StatusOK, gin.H{"message": "token revoked successfully"})
}

// ForgotPassword handles POST /auth/forgot-password
func (h *Handler) ForgotPassword(c *gin.Context) {
	tenantID, ok := getTenantID(c)
	if !ok {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "tenant not identified")
		return
	}

	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid JSON body")
		return
	}

	err := h.service.RequestPasswordReset(c.Request.Context(), tenantID, req.Email)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.AppError(c, appErr)
		} else {
			response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, err.Error())
		}
		return
	}

	// Always return success to prevent email enumeration
	response.Write(c, http.StatusOK, gin.H{
		"message": "If the email exists, a password reset link has been sent",
	})
}

// ResetPassword handles POST /auth/reset-password
func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid JSON body")
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), req.Token, req.Password)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			response.AppError(c, appErr)
		} else {
			response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, err.Error())
		}
		return
	}

	response.Write(c, http.StatusOK, gin.H{
		"message": "Password has been reset successfully",
	})
}
