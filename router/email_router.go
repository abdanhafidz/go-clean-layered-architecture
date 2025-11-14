package router

import (
	"abdanhafidz.com/go-clean-layered-architecture/provider"
	"github.com/gin-gonic/gin"
)

func EmailVerificationRouter(router *gin.Engine, controller provider.ControllerProvider) {
	emailVerificationController := controller.ProvideEmailVerificationController()
	routerGroup := router.Group("/api/v1/email")
	{
		routerGroup.POST("/verify", emailVerificationController.Validate)
		routerGroup.POST("/create-verification", emailVerificationController.Create)
	}
}
