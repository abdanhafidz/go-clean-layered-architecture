package provider

import (
	"abdanhafidz.com/go-boilerplate/services"
)

type ServicesProvider interface {
	ProvideAuthenticationService() services.AuthenticationService
	ProvideJWTService() services.JWTService
}
type servicesProvider struct {
	authenticationService services.AuthenticationService
	jwtService            services.JWTService
}

func NewServicesProvider(repoProvider RepositoriesProvider, configProvider ConfigProvider) ServicesProvider {

	env := configProvider.ProvideEnvConfig()
	jwtService := services.NewJWTService(env.GetSalt())
	authenticationService := services.NewAuthenticationService(repoProvider.ProvideAccountRepository(), jwtService)

	return &servicesProvider{
		authenticationService: authenticationService,
		jwtService:            jwtService,
	}
}

func (s *servicesProvider) ProvideAuthenticationService() services.AuthenticationService {
	return s.authenticationService
}

func (s *servicesProvider) ProvideJWTService() services.JWTService {
	return s.jwtService
}
