package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func RunRouter(appProvider provider.AppProvider) {
	router, controller, config, middleware := appProvider.ProvideRouter(), appProvider.ProvideControllers(), appProvider.ProvideConfig(), appProvider.ProvideMiddlewares()

	router.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "OK",
			"message": "Service is up and running",
			"address": config.ProvideEnvConfig().GetTCPAddress(),
		})
	})

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "OK",
			"message": "Welcome to Quzuu API",
		})
	})

	AuthenticationRouter(router, middleware, controller)
	ForgotPasswordRouter(router, controller)
	AccountDetailRouter(router, middleware, controller)
	EmailVerificationRouter(router, controller)
	EventRouter(router, middleware, controller)
	OptionsRouter(router, controller)
	AcademyRouter(router, middleware, controller)
	ExamEventRouter(router, middleware, controller)
	AcademyExamRouter(router, middleware, controller)
	UploadRouter(router, middleware, controller)
	AdminRouter(router, middleware, controller)
	SwaggerRouter(router)
	router.Run(config.ProvideEnvConfig().GetTCPAddress())
}
