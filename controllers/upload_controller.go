package controllers

import (
    "compress/gzip"
    "errors"
    "io"
    "net/http"
    "path/filepath"
    "strings"
    "fmt"

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
    fmt.Println("👉 Content-Type:", ctx.GetHeader("Content-Type"))

    if !strings.Contains(ctx.GetHeader("Content-Type"), "multipart/form-data") {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":  "error",
            "code":    "INVALID_FORM",
            "message": "Content-Type must be multipart/form-data",
        })
        return
    }

    if strings.EqualFold(ctx.GetHeader("Content-Encoding"), "gzip") {
        gz, err := gzip.NewReader(ctx.Request.Body)
        if err != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "status":  "error",
                "code":    "INVALID_FORM",
                "message": "Failed to decode gzip request body",
            })
            return
        }
        ctx.Request.Body = io.NopCloser(gz)
    }

    // Gunakan limit 32MB
    if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
        
        // 🔴 DEBUG: Print error ASLI ke terminal
        fmt.Println("❌ ERROR ParseMultipartForm:", err.Error())

        // Respon sementara dengan error asli agar terlihat di Postman
        ctx.JSON(http.StatusBadRequest, gin.H{
            "status":      "error",
            "code":        "INVALID_FORM",
            "message":     "Failed to parse form data",
            "debug_error": err.Error(), // <--- Kita butuh baca ini
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

	uploadContext := ctx.PostForm("context")
	if uploadContext == "" {
		ext := strings.ToLower(filepath.Ext(files[0].Filename))
		uploadContext = c.inferContextFromExt(ext)
	}

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

	uploadedFiles, err := c.uploadService.UploadFiles(ctx, files, uploadContext, accountID)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid Compact JWS") {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Storage misconfiguration: invalid Supabase service key",
			})
			return
		}
		if errors.Is(err, http_error.FILE_TOO_LARGE) ||
			errors.Is(err, http_error.INVALID_FILE_TYPE) ||
			errors.Is(err, http_error.BAD_REQUEST_ERROR) ||
			errors.Is(err, http_error.INVALID_DATA_PAYLOAD) {
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
				"data":    uploadedFiles,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

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
	fileIDStr := ctx.Param("id")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid file ID format",
		})
		return
	}

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

func (c *UploadController) inferContextFromExt(ext string) string {
	images := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true,
	}
	isSourceCode := map[string]bool{
		".cpp": true, ".c": true, ".py": true, ".java": true,
		".go": true, ".js": true, ".txt": true,
	}
	isDocument := map[string]bool{
		".pdf": true,
	}

	if images[ext] {
		return "image"
	}
	if isSourceCode[ext] {
		return "submission"
	}
	if isDocument[ext] {
		return "material"
	}
	return ""
}