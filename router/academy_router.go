package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func AcademyRouter(router *gin.Engine, middleware provider.MiddlewareProvider, controller provider.ControllerProvider) {
	academyController := controller.ProvideAcademyController()
	authenticationMiddleware := middleware.ProvideAuthenticationMiddleware()

	routerGroup := router.Group("/api/v1/academy")

	routerGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		routerGroup.GET("/", authenticationMiddleware.VerifyAccount, academyController.ListAcademy)
		routerGroup.POST("/join", authenticationMiddleware.VerifyAccount, academyController.JoinAcademyByCode)
		routerGroup.GET("/:academy_slug", authenticationMiddleware.VerifyAccount, academyController.GetAcademy)
		routerGroup.GET("/:academy_slug/:material_slug", authenticationMiddleware.VerifyAccount, academyController.GetMaterial)
		routerGroup.GET("/:academy_slug/:material_slug/:order", authenticationMiddleware.VerifyAccount, academyController.GetContent)
		routerGroup.POST("/:academy_slug/:material_slug/:order", authenticationMiddleware.VerifyAccount, academyController.UpdateContentProgress)
	}
}
