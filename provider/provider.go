package provider

import (
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
	ginRouter := gin.Default()
	configProvider := NewConfigProvider()
	repositoriesProvider := NewRepositoriesProvider(configProvider)
	supabaseCfg := configProvider.ProvideSupabaseConfig()
	storageDriver := NewSupabaseStorage(supabaseCfg.URL, supabaseCfg.ServiceKey, supabaseCfg.BucketName)
	servicesProvider := NewServicesProvider(repositoriesProvider, configProvider, storageDriver)

	controllerProvider := NewControllerProvider(servicesProvider)
	middlewareProvider := NewMiddlewareProvider(servicesProvider)

	// Database Migrations
	configProvider.ProvideDatabaseConfig().AutoMigrateAll(
		// Accounts & Auth
		&entity.Account{},
		&entity.AccountDetail{},
		&entity.EmailVerification{},
		&entity.ExternalAuth{},
		&entity.FCM{},
		&entity.ForgotPassword{},

		// Events
		&entity.Events{},
		&entity.EventAssign{},
		&entity.Announcement{},

		// Problemset & Exam
		&entity.ProblemSet{},
		&entity.Questions{},
		&entity.Exam{},
		&entity.ProblemSetExamAssign{},
		&entity.ExamEventAssign{},

		// Exam Attempt & Result
		&entity.ExamEventAnswer{},
		&entity.ExamEventAttempt{},
		&entity.Result{},

		// Academy LMS
		&entity.Academy{},
		&entity.AcademyMaterial{},
		&entity.AcademyContent{},
		&entity.AcademyMaterialProgress{},
		&entity.AcademyContentProgress{},
		&entity.AcademyProgress{},

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