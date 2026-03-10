package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func AdminRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()
	authenticationController := controller.ProvideAuthenticationController()

	// Authentication Admin Routes
	authAdminGroup := router.Group("/api/v1/admin/authentication", authenticationMiddleware.VerifyAccount)
	{
		authAdminGroup.PUT("/:account_id/assign", authenticationController.UpdateUserRole)
	}

}
