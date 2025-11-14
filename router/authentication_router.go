package router

import (
	"abdanhafidz.com/go-clean-layered-architecture/provider"
	"github.com/gin-gonic/gin"
)

func AuthenticationRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	routerGroup := router.Group("/api/v1/authentication")
	authenticationController := controller.ProvideAuthenticationController()
	authenticationmiddleware := middleware.ProvideAuthenticationMiddleware()
	{
		routerGroup.POST("/external-login", authenticationController.ExternalAuth)
		routerGroup.POST("/login", authenticationController.SignIn)
		routerGroup.POST("/register", authenticationController.SignUp)
		routerGroup.PUT("/change-password", authenticationmiddleware.VerifyAccount, authenticationController.ChangePassword)
	}
}
