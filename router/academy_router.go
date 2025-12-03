package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func AcademyRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	academyController := controller.ProvideAcademyController()
	authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()
	routerGroup := router.Group("/api/v1/academies")
	{

		routerGroup.GET("/", academyController.ListAcademies)
		routerGroup.POST("/", authenticationMiddleware.VerifyAccount, academyController.CreateAcademy)
		routerGroup.PUT("/:id", authenticationMiddleware.VerifyAccount, academyController.UpdateAcademy)
		routerGroup.DELETE("/:id", authenticationMiddleware.VerifyAccount, academyController.DeleteAcademy)
		routerGroup.GET("/:id", academyController.GetAcademy)
		routerGroup.GET("/:id/detail", academyController.GetAcademyDetail)
		routerGroup.POST("/materials", authenticationMiddleware.VerifyAccount, academyController.CreateMaterial)
		routerGroup.POST("/contents", authenticationMiddleware.VerifyAccount, academyController.CreateContent)
	}
}
