package provider

import "abdanhafidz.com/go-clean-layered-architecture/services"

type ServicesProvider interface {
	ProvideJWTService() services.JWTService
	ProvideRegionService() services.RegionService
	ProvideOptionService() services.OptionService
	ProvideAccountService() services.AccountService
	ProvideForgotPasswordService() services.ForgotPasswordService
	ProvideEmailVerificationService() services.EmailVerificationService
	ProvideExternalAuthService() services.ExternalAuthService
}

type servicesProvider struct {
	jWTService services.JWTService
	regionService services.RegionService
	optionService services.OptionService
	accountService services.AccountService
	forgotPasswordService services.ForgotPasswordService
	emailVerificationService services.EmailVerificationService
	externalAuthService services.ExternalAuthService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	jWTService := services.NewJWTService(configProvider.ProvideJWTConfig().GetSecretKey())
	regionService := services.NewRegionService(repoProvider.ProvideRegionRepository())
	optionService := services.NewOptionService(repoProvider.ProvideOptionRepository())
	accountService := services.NewAccountService(jWTService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideAccountDetailRepository())
	forgotPasswordService := services.NewForgotPasswordService(jWTService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideForgotPasswordRepository())
	emailVerificationService := services.NewEmailVerificationService(accountService, repoProvider.ProvideEmailVerificationRepository())
	externalAuthService := services.NewExternalAuthService(jWTService, accountService, repoProvider.ProvideExternalAuthRepository())

	return &servicesProvider{
		jWTService: jWTService,
		regionService: regionService,
		optionService: optionService,
		accountService: accountService,
		forgotPasswordService: forgotPasswordService,
		emailVerificationService: emailVerificationService,
		externalAuthService: externalAuthService,
	}
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

