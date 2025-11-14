package provider

import "abdanhafidz.com/go-clean-layered-architecture/repositories"

type RepositoriesProvider interface {
	ProvideAccountDetailRepository() repositories.AccountDetailRepository
	ProvideAccountRepository() repositories.AccountRepository
	ProvideEmailVerificationRepository() repositories.EmailVerificationRepository
	ProvideExternalAuthRepository() repositories.ExternalAuthRepository
	ProvideFCMRepository() repositories.FCMRepository
	ProvideForgotPasswordRepository() repositories.ForgotPasswordRepository
	ProvideOptionRepository() repositories.OptionRepository
	ProvideRegionRepository() repositories.RegionRepository
}

type repositoriesProvider struct {
	accountDetailRepository     repositories.AccountDetailRepository
	accountRepository           repositories.AccountRepository
	emailVerificationRepository repositories.EmailVerificationRepository
	externalAuthRepository      repositories.ExternalAuthRepository
	fCMRepository               repositories.FCMRepository
	forgotPasswordRepository    repositories.ForgotPasswordRepository
	optionRepository            repositories.OptionRepository
	regionRepository            repositories.RegionRepository
}

func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	dbConfig := cfg.ProvideDatabaseConfig()
	db := dbConfig.GetInstance()

	accountDetailRepository := repositories.NewAccountDetailRepository(db)
	accountRepository := repositories.NewAccountRepository(db)
	emailVerificationRepository := repositories.NewEmailVerificationRepository(db)
	externalAuthRepository := repositories.NewExternalAuthRepository(db)
	fCMRepository := repositories.NewFCMRepository(db)
	forgotPasswordRepository := repositories.NewForgotPasswordRepository(db)
	optionRepository := repositories.NewOptionRepository(db)
	regionRepository := repositories.NewRegionRepository(db)

	return &repositoriesProvider{
		accountDetailRepository:     accountDetailRepository,
		accountRepository:           accountRepository,
		emailVerificationRepository: emailVerificationRepository,
		externalAuthRepository:      externalAuthRepository,
		fCMRepository:               fCMRepository,
		forgotPasswordRepository:    forgotPasswordRepository,
		optionRepository:            optionRepository,
		regionRepository:            regionRepository,
	}
}

func (r *repositoriesProvider) ProvideAccountDetailRepository() repositories.AccountDetailRepository {
	return r.accountDetailRepository
}

func (r *repositoriesProvider) ProvideAccountRepository() repositories.AccountRepository {
	return r.accountRepository
}

func (r *repositoriesProvider) ProvideEmailVerificationRepository() repositories.EmailVerificationRepository {
	return r.emailVerificationRepository
}

func (r *repositoriesProvider) ProvideExternalAuthRepository() repositories.ExternalAuthRepository {
	return r.externalAuthRepository
}

func (r *repositoriesProvider) ProvideFCMRepository() repositories.FCMRepository {
	return r.fCMRepository
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
