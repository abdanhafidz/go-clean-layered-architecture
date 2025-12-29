package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func EventExamRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	eventExamController := controller.ProvideEventExamController()
	auth := middleware.ProvideAuthenticationMiddleware()
	routerGroup := router.Group("api/v1/events")
	{
		routerGroup.GET("/:event_slug/exam", auth.VerifyAccount, eventExamController.List)
		routerGroup.GET("/:event_slug/exam/:exam_slug/attempt", auth.VerifyAccount, eventExamController.Attempt)
		routerGroup.POST("/:event_slug/exam/:attempt_id/answer_question", auth.VerifyAccount, eventExamController.Answer)
		routerGroup.POST("/:event_slug/exam/:attempt_id/submit", auth.VerifyAccount, eventExamController.Submit)
	}
}
