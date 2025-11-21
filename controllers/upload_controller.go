package controllers

import (
	"net/http"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploadService *services.UploadService
}

func NewUploadController(s *services.UploadService) *UploadController {
	return &UploadController{uploadService: s}
}

func (c *UploadController) Upload(ctx *gin.Context) {
	if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse form"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}
	files := form.File["files"] 
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	urls, err := c.uploadService.UploadFiles(ctx, files)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Return Response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "files uploaded successfully",
		"urls":    urls,
	})
}