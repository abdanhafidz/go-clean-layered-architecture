package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func EventRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	eventController := controller.ProvideEventController()
	authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()
	authorizationMiddleware := middleware.ProvideAuthorizationMiddleware()
	routerGroup := router.Group("api/v1/events")
	{
		routerGroup.GET("/", eventController.List)
		routerGroup.GET("/:slug", authenticationMiddleware.VerifyAccount, eventController.DetailBySlug)
		routerGroup.POST("/register-event", authenticationMiddleware.VerifyAccount, eventController.Join)
		routerGroup.GET("/:slug/quizzes", authenticationMiddleware.VerifyAccount, authorizationMiddleware.AuthorizeUserToEvent, eventController.QuizListByEvent)
	}
}
