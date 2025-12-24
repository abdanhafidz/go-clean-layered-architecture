package provider

import (
	"abdanhafidz.com/go-boilerplate/config"
	"abdanhafidz.com/go-boilerplate/services"
)

type ServicesProvider interface {
	ProvideRegionService() services.RegionService
	ProvideJWTService() services.JWTService
	ProvideAcademyService() services.AcademyService
	ProvidePaymentService() services.PaymentService
	ProvideUploadService() services.UploadService
	ProvideProblemSetService() services.ProblemSetService
	ProvideOptionService() services.OptionService
	ProvideAccountService() services.AccountService
	ProvideForgotPasswordService() services.ForgotPasswordService
	ProvideEventService() services.EventService
	ProvideAcademyExamService() services.AcademyExamService
	ProvideEmailVerificationService() services.EmailVerificationService
	ProvideExternalAuthService() services.ExternalAuthService
	ProvideExamService() services.ExamService
}

type servicesProvider struct {
	regionService            services.RegionService
	jWTService               services.JWTService
	academyService           services.AcademyService
	paymentService           services.PaymentService
	uploadService            services.UploadService
	problemSetService        services.ProblemSetService
	optionService            services.OptionService
	accountService           services.AccountService
	forgotPasswordService    services.ForgotPasswordService
	eventService             services.EventService
	academyExamService       services.AcademyExamService
	emailVerificationService services.EmailVerificationService
	externalAuthService      services.ExternalAuthService
	examService              services.ExamService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	regionService := services.NewRegionService(repoProvider.ProvideRegionRepository())
	jWTService := services.NewJWTService(configProvider.ProvideJWTConfig().GetSecretKey())
	paymentService := services.NewPaymentService(configProvider.ProvideXenditConfig().GetClient(), repoProvider.ProvideEventPaymentRepository(), repoProvider.ProvideAcademyPaymentRepository())
	academyService := services.NewAcademyService(paymentService, repoProvider.ProvideAcademyRepository())
	storageService := services.NewSupabaseStorageService(configProvider.ProvideSupabaseConfig().GetURL(), configProvider.ProvideSupabaseConfig().GetServiceKey(), configProvider.ProvideSupabaseConfig().GetBucketName())
	uploadService := services.NewUploadService(
		storageService,
		repoProvider.ProvideFileRepository(),
		repoProvider.ProvideAccountRepository(),
		config.NewUploadConfig(),
	)
	problemSetService := services.NewProblemSetService(repoProvider.ProvideProblemSetRepository(), repoProvider.ProvideQuestionsRepository(), repoProvider.ProvideProblemSetExamAssignRepository())
	optionService := services.NewOptionService(repoProvider.ProvideOptionRepository())
	accountService := services.NewAccountService(jWTService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideAccountDetailRepository())
	forgotPasswordService := services.NewForgotPasswordService(jWTService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideForgotPasswordRepository())
	eventService := services.NewEventService(paymentService, repoProvider.ProvideEventsRepository(), repoProvider.ProvideEventAssignRepository())
	academyExamService := services.NewAcademyExamService(academyService, problemSetService, repoProvider.ProvideExamRepository(), repoProvider.ProvideExamAcademyAttemptRepository(), repoProvider.ProvideExamAcademyAssignRepository(), repoProvider.ProvideExamAcademyAnswerRepository(), repoProvider.ProvideAcademyResultRepository())
	emailVerificationService := services.NewEmailVerificationService(accountService, repoProvider.ProvideEmailVerificationRepository())
	externalAuthService := services.NewExternalAuthService(jWTService, accountService, repoProvider.ProvideExternalAuthRepository())
	examService := services.NewExamService(eventService, problemSetService, repoProvider.ProvideProblemSetExamAssignRepository(), repoProvider.ProvideExamRepository(), repoProvider.ProvideExamEventAttemptRepository(), repoProvider.ProvideExamEventAssignRepository(), repoProvider.ProvideExamEventAnswerRepository(), repoProvider.ProvideResultRepository())

	return &servicesProvider{
		regionService:            regionService,
		jWTService:               jWTService,
		academyService:           academyService,
		paymentService:           paymentService,
		uploadService:            uploadService,
		problemSetService:        problemSetService,
		optionService:            optionService,
		accountService:           accountService,
		forgotPasswordService:    forgotPasswordService,
		eventService:             eventService,
		academyExamService:       academyExamService,
		emailVerificationService: emailVerificationService,
		externalAuthService:      externalAuthService,
		examService:              examService,
	}
}

func (s *servicesProvider) ProvideRegionService() services.RegionService {
	return s.regionService
}

func (s *servicesProvider) ProvideJWTService() services.JWTService {
	return s.jWTService
}

func (s *servicesProvider) ProvideAcademyService() services.AcademyService {
	return s.academyService
}

func (s *servicesProvider) ProvidePaymentService() services.PaymentService {
	return s.paymentService
}

func (s *servicesProvider) ProvideUploadService() services.UploadService {
	return s.uploadService
}

func (s *servicesProvider) ProvideProblemSetService() services.ProblemSetService {
	return s.problemSetService
}

func (s *servicesProvider) ProvideOptionService() services.OptionService {
	return s.optionService
}

func (s *servicesProvider) ProvideAccountService() services.AccountService {
	return s.accountService
}

func (s *servicesProvider) ProvideForgotPasswordService() services.ForgotPasswordService {
	return s.forgotPasswordService
}

func (s *servicesProvider) ProvideEventService() services.EventService {
	return s.eventService
}

func (s *servicesProvider) ProvideAcademyExamService() services.AcademyExamService {
	return s.academyExamService
}

func (s *servicesProvider) ProvideEmailVerificationService() services.EmailVerificationService {
	return s.emailVerificationService
}

func (s *servicesProvider) ProvideExternalAuthService() services.ExternalAuthService {
	return s.externalAuthService
}

func (s *servicesProvider) ProvideExamService() services.ExamService {
	return s.examService
}
