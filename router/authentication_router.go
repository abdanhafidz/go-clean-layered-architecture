package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func AuthenticationRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	routerGroup := router.Group("/api/v1/authentication")
	authenticationController := controller.ProvideAuthenticationController()
	authenticationmiddleware := middleware.ProvideAuthenticationMiddleware()
	authorizationMiddleware := middleware.ProvideAuthorizationMiddleware()
	{
		routerGroup.POST("/external-login", authenticationController.ExternalAuth)
		routerGroup.POST("/login", authenticationController.SignIn)
		routerGroup.POST("/register", authenticationController.SignUp)
		routerGroup.PUT("/change-password", authenticationmiddleware.VerifyAccount, authenticationController.ChangePassword)
	}

	userGroup := router.Group("/api/v1/users")
	{
		userGroup.PUT("/:account_id/assign", authorizationMiddleware.VerifySuperAdmin, authenticationController.UpdateUserRole)
	}
}
