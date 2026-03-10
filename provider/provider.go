package provider

import (
	"log"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
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
	log.Println("[BOOT] Initializing App Provider...")

	log.Println("[BOOT] Creating Gin Router")
	ginRouter := gin.Default()

	log.Println("[BOOT] Initializing Config Provider")
	configProvider := NewConfigProvider()

	log.Println("[BOOT] Initializing Repositories Provider")
	repositoriesProvider := NewRepositoriesProvider(configProvider)

	log.Println("[BOOT] Initializing Services Provider")
	servicesProvider := NewServicesProvider(repositoriesProvider, configProvider)

	log.Println("[BOOT] Initializing Controller Provider")
	controllerProvider := NewControllerProvider(servicesProvider)

	log.Println("[BOOT] Initializing Middleware Provider")
	middlewareProvider := NewMiddlewareProvider(servicesProvider)

	// ===============================
	// DATABASE MIGRATION
	// ===============================
	log.Println("[BOOT][DB] Starting database migration...")

	dbConfig := configProvider.ProvideDatabaseConfig()
	log.Println("[BOOT][DB] Database config acquired")

	err := dbConfig.AutoMigrateAll(
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

		// Files Storage
		&entity.File{},
	)

	if err != nil {
		log.Fatalf("[BOOT][DB] ❌ Database migration failed: %v", err)
	}

	log.Println("[BOOT][DB] ✅ Database migration completed")

	log.Println("[BOOT] App Provider initialized successfully")

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
