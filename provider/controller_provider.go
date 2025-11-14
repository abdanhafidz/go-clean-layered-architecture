package provider

import "abdanhafidz.com/go-clean-layered-architecture/controllers"

type ControllerProvider interface {
	ProvideAccountDetailController() controllers.AccountDetailController
	ProvideAuthenticationController() controllers.AuthenticationController
	ProvideEmailVerificationController() controllers.EmailVerificationController
	ProvideForgotPasswordController() controllers.ForgotPasswordController
	ProvideOptionController() controllers.OptionController
	ProvideRegionController() controllers.RegionController
}

type controllerProvider struct {
	accountDetailController     controllers.AccountDetailController
	authenticationController    controllers.AuthenticationController
	emailVerificationController controllers.EmailVerificationController
	forgotPasswordController    controllers.ForgotPasswordController
	optionController            controllers.OptionController
	regionController            controllers.RegionController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	accountDetailController := controllers.NewAccountDetailController(servicesProvider.ProvideAccountService())
	authenticationController := controllers.NewAuthenticationController(servicesProvider.ProvideAccountService(), servicesProvider.ProvideExternalAuthService())
	emailVerificationController := controllers.NewEmailVerificationController(servicesProvider.ProvideEmailVerificationService())
	forgotPasswordController := controllers.NewForgotPasswordController(servicesProvider.ProvideForgotPasswordService())
	optionController := controllers.NewOptionController(servicesProvider.ProvideOptionService())
	regionController := controllers.NewRegionController(servicesProvider.ProvideRegionService())
	return &controllerProvider{
		accountDetailController:     accountDetailController,
		authenticationController:    authenticationController,
		emailVerificationController: emailVerificationController,
		forgotPasswordController:    forgotPasswordController,
		optionController:            optionController,
		regionController:            regionController,
	}
}

// --- Getter Methods ---

func (c *controllerProvider) ProvideAccountDetailController() controllers.AccountDetailController {
	return c.accountDetailController
}

func (c *controllerProvider) ProvideAuthenticationController() controllers.AuthenticationController {
	return c.authenticationController
}

func (c *controllerProvider) ProvideEmailVerificationController() controllers.EmailVerificationController {
	return c.emailVerificationController
}

func (c *controllerProvider) ProvideForgotPasswordController() controllers.ForgotPasswordController {
	return c.forgotPasswordController
}

func (c *controllerProvider) ProvideOptionController() controllers.OptionController {
	return c.optionController
}

func (c *controllerProvider) ProvideRegionController() controllers.RegionController {
	return c.regionController
}
