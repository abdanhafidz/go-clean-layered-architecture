package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func AdminRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()
	authorizationMiddleware := middleware.ProvideAuthorizationMiddleware()

	eventController := controller.ProvideEventController()
	academyController := controller.ProvideAcademyController()
	authenticationController := controller.ProvideAuthenticationController()

	// Event Admin Routes
	eventAdminGroup := router.Group("api/v1/admin/events", authenticationMiddleware.VerifyAccount, authorizationMiddleware.VerifyAdmin)
	{
		eventAdminGroup.POST("/", eventController.CreateEvent)
		eventAdminGroup.PUT("/:id", eventController.UpdateEvent)
		eventAdminGroup.DELETE("/:id", eventController.DeleteEvent)
	}

	// Academy Admin Routes
	academyAdminGroup := router.Group("/api/v1/admin/academy", authenticationMiddleware.VerifyAccount, authorizationMiddleware.VerifyAdmin)
	{
		academyAdminGroup.POST("/", academyController.CreateAcademy)
		academyAdminGroup.GET("/id/:id/detail", academyController.GetAcademyDetail)
		academyAdminGroup.PUT("/id/:id", academyController.UpdateAcademy)
		academyAdminGroup.DELETE("/id/:id", academyController.DeleteAcademy)

		academyAdminGroup.POST("/materials", academyController.CreateMaterial)
		academyAdminGroup.DELETE("/materials/:id", academyController.DeleteMaterial)

		academyAdminGroup.POST("/contents", academyController.CreateContent)
		academyAdminGroup.DELETE("/contents/:id", academyController.DeleteContent)

		academyAdminGroup.POST("/assign", academyController.AssignAccountToAcademy)
		academyAdminGroup.DELETE("/assign/:id", academyController.UnassignAccountFromAcademy)
		academyAdminGroup.GET("/assign/:academy_id", academyController.ListAssignmentsByAcademy)
	}

	// Authentication Admin Routes 
	authAdminGroup := router.Group("/api/v1/admin/authentication", authenticationMiddleware.VerifyAccount, authorizationMiddleware.VerifySuperAdmin)
	{
		authAdminGroup.PUT("/:account_id/assign", authenticationController.UpdateUserRole)
	}
}
