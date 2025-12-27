package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func EventExamProctoringRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	proctoringController := controller.ProvideEventExamProctoringController()
	auth := middleware.ProvideAuthenticationMiddleware()

	// Group under api/v1/events to match existing event structure
	routerGroup := router.Group("api/v1/events")
	{
		// :event_slug and :exam_slug are kept for context/consistency with other routes,
		// even if not strictly used for lookup (IDs are in body/query).

		routerGroup.POST("/:event_slug/exam/:exam_slug/proctoring", auth.VerifyAccount, proctoringController.CreateLog)
		routerGroup.GET("/:event_slug/exam/:exam_slug/proctoring", auth.VerifyAccount, proctoringController.ListLogs)
		routerGroup.GET("/:event_slug/exam/:exam_slug/proctoring/:log_id", auth.VerifyAccount, proctoringController.GetLogById)
		routerGroup.PUT("/:event_slug/exam/:exam_slug/proctoring/:log_id", auth.VerifyAccount, proctoringController.UpdateLog)
		routerGroup.DELETE("/:event_slug/exam/:exam_slug/proctoring/:log_id", auth.VerifyAccount, proctoringController.DeleteLog)
	}
}
