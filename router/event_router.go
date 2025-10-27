package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func EventRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	eventController := controller.ProvideEventController()
	authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()
	routerGroup := router.Group("api/v1/events")
	{
		routerGroup.GET("/", eventController.List)
		routerGroup.GET("/details/:slug", authenticationMiddleware.VerifyAccount, eventController.DetailBySlug)
		routerGroup.POST("/register-event", authenticationMiddleware.VerifyAccount, eventController.Join)
	}
}
