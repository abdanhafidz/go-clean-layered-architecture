package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
)

func RunRouter(appProvider provider.AppProvider) {
	router, controller, config := appProvider.ProvideRouter(), appProvider.ProvideControllers(), appProvider.ProvideConfig()
	AuthenticationRouter(router, controller)
	router.Run(config.ProvideEnvConfig().GetTCPAddress())
}
