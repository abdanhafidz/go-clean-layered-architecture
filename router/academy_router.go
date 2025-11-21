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
        // ADMIN
        adminGroup := routerGroup.Group("/admin",  authenticationMiddleware.VerifyAccount)
        {
            // CRUD for Academies
            adminGroup.POST("/", academyController.CreateAcademy)
            adminGroup.GET("/id/:id/detail", academyController.GetAcademyDetail)
            adminGroup.PUT("/id/:id", academyController.UpdateAcademy)
            adminGroup.DELETE("/id/:id", academyController.DeleteAcademy)

            // Creation of Sub-resources
            adminGroup.POST("/materials", academyController.CreateMaterial)
            adminGroup.POST("/contents", academyController.CreateContent)
        }

        // USER
        routerGroup.GET("/", authenticationMiddleware.VerifyAccount, academyController.ListAcademies) 
        routerGroup.GET("/:academy_slug", authenticationMiddleware.VerifyAccount, academyController.GetAcademy)
        routerGroup.GET("/:academy_slug/:material_slug", authenticationMiddleware.VerifyAccount, academyController.GetMaterial)
        routerGroup.GET("/:academy_slug/:material_slug/:order", authenticationMiddleware.VerifyAccount, academyController.GetContent)
        routerGroup.POST("/:academy_slug/:material_slug/:order", authenticationMiddleware.VerifyAccount, academyController.UpdateContentProgress)
    }
}