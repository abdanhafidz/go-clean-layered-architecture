package router

import (
	"abdanhafidz.com/go-clean-layered-architecture/provider"
)

func RunRouter(appProvider provider.AppProvider) {
	router, controller, config, middleware := appProvider.ProvideRouter(), appProvider.ProvideControllers(), appProvider.ProvideConfig(), appProvider.ProvideMiddlewares()
	AuthenticationRouter(router, middleware, controller)
	ForgotPasswordRouter(router, controller)
	AccountDetailRouter(router, middleware, controller)
	EmailVerificationRouter(router, controller)
	EventRouter(router, middleware, controller)
	OptionsRouter(router, controller)
	AcademyRouter(router, middleware, controller)
	ExamEventRouter(router, middleware, controller)
	router.Run(config.ProvideEnvConfig().GetTCPAddress())
}
