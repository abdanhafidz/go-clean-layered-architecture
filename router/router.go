package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
)

func RunRouter(appProvider provider.AppProvider) {
	router, controller, config, middleware := appProvider.ProvideRouter(), appProvider.ProvideControllers(), appProvider.ProvideConfig(), appProvider.ProvideMiddlewares()
	AuthenticationRouter(router, middleware, controller)
	AccountDetailRouter(router, middleware, controller)
	EmailVerificationRoute(router, controller)
	EventRouter(router, middleware, controller)
	OptionsRouter(router, controller)
	router.Run(config.ProvideEnvConfig().GetTCPAddress())
}
