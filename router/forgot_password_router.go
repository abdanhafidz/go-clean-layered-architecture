package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-contrib/gzip" 
	"github.com/gin-gonic/gin"
)

func ForgotPasswordRouter(router *gin.Engine, controller provider.ControllerProvider) {
	routerGroup := router.Group("/api/v1/forgot-password")
	forgotPasswordController := controller.ProvideForgotPasswordController()
	routerGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		routerGroup.POST("/", forgotPasswordController.Request)
	}
}
