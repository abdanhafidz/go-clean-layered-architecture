package provider

import (
	"abdanhafidz.com/go-boilerplate/repositories"
)

type RepositoriesProvider interface {
	ProvideAccountRepository() repositories.AccountRepository
	ProvideAccountDetailRepository() repositories.AccountDetailRepository
	ProvideEmailVerificationRepository() repositories.EmailVerificationRepository
	ProvideEventAssignRepository() repositories.EventAssignRepository
	ProvideEventsRepository() repositories.EventsRepository
	ProvideExternalAuthRepository() repositories.ExternalAuthRepository
	ProvideFCMRepository() repositories.FCMRepository
	ProvideForgotPasswordRepository() repositories.ForgotPasswordRepository
	ProvideOptionRepository() repositories.OptionRepository
	ProvideRegionRepository() repositories.RegionRepository
	ProvideAcademyRepository() repositories.AcademyRepository
}

type repositoriesProvider struct {
	accountRepository           repositories.AccountRepository
	accountDetailRepository     repositories.AccountDetailRepository
	emailVerificationRepository repositories.EmailVerificationRepository
	eventAssignRepository       repositories.EventAssignRepository
	eventsRepository            repositories.EventsRepository
	externalAuthRepository      repositories.ExternalAuthRepository
	fcmRepository               repositories.FCMRepository
	forgotPasswordRepository    repositories.ForgotPasswordRepository
	optionRepository            repositories.OptionRepository
	regionRepository            repositories.RegionRepository
	academyRepository           repositories.AcademyRepository
}

// NewRepositoriesProvider akan membuat semua repository dan inject konfigurasi database-nya.
func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	dbConfig := cfg.ProvideDatabaseConfig()
	db := dbConfig.GetInstance()

	accountRepo := repositories.NewAccountRepository(db)
	accountDetailRepo := repositories.NewAccountDetailRepository(db)
	emailVerificationRepo := repositories.NewEmailVerificationRepository(db)
	eventAssignRepo := repositories.NewEventAssignRepository(db)
	eventsRepo := repositories.NewEventsRepository(db)
	externalAuthRepo := repositories.NewExternalAuthRepository(db)
	fcmRepo := repositories.NewFCMRepository(db)
	forgotPasswordRepo := repositories.NewForgotPasswordRepository(db)
	optionRepo := repositories.NewOptionRepository(db)
	regionRepo := repositories.NewRegionRepository(db)
	academyRepo := repositories.NewAcademyRepository(db)

	return &repositoriesProvider{
		accountRepository:           accountRepo,
		accountDetailRepository:     accountDetailRepo,
		emailVerificationRepository: emailVerificationRepo,
		eventAssignRepository:       eventAssignRepo,
		eventsRepository:            eventsRepo,
		externalAuthRepository:      externalAuthRepo,
		fcmRepository:               fcmRepo,
		forgotPasswordRepository:    forgotPasswordRepo,
		optionRepository:            optionRepo,
		regionRepository:            regionRepo,
		academyRepository:           academyRepo,
	}
}

// --- Getter methods (implementasi interface) ---

func (r *repositoriesProvider) ProvideAccountRepository() repositories.AccountRepository {
	return r.accountRepository
}

func (r *repositoriesProvider) ProvideAccountDetailRepository() repositories.AccountDetailRepository {
	return r.accountDetailRepository
}

func (r *repositoriesProvider) ProvideEmailVerificationRepository() repositories.EmailVerificationRepository {
	return r.emailVerificationRepository
}

func (r *repositoriesProvider) ProvideEventAssignRepository() repositories.EventAssignRepository {
	return r.eventAssignRepository
}

func (r *repositoriesProvider) ProvideEventsRepository() repositories.EventsRepository {
	return r.eventsRepository
}

func (r *repositoriesProvider) ProvideExternalAuthRepository() repositories.ExternalAuthRepository {
	return r.externalAuthRepository
}

func (r *repositoriesProvider) ProvideFCMRepository() repositories.FCMRepository {
	return r.fcmRepository
}

func (r *repositoriesProvider) ProvideForgotPasswordRepository() repositories.ForgotPasswordRepository {
	return r.forgotPasswordRepository
}

func (r *repositoriesProvider) ProvideOptionRepository() repositories.OptionRepository {
	return r.optionRepository
}

func (r *repositoriesProvider) ProvideRegionRepository() repositories.RegionRepository {
	return r.regionRepository
}

func (r *repositoriesProvider) ProvideAcademyRepository() repositories.AcademyRepository {
	return r.academyRepository
}
