package provider

import (
	"abdanhafidz.com/go-boilerplate/config"
	"abdanhafidz.com/go-boilerplate/services"
)

type ServicesProvider interface {
	ProvideRegionService() services.RegionService
	ProvideJWTService() services.JWTService
	ProvidePaymentService() services.PaymentService
	ProvideUploadService() services.UploadService
	ProvideOptionService() services.OptionService
	ProvideAccountService() services.AccountService
	ProvideForgotPasswordService() services.ForgotPasswordService
	ProvideEmailVerificationService() services.EmailVerificationService
	ProvideExternalAuthService() services.ExternalAuthService
}

type servicesProvider struct {
	regionService            services.RegionService
	jWTService               services.JWTService
	paymentService           services.PaymentService
	uploadService            services.UploadService
	optionService            services.OptionService
	accountService           services.AccountService
	forgotPasswordService    services.ForgotPasswordService
	emailVerificationService services.EmailVerificationService
	externalAuthService      services.ExternalAuthService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	regionService := services.NewRegionService(repoProvider.ProvideRegionRepository())
	jWTService := services.NewJWTService(configProvider.ProvideJWTConfig().GetSecretKey())
	paymentService := services.NewPaymentService(configProvider.ProvideXenditConfig().GetClient())
	storageService := services.NewSupabaseStorageService(configProvider.ProvideSupabaseConfig().GetURL(), configProvider.ProvideSupabaseConfig().GetServiceKey(), configProvider.ProvideSupabaseConfig().GetBucketName())
	uploadService := services.NewUploadService(
		storageService,
		repoProvider.ProvideFileRepository(),
		repoProvider.ProvideAccountRepository(),
		config.NewUploadConfig(),
	)
	optionService := services.NewOptionService(repoProvider.ProvideOptionRepository())
	accountService := services.NewAccountService(jWTService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideAccountDetailRepository())
	forgotPasswordService := services.NewForgotPasswordService(jWTService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideForgotPasswordRepository())
	emailVerificationService := services.NewEmailVerificationService(accountService, repoProvider.ProvideEmailVerificationRepository())
	externalAuthService := services.NewExternalAuthService(jWTService, accountService, repoProvider.ProvideExternalAuthRepository())
	return &servicesProvider{
		regionService:            regionService,
		jWTService:               jWTService,
		paymentService:           paymentService,
		uploadService:            uploadService,
		optionService:            optionService,
		accountService:           accountService,
		forgotPasswordService:    forgotPasswordService,
		emailVerificationService: emailVerificationService,
		externalAuthService:      externalAuthService,
	}
}

func (s *servicesProvider) ProvideRegionService() services.RegionService {
	return s.regionService
}

func (s *servicesProvider) ProvideJWTService() services.JWTService {
	return s.jWTService
}

func (s *servicesProvider) ProvidePaymentService() services.PaymentService {
	return s.paymentService
}

func (s *servicesProvider) ProvideUploadService() services.UploadService {
	return s.uploadService
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

func (s *servicesProvider) ProvideEmailVerificationService() services.EmailVerificationService {
	return s.emailVerificationService
}

func (s *servicesProvider) ProvideExternalAuthService() services.ExternalAuthService {
	return s.externalAuthService
}
