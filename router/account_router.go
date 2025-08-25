package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func AuthenticationRouter(router *gin.Engine, controller provider.ControllerProvider) {
	routerGroup := router.Group("/api/v1/account")
	accountController := controller.ProvideAccountController()
	{
		routerGroup.POST("/create", accountController.CreateAccount)
		routerGroup.POST("/verify", accountController.VerifyAccount)
		routerGroup.POST("/reset-pin", accountController.ResetPIN)
		routerGroup.POST("/block-account", accountController.BlockAccount)
	}
}
