package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-contrib/gzip" 
	"github.com/gin-gonic/gin"
)

func EmailVerificationRouter(router *gin.Engine, controller provider.ControllerProvider) {
	emailVerificationController := controller.ProvideEmailVerificationController()
	routerGroup := router.Group("/api/v1/email")
	routerGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		routerGroup.POST("/verify", emailVerificationController.Validate)
		routerGroup.POST("/create-verification", emailVerificationController.Create)
	}
}
