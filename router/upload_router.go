package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func UploadRouter(r *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	uploadController := controller.ProvideUploadController()
	authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()

	routerGroup := r.Group("/api/v1/files")
	routerGroup.Use(gzip.Gzip(gzip.DefaultCompression), authenticationMiddleware.VerifyAccount)

	{
		routerGroup.POST("/", uploadController.Upload)
		routerGroup.GET("/:id", uploadController.GetFileByID)
	}
}