package provider

import "abdanhafidz.com/go-boilerplate/controllers"

type ControllerProvider interface {
	ProvideAcademyController() controllers.AcademyController
	ProvideAccountDetailController() controllers.AccountDetailController
	ProvideAuthenticationController() controllers.AuthenticationController
	ProvideEmailVerificationController() controllers.EmailVerificationController
	ProvideEventController() controllers.EventController
	ProvideExamController() controllers.ExamController
	ProvideForgotPasswordController() controllers.ForgotPasswordController
	ProvideOptionController() controllers.OptionController
	ProvideRegionController() controllers.RegionController
	
	// UPDATE: Menggunakan Pointer (*)
	ProvideUploadController() *controllers.UploadController 
}

type controllerProvider struct {
	academyController           controllers.AcademyController
	accountDetailController     controllers.AccountDetailController
	authenticationController    controllers.AuthenticationController
	emailVerificationController controllers.EmailVerificationController
	eventController             controllers.EventController
	examController              controllers.ExamController
	forgotPasswordController    controllers.ForgotPasswordController
	optionController            controllers.OptionController
	regionController            controllers.RegionController
	
	// UPDATE: Menggunakan Pointer (*)
	uploadController            *controllers.UploadController 
}

func NewControllerProvider(servicesProvider ServicesProvider) ControllerProvider {

	academyController := controllers.NewAcademyController(servicesProvider.ProvideAcademyService())
	accountDetailController := controllers.NewAccountDetailController(servicesProvider.ProvideAccountService())
	authenticationController := controllers.NewAuthenticationController(servicesProvider.ProvideAccountService(), servicesProvider.ProvideExternalAuthService())
	emailVerificationController := controllers.NewEmailVerificationController(servicesProvider.ProvideEmailVerificationService())
	eventController := controllers.NewEventController(servicesProvider.ProvideEventService())
	examController := controllers.NewExamController(servicesProvider.ProvideExamService())
	forgotPasswordController := controllers.NewForgotPasswordController(servicesProvider.ProvideForgotPasswordService())
	optionController := controllers.NewOptionController(servicesProvider.ProvideOptionService())
	regionController := controllers.NewRegionController(servicesProvider.ProvideRegionService())

	// UPDATE: Inisialisasi Upload Controller
	// servicesProvider.ProvideUploadService() sekarang sudah return Pointer (*), jadi aman.
	uploadController := controllers.NewUploadController(servicesProvider.ProvideUploadService())

	return &controllerProvider{
		academyController:           academyController,
		accountDetailController:     accountDetailController,
		authenticationController:    authenticationController,
		emailVerificationController: emailVerificationController,
		eventController:             eventController,
		examController:              examController,
		forgotPasswordController:    forgotPasswordController,
		optionController:            optionController,
		regionController:            regionController,
		uploadController:            uploadController, // Pointer assign ke Pointer
	}
}

// --- Getter Methods ---

func (c *controllerProvider) ProvideAcademyController() controllers.AcademyController {
	return c.academyController
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

// UPDATE: Return Pointer (*)
func (c *controllerProvider) ProvideUploadController() *controllers.UploadController {
	return c.uploadController
}