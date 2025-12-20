package router

import (
    "abdanhafidz.com/go-boilerplate/provider"
    "github.com/gin-contrib/gzip"
    "github.com/gin-gonic/gin"
)

func AcademyExamRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
    academyExamController := controller.ProvideAcademyExamController()
    authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()
    routerGroup := router.Group("/api/v1/academy")
    routerGroup.Use(gzip.Gzip(gzip.DefaultCompression)) 

    {
        routerGroup.GET("/:academy_slug/exam", authenticationMiddleware.VerifyAccount, academyExamController.List)
        routerGroup.GET("/:academy_slug/exam/:exam_slug/attempt", authenticationMiddleware.VerifyAccount, academyExamController.Attempt)
        routerGroup.POST("/:academy_slug/exam/:attempt_id/answer_question", authenticationMiddleware.VerifyAccount, academyExamController.Answer)
        routerGroup.POST("/:academy_slug/exam/:attempt_id/submit", authenticationMiddleware.VerifyAccount, academyExamController.Submit)
    }
}