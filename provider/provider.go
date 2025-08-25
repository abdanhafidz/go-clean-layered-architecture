package provider

import (
	"abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/gin-gonic/gin"
)

type AppProvider interface {
	ProvideRouter() *gin.Engine
	ProvideConfig() ConfigProvider
	ProvideRepositories() RepositoriesProvider
	ProvideServices() ServicesProvider
	ProvideControllers() ControllerProvider
}
type appProvider struct {
	ginRouter            *gin.Engine
	configProvider       ConfigProvider
	repositoriesProvider RepositoriesProvider
	servicesProvider     ServicesProvider
	controllerProvider   ControllerProvider
}

func NewAppProvider() AppProvider {
	ginRouter := gin.Default()
	configProvider := NewConfigProvider()
	repositoriesProvider := NewRepositoriesProvider(configProvider)
	servicesProvider := NewServicesProvider(repositoriesProvider, configProvider)
	controllerProvider := NewControllerProvider(servicesProvider)

	configProvider.ProvideDatabaseConfig().AutoMigrateAll(
		&entity.Account{},
	)

	return &appProvider{
		ginRouter:            ginRouter,
		configProvider:       configProvider,
		repositoriesProvider: repositoriesProvider,
		servicesProvider:     servicesProvider,
		controllerProvider:   controllerProvider,
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
