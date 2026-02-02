package upload

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"werk-ticketing/internal/errors"
	"werk-ticketing/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "no file provided")
		return
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".svg":  true,
		".webp": true,
		".gif":  true,
	}

	if !allowed[ext] {
		response.ErrorWithCode(c, http.StatusBadRequest, errors.ErrCodeInvalidInput, "invalid file type (allowed: jpg, jpeg, png, svg, webp, gif)")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	uploadDir := "./uploads"

	// Create dir if not exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to create upload directory")
			return
		}
	}

	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.ErrorWithCode(c, http.StatusInternalServerError, errors.ErrCodeInternal, "failed to save file")
		return
	}

	// Return public URL
	// Note: In production, this should be a full URL if serving from CDN or separate domain.
	// For local/simple setup, relative path works if static middleware is configured.
	url := fmt.Sprintf("/uploads/%s", filename)

	response.Success(c, http.StatusOK, gin.H{
		"url": url,
	})
}
