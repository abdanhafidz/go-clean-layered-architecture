package provider

import "abdanhafidz.com/go-boilerplate/controller"

type ControllerProvider interface {
	ProvideAuthenticationController() controller.AuthenticationController
}

type controllerProvider struct {
	authenticationController controller.AuthenticationController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	authenticationController := controller.NewAuthenticationController(servicesProvider.ProvideAuthenticationService())

	return &controllerProvider{
		authenticationController: authenticationController,
	}
}

func (c *controllerProvider) ProvideAuthenticationController() controller.AuthenticationController {
	return c.authenticationController
}
