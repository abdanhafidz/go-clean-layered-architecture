package provider

import "abdanhafidz.com/go-boilerplate/controllers"

type ControllerProvider interface {
	ProvideAccountDetailController() controllers.AccountDetailController
	ProvideAuthenticationController() controllers.AuthenticationController
	ProvideEmailVerificationController() controllers.EmailVerificationController
	ProvidePaymentCallbackController() controllers.PaymentCallbackController
	ProvideForgotPasswordController() controllers.ForgotPasswordController
	ProvideOptionController() controllers.OptionController
	ProvideRegionController() controllers.RegionController
	ProvideUploadController() controllers.UploadController
}

type controllerProvider struct {
	accountDetailController     controllers.AccountDetailController
	authenticationController    controllers.AuthenticationController
	emailVerificationController controllers.EmailVerificationController
	paymentCallbackController   controllers.PaymentCallbackController
	forgotPasswordController    controllers.ForgotPasswordController
	optionController            controllers.OptionController
	regionController            controllers.RegionController
	uploadController            controllers.UploadController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	accountDetailController := controllers.NewAccountDetailController(servicesProvider.ProvideAccountService())
	authenticationController := controllers.NewAuthenticationController(servicesProvider.ProvideAccountService(), servicesProvider.ProvideExternalAuthService())
	emailVerificationController := controllers.NewEmailVerificationController(servicesProvider.ProvideEmailVerificationService())
	paymentCallbackController := controllers.NewPaymentCallbackController(
		servicesProvider.ProvidePaymentService(),
	)
	forgotPasswordController := controllers.NewForgotPasswordController(servicesProvider.ProvideForgotPasswordService())
	optionController := controllers.NewOptionController(servicesProvider.ProvideOptionService())
	regionController := controllers.NewRegionController(servicesProvider.ProvideRegionService())
	uploadController := controllers.NewUploadController(servicesProvider.ProvideUploadService())
	return &controllerProvider{
		accountDetailController:     accountDetailController,
		authenticationController:    authenticationController,
		emailVerificationController: emailVerificationController,
		paymentCallbackController:   paymentCallbackController,
		forgotPasswordController:    forgotPasswordController,
		optionController:            optionController,
		regionController:            regionController,
		uploadController:            uploadController,
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

func (c *controllerProvider) ProvidePaymentCallbackController() controllers.PaymentCallbackController {
	return c.paymentCallbackController
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

func (c *controllerProvider) ProvideUploadController() controllers.UploadController {
	return c.uploadController
}
