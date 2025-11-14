package provider

import "abdanhafidz.com/go-clean-layered-architecture/middleware"

type MiddlewareProvider interface {
	ProvideAuthenticationMiddleware() middleware.AuthenticationMiddleware
}

type middlewareProvider struct {
	authenticationMiddleware middleware.AuthenticationMiddleware
}

func NewMiddlewareProvider(servicesProvider ServicesProvider) MiddlewareProvider {
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(servicesProvider.ProvideJWTService())
	return &middlewareProvider{
		authenticationMiddleware: authenticationMiddleware,
	}
}

func (p *middlewareProvider) ProvideAuthenticationMiddleware() middleware.AuthenticationMiddleware {
	return p.authenticationMiddleware
}
