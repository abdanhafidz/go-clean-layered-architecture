package router

import (
	"abdanhafidz.com/go-clean-layered-architecture/provider"
	"github.com/gin-gonic/gin"
)

func ForgotPasswordRouter(router *gin.Engine, controller provider.ControllerProvider) {
	routerGroup := router.Group("/api/v1/forgot-password")
	forgotPasswordController := controller.ProvideForgotPasswordController()
	{
		routerGroup.POST("/", forgotPasswordController.Request)
	}
}
