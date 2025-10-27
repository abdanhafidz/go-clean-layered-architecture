package provider

import (
	"abdanhafidz.com/go-boilerplate/services"
)

type ServicesProvider interface {
	ProvideAccountService() services.AccountService
	ProvideEmailVerificationService() services.EmailVerificationService
	ProvideEventService() services.EventService
	ProvideForgotPasswordService() services.ForgotPasswordService
	ProvideJWTService() services.JWTService
	ProvideOptionService() services.OptionService
	ProvideRegionService() services.RegionService
}

type servicesProvider struct {
	accountService           services.AccountService
	emailVerificationService services.EmailVerificationService
	eventService             services.EventService
	forgotPasswordService    services.ForgotPasswordService
	jwtService               services.JWTService
	optionService            services.OptionService
	regionService            services.RegionService
}

// Konstruktor utama yang menginisialisasi semua service
func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	jwtService := services.NewJWTService(configProvider.ProvideJWTConfig().GetSecretKey())
	accountService := services.NewAccountService(jwtService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideAccountDetailRepository())
	emailVerificationService := services.NewEmailVerificationService(repoProvider.ProvideAccountRepository(), repoProvider.ProvideEmailVerificationRepository())
	eventService := services.NewEventService(repoProvider.ProvideEventsRepository(), repoProvider.ProvideEventAssignRepository())
	forgotPasswordService := services.NewForgotPasswordService(jwtService, repoProvider.ProvideAccountRepository(), repoProvider.ProvideForgotPasswordRepository())
	optionService := services.NewOptionService(repoProvider.ProvideOptionRepository())
	regionService := services.NewRegionService(repoProvider.ProvideRegionRepository())
	return &servicesProvider{
		accountService:           accountService,
		emailVerificationService: emailVerificationService,
		eventService:             eventService,
		forgotPasswordService:    forgotPasswordService,
		jwtService:               jwtService,
		optionService:            optionService,
		regionService:            regionService,
	}
}

// Getter methods (implementasi interface)
func (s *servicesProvider) ProvideAccountService() services.AccountService {
	return s.accountService
}

func (s *servicesProvider) ProvideEmailVerificationService() services.EmailVerificationService {
	return s.emailVerificationService
}

func (s *servicesProvider) ProvideEventService() services.EventService {
	return s.eventService
}

func (s *servicesProvider) ProvideForgotPasswordService() services.ForgotPasswordService {
	return s.forgotPasswordService
}

func (s *servicesProvider) ProvideJWTService() services.JWTService {
	return s.jwtService
}

func (s *servicesProvider) ProvideOptionService() services.OptionService {
	return s.optionService
}

func (s *servicesProvider) ProvideRegionService() services.RegionService {
	return s.regionService
}
