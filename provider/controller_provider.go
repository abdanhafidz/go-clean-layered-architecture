package provider

import "abdanhafidz.com/go-boilerplate/controllers"

type ControllerProvider interface {
	ProvideAcademyController() controllers.AcademyController
	ProvideAcademyExamController() controllers.AcademyExamController
	ProvideAccountDetailController() controllers.AccountDetailController
	ProvideAuthenticationController() controllers.AuthenticationController
	ProvideEmailVerificationController() controllers.EmailVerificationController
	ProvideEventController() controllers.EventController
	ProvideEventExamController() controllers.EventExamController
	ProvideEventExamProctoringController() controllers.EventExamProctoringController
	ProvidePaymentCallbackController() controllers.PaymentCallbackController
	ProvideExamController() controllers.ExamController
	ProvideForgotPasswordController() controllers.ForgotPasswordController
	ProvideOptionController() controllers.OptionController
	ProvideRegionController() controllers.RegionController
	ProvideUploadController() controllers.UploadController
}

type controllerProvider struct {
	academyController             controllers.AcademyController
	academyExamController         controllers.AcademyExamController
	accountDetailController       controllers.AccountDetailController
	authenticationController      controllers.AuthenticationController
	emailVerificationController   controllers.EmailVerificationController
	eventController               controllers.EventController
	eventExamController           controllers.EventExamController
	eventExamProctoringController controllers.EventExamProctoringController
	paymentCallbackController     controllers.PaymentCallbackController
	examController                controllers.ExamController
	forgotPasswordController      controllers.ForgotPasswordController
	optionController              controllers.OptionController
	regionController              controllers.RegionController
	uploadController              controllers.UploadController
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	academyController := controllers.NewAcademyController(servicesProvider.ProvideAcademyService())
	academyExamController := controllers.NewAcademyExamController(servicesProvider.ProvideAcademyExamService())
	accountDetailController := controllers.NewAccountDetailController(servicesProvider.ProvideAccountService())
	authenticationController := controllers.NewAuthenticationController(servicesProvider.ProvideAccountService(), servicesProvider.ProvideExternalAuthService())
	emailVerificationController := controllers.NewEmailVerificationController(servicesProvider.ProvideEmailVerificationService())
	eventController := controllers.NewEventController(servicesProvider.ProvideEventService())
	eventExamController := controllers.NewEventExamController(servicesProvider.ProvideEventExamService())
	eventExamProctoringController := controllers.NewEventExamProctoringController(servicesProvider.ProvideEventExamProctoringService())
	paymentCallbackController := controllers.NewPaymentCallbackController(
		servicesProvider.ProvidePaymentService(),
		servicesProvider.ProvideEventService(),
		servicesProvider.ProvideAcademyService(),
	)
	examController := controllers.NewExamController(servicesProvider.ProvideExamService())
	forgotPasswordController := controllers.NewForgotPasswordController(servicesProvider.ProvideForgotPasswordService())
	optionController := controllers.NewOptionController(servicesProvider.ProvideOptionService())
	regionController := controllers.NewRegionController(servicesProvider.ProvideRegionService())
	uploadController := controllers.NewUploadController(servicesProvider.ProvideUploadService())
	return &controllerProvider{
		academyController:             academyController,
		academyExamController:         academyExamController,
		accountDetailController:       accountDetailController,
		authenticationController:      authenticationController,
		emailVerificationController:   emailVerificationController,
		eventController:               eventController,
		eventExamController:           eventExamController,
		eventExamProctoringController: eventExamProctoringController,
		paymentCallbackController:     paymentCallbackController,
		examController:                examController,
		forgotPasswordController:      forgotPasswordController,
		optionController:              optionController,
		regionController:              regionController,
		uploadController:              uploadController,
	}
}

// --- Getter Methods ---

func (c *controllerProvider) ProvideAcademyController() controllers.AcademyController {
	return c.academyController
}

func (c *controllerProvider) ProvideAcademyExamController() controllers.AcademyExamController {
	return c.academyExamController
}

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

func (c *controllerProvider) ProvideEventExamController() controllers.EventExamController {
	return c.eventExamController
}

func (c *controllerProvider) ProvideEventExamProctoringController() controllers.EventExamProctoringController {
	return c.eventExamProctoringController
}

func (c *controllerProvider) ProvidePaymentCallbackController() controllers.PaymentCallbackController {
	return c.paymentCallbackController
}

func (c *controllerProvider) ProvideExamController() controllers.ExamController {
	return c.examController
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
