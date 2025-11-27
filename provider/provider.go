package provider

import (
    "log"
    "strings"
    entity "abdanhafidz.com/go-boilerplate/models/entity"
    "github.com/gin-gonic/gin"
    "abdanhafidz.com/go-boilerplate/services"
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
    if supabaseCfg.GetURL() == "" || supabaseCfg.GetServiceKey() == "" || supabaseCfg.GetBucketName() == "" {
        log.Fatalf("Supabase configuration is invalid: URL, ServiceKey, or BucketName is empty")
    }
    if !strings.HasPrefix(supabaseCfg.GetURL(), "https://") || !strings.Contains(supabaseCfg.GetURL(), ".supabase.co") {
        log.Fatalf("Supabase URL is invalid")
    }
    if strings.Count(supabaseCfg.GetServiceKey(), ".") != 2 {
        log.Fatalf("Supabase service key is not a valid compact JWS")
    }
    storageDriver := services.NewSupabaseStorageService(supabaseCfg.GetURL(), supabaseCfg.GetServiceKey(), supabaseCfg.GetBucketName())
    servicesProvider := NewServicesProvider(repositoriesProvider, configProvider, storageDriver)
    controllerProvider := NewControllerProvider(servicesProvider)
    middlewareProvider := NewMiddlewareProvider(servicesProvider)

    // Database Migrations with error handling
    err := configProvider.ProvideDatabaseConfig().AutoMigrateAll(
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
        &entity.ExamAcademyAssign{},
        &entity.ExamAcademyAnswer{},
        &entity.ExamAcademyAttempt{},
        &entity.AcademyExamResult{},

        // Options & Regions
        &entity.OptionCategory{},
        &entity.OptionValues{},
        &entity.RegionProvince{},
        &entity.RegionCity{},

        // Files Storage
        &entity.File{},
    )
    
    if err != nil {
        log.Fatalf("Database migration failed: %v", err)
    }

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
