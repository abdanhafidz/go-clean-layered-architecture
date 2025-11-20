package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-contrib/gzip" 
	"github.com/gin-gonic/gin"
)

func ExamEventRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	examController := controller.ProvideExamController()
	auth := middleware.ProvideAuthenticationMiddleware()
	routerGroup := router.Group("api/v1/events")
	routerGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		routerGroup.GET("/:event_slug/exam", auth.VerifyAccount, examController.List)
		routerGroup.GET("/:event_slug/exam/:exam_slug/attempt", auth.VerifyAccount, examController.Attempt)
		routerGroup.POST("/:event_slug/exam/:attempt_id/answer_question", auth.VerifyAccount, examController.Answer)
		routerGroup.POST("/:event_slug/exam/:attempt_id/submit", auth.VerifyAccount, examController.Submit)
	}
}
