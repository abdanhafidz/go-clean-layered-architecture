package router

import (
	"abdanhafidz.com/go-clean-layered-architecture/provider"
	"github.com/gin-gonic/gin"
)

func AccountDetailRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	routerGroup := router.Group("/api/v1/account")
	accountDetailController := controller.ProvideAccountDetailController()
	authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()
	{
		routerGroup.GET("/me", authenticationMiddleware.VerifyAccount, accountDetailController.GetDetail)
		routerGroup.PUT("/me", authenticationMiddleware.VerifyAccount, accountDetailController.UpdateDetail)
	}
}
