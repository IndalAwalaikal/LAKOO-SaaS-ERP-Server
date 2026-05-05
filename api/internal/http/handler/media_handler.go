package handler

import (
	"path/filepath"
	"net/http"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	
	"lakoo/backend/pkg/storage"
	"lakoo/backend/pkg/response"
)

type MediaHandler struct {
	minioService storage.MinioService
}

func NewMediaHandler(ms storage.MinioService) *MediaHandler {
	return &MediaHandler{minioService: ms}
}

func (h *MediaHandler) Upload(c *gin.Context) {
	tenantID, exists := c.Get("tenant_id")
	if !exists {
		response.Error(c, 403, "FORBIDDEN", "Tenant Context Missing")
		return
	}

	// 1. Enforce 5MB limit
	const MaxFileSize = 5 * 1024 * 1024 
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxFileSize)

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, 400, "BAD_REQUEST", "File is required or exceeds size limit (5MB)")
		return
	}

	if file.Size > MaxFileSize {
		response.Error(c, 400, "BAD_REQUEST", "File size exceeds 5MB limit")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Error(c, 500, "INTERNAL_ERROR", "Could not process file")
		return
	}
	defer src.Close()

	// 2. Robust MIME Validation using http.DetectContentType
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		response.Error(c, 500, "INTERNAL_ERROR", "Failed to read file for validation")
		return
	}
	
	// Reset pointer after reading buffer
	src.Seek(0, 0)
	
	detectedType := http.DetectContentType(buffer)
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	if !allowedTypes[detectedType] {
		response.Error(c, 400, "BAD_REQUEST", "Invalid file type. Only JPEG, PNG, and WEBP are allowed.")
		return
	}

	// Generate secure filename
	ext := filepath.Ext(file.Filename)
	uniqueName := tenantID.(string) + "/media/" + uuid.New().String() + ext

	objectURL, err := h.minioService.UploadFile(c.Request.Context(), uniqueName, src, file.Size, detectedType)
	if err != nil {
		response.Error(c, 500, "INTERNAL_ERROR", "Failed to upload file: " + err.Error())
		return
	}

	response.Success(c, 201, gin.H{
		"url": objectURL,
		"filename": uniqueName,
		"original_name": file.Filename,
	})
}
