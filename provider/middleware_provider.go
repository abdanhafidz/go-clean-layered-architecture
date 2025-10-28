package provider

import "abdanhafidz.com/go-boilerplate/middleware"

type MiddlewareProvider interface {
	ProvideAuthenticationMiddleware() middleware.AuthenticationMiddleware
	ProvideAuthorizationMiddleware() middleware.AuthorizationMiddleware
}

type middlewareProvider struct {
	authorizationMiddleware  middleware.AuthorizationMiddleware
	authenticationMiddleware middleware.AuthenticationMiddleware
}

func NewMiddlewareProvider(servicesProvider ServicesProvider) MiddlewareProvider {
	authorizationMiddleware := middleware.NewAuthorizationMiddleware(servicesProvider.ProvideEventService())
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(servicesProvider.ProvideJWTService())
	return &middlewareProvider{
		authorizationMiddleware:  authorizationMiddleware,
		authenticationMiddleware: authenticationMiddleware,
	}
}

func (p *middlewareProvider) ProvideAuthenticationMiddleware() middleware.AuthenticationMiddleware {
	return p.authenticationMiddleware
}

func (p *middlewareProvider) ProvideAuthorizationMiddleware() middleware.AuthorizationMiddleware {
	return p.authorizationMiddleware
}
