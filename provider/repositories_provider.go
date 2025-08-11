package provider

import (
	"abdanhafidz.com/go-boilerplate/repositories"
)

type RepositoriesProvider interface {
	ProvideAccountRepository() repositories.AccountRepository
}
type repositoriesProvider struct {
	accountRepository repositories.AccountRepository
}

func NewRepositoriesProvider(cfg ConfigProvider) RepositoriesProvider {
	dbConfig := cfg.ProvideDatabaseConfig()
	accountRepository := repositories.NewAccountRepository(dbConfig.GetInstance())

	return &repositoriesProvider{
		accountRepository: accountRepository,
	}

}

func (r *repositoriesProvider) ProvideAccountRepository() repositories.AccountRepository {
	return r.accountRepository
}
