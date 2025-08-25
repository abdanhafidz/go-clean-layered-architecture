package provider

import (
	"abdanhafidz.com/go-boilerplate/services"
)

type ServicesProvider interface {
	ProvideAccountService() services.AccountService
}
type servicesProvider struct {
	accountService services.AccountService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {
	accountService := services.NewAccountService(repoProvider.ProvideAccountRepository())

	return &servicesProvider{
		accountService: accountService,
	}
}

func (s *servicesProvider) ProvideAccountService() services.AccountService {
	return s.accountService
}
