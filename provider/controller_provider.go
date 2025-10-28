package provider

import (
	"abdanhafidz.com/go-boilerplate/controllers"
)

type ControllerProvider interface {
	ProvideAccountDetailController() controllers.AccountDetailController
	ProvideAuthenticationController() controllers.AuthenticationController
	ProvideEmailVerificationController() controllers.EmailVerificationController
	ProvideEventController() controllers.EventController
	ProvideForgotPasswordController() controllers.ForgotPasswordController
	ProvideOptionController() controllers.OptionController
	ProvideRegionController() controllers.RegionController
	ProvideAcademyController() controllers.AcademyController
}

type controllerProvider struct {
	accountDetailController     controllers.AccountDetailController
	authenticationController    controllers.AuthenticationController
	emailVerificationController controllers.EmailVerificationController
	eventController             controllers.EventController
	forgotPasswordController    controllers.ForgotPasswordController
	optionController            controllers.OptionController
	regionController            controllers.RegionController
	academyController           controllers.AcademyController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	accountDetailController := controllers.NewAccountDetailController(servicesProvider.ProvideAccountService())
	authenticationController := controllers.NewAuthenticationController(servicesProvider.ProvideAccountService())
	emailVerificationController := controllers.NewEmailVerificationController(servicesProvider.ProvideEmailVerificationService())
	eventController := controllers.NewEventController(servicesProvider.ProvideEventService())
	forgotPasswordController := controllers.NewForgotPasswordController(servicesProvider.ProvideForgotPasswordService())
	optionController := controllers.NewOptionController(servicesProvider.ProvideOptionService())
	regionController := controllers.NewRegionController(servicesProvider.ProvideRegionService())
	academyController := controllers.NewAcademyController(servicesProvider.ProvideAcademyService())
	return &controllerProvider{
		accountDetailController:     accountDetailController,
		authenticationController:    authenticationController,
		emailVerificationController: emailVerificationController,
		eventController:             eventController,
		forgotPasswordController:    forgotPasswordController,
		optionController:            optionController,
		regionController:            regionController,
		academyController:           academyController,
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

func (c *controllerProvider) ProvideEventController() controllers.EventController {
	return c.eventController
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

func (c *controllerProvider) ProvideAcademyController() controllers.AcademyController {
	return c.academyController
}
