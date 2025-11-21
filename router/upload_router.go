package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)


func UploadRouter(r *gin.Engine, controller provider.ControllerProvider) {
	uploadController := controller.ProvideUploadController()
	routerGroup := r.Group("/api/v1/files")
	routerGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		routerGroup.POST("/", uploadController.Upload)
	}
}