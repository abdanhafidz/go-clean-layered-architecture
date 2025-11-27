package provider

import (
    "abdanhafidz.com/go-boilerplate/config"
    "abdanhafidz.com/go-boilerplate/services"
)

type ServicesProvider interface {
    ProvideEventService() services.EventService
    ProvideAcademyService() services.AcademyService
    ProvideProblemSetService() services.ProblemSetService
    ProvideJWTService() services.JWTService
    ProvideRegionService() services.RegionService
    ProvideOptionService() services.OptionService
    ProvideExamService() services.ExamService
    ProvideAccountService() services.AccountService
    ProvideForgotPasswordService() services.ForgotPasswordService
    ProvideEmailVerificationService() services.EmailVerificationService
    ProvideExternalAuthService() services.ExternalAuthService
    ProvideUploadService() services.UploadService 
    ProvideAcademyExamService() services.AcademyExamService
}

type servicesProvider struct {
    eventService             services.EventService
    academyService           services.AcademyService
    problemSetService        services.ProblemSetService
    jWTService               services.JWTService
    regionService            services.RegionService
    optionService            services.OptionService
    examService              services.ExamService
    accountService           services.AccountService
    forgotPasswordService    services.ForgotPasswordService
    emailVerificationService services.EmailVerificationService
    externalAuthService      services.ExternalAuthService
    uploadService            services.UploadService 
    academyExamService       services.AcademyExamService
}

func NewServicesProvider(
    repoProvider RepositoriesProvider, 
    configProvider ConfigProvider, 
    storageProvider StorageProvider, 
) ServicesProvider {

    eventService := services.NewEventService(repoProvider.ProvideEventsRepository(), repoProvider.ProvideEventAssignRepository())
    academyService := services.NewAcademyService(repoProvider.ProvideAcademyRepository())
    problemSetService := services.NewProblemSetService(repoProvider.ProvideProblemSetRepository(), repoProvider.ProvideQuestionsRepository(), repoProvider.ProvideProblemSetExamAssignRepository())
    jWTService := services.NewJWTService(configProvider.ProvideJWTConfig().GetSecretKey())
    regionService := services.NewRegionService(repoProvider.ProvideRegionRepository())
    optionService := services.NewOptionService(repoProvider.ProvideOptionRepository())
    examService := services.NewExamService(eventService, problemSetService, repoProvider.ProvideProblemSetExamAssignRepository(), repoProvider.ProvideExamRepository(), repoProvider.ProvideExamEventAttemptRepository(), repoProvider.ProvideExamEventAssignRepository(), repoProvider.ProvideExamEventAnswerRepository(), repoProvider.ProvideResultRepository())
    accountService := services.NewAccountService(jWTService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideAccountDetailRepository())
    forgotPasswordService := services.NewForgotPasswordService(jWTService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideForgotPasswordRepository())
    emailVerificationService := services.NewEmailVerificationService(accountService, repoProvider.ProvideEmailVerificationRepository())
    externalAuthService := services.NewExternalAuthService(jWTService, accountService, repoProvider.ProvideExternalAuthRepository())
    academyExamService := services.NewAcademyExamService(
        academyService,
        problemSetService,
        repoProvider.ProvideExamRepository(),
        repoProvider.ProvideExamAcademyAttemptRepository(),
        repoProvider.ProvideExamAcademyAssignRepository(),
        repoProvider.ProvideExamAcademyAnswerRepository(),
        repoProvider.ProvideAcademyResultRepository(),
    )

    uploadService := services.NewUploadService(
        storageProvider, 
        repoProvider.ProvideFileRepository(), 
        config.NewUploadConfig(),
    )

    return &servicesProvider{
        eventService:             eventService,
        academyService:           academyService,
        problemSetService:        problemSetService,
        jWTService:               jWTService,
        regionService:            regionService,
        optionService:            optionService,
        examService:              examService,
        accountService:           accountService,
        forgotPasswordService:    forgotPasswordService,
        emailVerificationService: emailVerificationService,
        externalAuthService:      externalAuthService,
        uploadService:            uploadService, 
        academyExamService:       academyExamService,
    }
}

func (s *servicesProvider) ProvideEventService() services.EventService {
    return s.eventService
}
func (s *servicesProvider) ProvideAcademyService() services.AcademyService {
    return s.academyService
}
func (s *servicesProvider) ProvideProblemSetService() services.ProblemSetService {
    return s.problemSetService
}
func (s *servicesProvider) ProvideJWTService() services.JWTService {
    return s.jWTService
}
func (s *servicesProvider) ProvideRegionService() services.RegionService {
    return s.regionService
}
func (s *servicesProvider) ProvideOptionService() services.OptionService {
    return s.optionService
}
func (s *servicesProvider) ProvideExamService() services.ExamService {
    return s.examService
}
func (s *servicesProvider) ProvideAccountService() services.AccountService {
    return s.accountService
}
func (s *servicesProvider) ProvideForgotPasswordService() services.ForgotPasswordService {
    return s.forgotPasswordService
}
func (s *servicesProvider) ProvideEmailVerificationService() services.EmailVerificationService {
    return s.emailVerificationService
}
func (s *servicesProvider) ProvideExternalAuthService() services.ExternalAuthService {
    return s.externalAuthService
}

func (s *servicesProvider) ProvideUploadService() services.UploadService {
    return s.uploadService
}

func (s *servicesProvider) ProvideAcademyExamService() services.AcademyExamService { return s.academyExamService }
