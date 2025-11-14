package provider

import (
	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"github.com/gin-gonic/gin"
)

type AppProvider interface {
	ProvideRouter() *gin.Engine
	ProvideConfig() ConfigProvider
	ProvideRepositories() RepositoriesProvider
	ProvideServices() ServicesProvider
	ProvideControllers() ControllerProvider
	ProvideMiddlewares() MiddlewareProvider
}
type appProvider struct {
	ginRouter            *gin.Engine
	configProvider       ConfigProvider
	repositoriesProvider RepositoriesProvider
	servicesProvider     ServicesProvider
	controllerProvider   ControllerProvider
	middlewareProvider   MiddlewareProvider
}

func NewAppProvider() AppProvider {
	ginRouter := gin.Default()
	configProvider := NewConfigProvider()
	repositoriesProvider := NewRepositoriesProvider(configProvider)
	servicesProvider := NewServicesProvider(repositoriesProvider, configProvider)
	controllerProvider := NewControllerProvider(servicesProvider)
	middlewareProvider := NewMiddlewareProvider(servicesProvider)
	configProvider.ProvideDatabaseConfig().AutoMigrateAll(
		// Accounts & Auth
		&entity.Account{},
		&entity.AccountDetail{},
		&entity.EmailVerification{},
		&entity.ExternalAuth{},
		&entity.FCM{},
		&entity.ForgotPassword{},

		// Options & Regions
		&entity.OptionCategory{},
		&entity.OptionValues{},
		&entity.RegionProvince{},
		&entity.RegionCity{},
	)

	return &appProvider{
		ginRouter:            ginRouter,
		configProvider:       configProvider,
		repositoriesProvider: repositoriesProvider,
		servicesProvider:     servicesProvider,
		controllerProvider:   controllerProvider,
		middlewareProvider:   middlewareProvider,
	}
}
func (a *appProvider) ProvideRouter() *gin.Engine {
	return a.ginRouter
}
func (a *appProvider) ProvideConfig() ConfigProvider {
	return a.configProvider
}

func (a *appProvider) ProvideRepositories() RepositoriesProvider {
	return a.repositoriesProvider
}

func (a *appProvider) ProvideServices() ServicesProvider {
	return a.servicesProvider
}

func (a *appProvider) ProvideControllers() ControllerProvider {
	return a.controllerProvider
}

func (a *appProvider) ProvideMiddlewares() MiddlewareProvider {
	return a.middlewareProvider
}
