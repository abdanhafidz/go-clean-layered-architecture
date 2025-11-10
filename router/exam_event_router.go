package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func ExamEventRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	examController := controller.ProvideExamController()
	auth := middleware.ProvideAuthenticationMiddleware()

	routerGroup := router.Group("api/v1/exam")
	{
		routerGroup.GET("/:event_slug/:exam_slug/attempt",
			auth.VerifyAccount,
			examController.Attempt,
		)
		routerGroup.POST("/answer_question/:attempt_id",
			auth.VerifyAccount,
			examController.Answer,
		)
		routerGroup.POST("/submit/:attempt_id",
			auth.VerifyAccount,
			examController.Submit,
		)
	}
}
