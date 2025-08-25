package provider

import "abdanhafidz.com/go-boilerplate/controller"

type ControllerProvider interface {
	ProvideAccountController() controller.AccountController
}

type controllerProvider struct {
	accountController controller.AccountController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	accountController := controller.NewAccountController(servicesProvider.ProvideAccountService())

	return &controllerProvider{
		accountController: accountController,
	}
}

func (c *controllerProvider) ProvideAccountController() controller.AccountController {
	return c.accountController
}
