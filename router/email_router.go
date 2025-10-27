package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func EmailVerificationRoute(router *gin.Engine, controller provider.ControllerProvider) {
	emailVerificationController := controller.ProvideEmailVerificationController()
	routerGroup := router.Group("/api/v1/email")
	{
		routerGroup.POST("/verify", emailVerificationController.Validate)
		routerGroup.POST("/create-verification", emailVerificationController.Create)
	}
}
