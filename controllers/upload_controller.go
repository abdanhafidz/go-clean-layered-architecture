package controllers

import (
	"errors" 
	"net/http"
	"path/filepath"
	"strings" 

	"abdanhafidz.com/go-boilerplate/models/dto"
	http_error "abdanhafidz.com/go-boilerplate/models/error" 
	"abdanhafidz.com/go-boilerplate/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadController struct {
	uploadService *services.UploadService
}

func NewUploadController(s *services.UploadService) *UploadController {
	return &UploadController{uploadService: s}
}

func (c *UploadController) Upload(ctx *gin.Context) {
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

	// 2. Get Files
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
	
	// 🎯 KLASIFIKASI OTOMATIS JIKA CONTEXT KOSONG
	if uploadContext == "" {
        ext := strings.ToLower(filepath.Ext(files[0].Filename))
        
        isSourceCode := map[string]bool{
            ".cpp": true, ".c": true, ".py": true, ".java": true, 
            ".go": true, ".js": true, ".txt": true,
        }[ext]

        isDocument := map[string]bool{
            ".pdf": true,
        }[ext]
        
        if isSourceCode {
            uploadContext = "submission" 
        } else if isDocument {
            uploadContext = "material" 
        } else {

            uploadContext = "general"
        }
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