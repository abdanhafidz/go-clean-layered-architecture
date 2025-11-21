package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"abdanhafidz.com/go-boilerplate/models/dto"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
)

type UploadController struct {
	uploadService *services.UploadService
}

func NewUploadController(s *services.UploadService) *UploadController {
	return &UploadController{uploadService: s}
}

func (c *UploadController) Upload(ctx *gin.Context) {
	// 1. Parse Multipart Form (Limit 32MB)
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		if strings.Contains(err.Error(), "http: request body too large") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "File size exceeds the allowed limit of 32MB",
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "INVALID_FORM",
			"message": "Failed to parse form data",
		})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"code":    "INVALID_DATA",
			"message": "Invalid form data",
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "No files uploaded",
		})
		return
	}

	// 2. Determine Upload Context
	uploadContext := ctx.PostForm("context")
	if uploadContext == "" {
		// Auto-infer context based on first file extension if not provided
		ext := strings.ToLower(filepath.Ext(files[0].Filename))
		uploadContext = c.inferContextFromExt(ext)
	}

	// 3. Get Account ID from Context (Middleware)
	accountIDStr := ctx.GetString("account_id")
	if accountIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized: Missing account ID",
		})
		return
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized: Invalid UUID format",
		})
		return
	}

	// 4. Call Service
	uploadedFiles, err := c.uploadService.UploadFiles(ctx, files, uploadContext, accountID)
	if err != nil {
		// Map Service Errors to HTTP Status
		if errors.Is(err, http_error.FILE_TOO_LARGE) ||
			errors.Is(err, http_error.INVALID_FILE_TYPE) ||
			errors.Is(err, http_error.BAD_REQUEST_ERROR) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		if errors.Is(err, http_error.PARTIAL_UPLOAD_FAILURE) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"status":  "error",
				"message": err.Error(),
				// Opsional: Anda bisa mengembalikan file yang berhasil di sini jika perlu
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// 5. Prepare Response DTO
	var fileResponses []dto.FileResponse
	for _, f := range uploadedFiles {
		fileResponses = append(fileResponses, dto.FileResponse{
			Id:           f.Id,
			OriginalName: f.OriginalName,
			URL:          f.Path,
			MimeType:     f.MimeType,
			Size:         f.Size,
			CreatedAt:    f.CreatedAt,
		})
	}

	ctx.JSON(http.StatusCreated, dto.FileUploadResponse{
		Status:  "success",
		Message: "Files uploaded successfully",
		Data:    fileResponses,
	})
}

func (c *UploadController) GetFileByID(ctx *gin.Context) {
	// 1. Validate Param ID
	fileIDStr := ctx.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid file ID format",
		})
		return
	}

	// 2. Validate Account ID
	accountIDStr := ctx.GetString("account_id")
	if accountIDStr == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized: Missing account ID",
		})
		return
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Unauthorized: Invalid UUID format",
		})
		return
	}

	// 3. Call Service
	fileData, err := c.uploadService.GetFileByID(ctx, fileID, accountID)
	if err != nil {
		if errors.Is(err, http_error.NOT_FOUND_ERROR) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "File not found or access denied",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// 4. Response
	response := dto.FileResponse{
		Id:           fileData.Id,
		OriginalName: fileData.OriginalName,
		URL:          fileData.Path,
		MimeType:     fileData.MimeType,
		Size:         fileData.Size,
		CreatedAt:    fileData.CreatedAt,
	}

	ctx.JSON(http.StatusOK, dto.FileResponseSingle{
		Status:  "success",
		Message: "File retrieved successfully",
		Data:    response,
	})
}

// Helper untuk logika inferensi context (clean code)
func (c *UploadController) inferContextFromExt(ext string) string {
	isSourceCode := map[string]bool{
		".cpp": true, ".c": true, ".py": true, ".java": true,
		".go": true, ".js": true, ".txt": true,
	}
	isDocument := map[string]bool{
		".pdf": true,
	}

	if isSourceCode[ext] {
		return "submission"
	}
	if isDocument[ext] {
		return "material"
	}
	return "general"
}