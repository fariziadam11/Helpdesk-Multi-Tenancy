package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"werk-ticketing/internal/response"
)

// Handler handles HTTP requests for user operations
type Handler struct {
	repo Repository
}

// NewHandler creates a new user handler
func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// UpdateProfileRequest represents the request body for updating user profile
type UpdateProfileRequest struct {
	Name     *string `json:"name"`
	LastName *string `json:"lastname"`
	Password *string `json:"password"`
}

// UserResponse represents the user data returned to client
type UserResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	LastName      string `json:"lastname"`
	Email         string `json:"email"`
	InvGateUserID int    `json:"invgate_user_id"`
}

// UpdateProfile handles PUT /users/profile
func (h *Handler) UpdateProfile(c *gin.Context) {
	// Get user email from context (set by auth middleware)
	userEmail, exists := c.Get("userEmail")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Get existing user by email
	user, err := h.repo.GetByEmail(c.Request.Context(), userEmail.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get user")
		return
	}
	if user == nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	// Update fields if provided
	if req.Name != nil && *req.Name != "" {
		user.Name = *req.Name
	}
	if req.LastName != nil && *req.LastName != "" {
		user.LastName = *req.LastName
	}
	if req.Password != nil && *req.Password != "" {
		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to hash password")
			return
		}
		user.Password = string(hashedPassword)
	}

	// Update user in database
	if err := h.repo.Update(c.Request.Context(), user); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	// Return updated user data (without password)
	userResp := UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		LastName:      user.LastName,
		Email:         user.Email,
		InvGateUserID: user.InvGateUserID,
	}

	response.Success(c, http.StatusOK, userResp)
}
