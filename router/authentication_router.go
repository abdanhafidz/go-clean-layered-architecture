package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func AuthenticationRouter(router *gin.Engine, controller provider.ControllerProvider) {
	routerGroup := router.Group("/api/v1/auth")
	{
		routerGroup.POST("/login", controller.ProvideAuthenticationController().Login)
		routerGroup.POST("/register", controller.ProvideAuthenticationController().Register)
	}
}
