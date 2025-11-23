package provider

import "abdanhafidz.com/go-boilerplate/middleware"

type MiddlewareProvider interface {
	ProvideAuthenticationMiddleware() middleware.AuthenticationMiddleware
	ProvideAuthorizationMiddleware() middleware.AuthorizationMiddleware
}

type middlewareProvider struct {
	authenticationMiddleware middleware.AuthenticationMiddleware
	authorizationMiddleware  middleware.AuthorizationMiddleware
}

func NewMiddlewareProvider(servicesProvider ServicesProvider) MiddlewareProvider {
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(servicesProvider.ProvideJWTService())
	authorizationMiddleware := middleware.NewAuthorizationMiddleware(servicesProvider.ProvideEventService())
	return &middlewareProvider{
		authenticationMiddleware: authenticationMiddleware,
		authorizationMiddleware:  authorizationMiddleware,
	}
}

func (p *middlewareProvider) ProvideAuthenticationMiddleware() middleware.AuthenticationMiddleware {
	return p.authenticationMiddleware
}

func (p *middlewareProvider) ProvideAuthorizationMiddleware() middleware.AuthorizationMiddleware {
	return p.authorizationMiddleware
}
