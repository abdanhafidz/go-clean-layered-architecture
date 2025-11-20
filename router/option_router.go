package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-contrib/gzip" 
	"github.com/gin-gonic/gin"
)

func OptionsRouter(router *gin.Engine, controller provider.ControllerProvider) {
	optionsController := controller.ProvideOptionController()
	regionController := controller.ProvideRegionController()
	routerGroup := router.Group("/api/v1/options")
	routerGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		routerGroup.POST("/create", optionsController.CreateBulk)
		routerGroup.GET("/list/:slug", optionsController.GetBySlug)
		routerGroup.GET("/region/provinces", regionController.ListProvinces)
		routerGroup.GET("/region/cities", regionController.ListCitiesByProvince)
		routerGroup.POST("/region/seed-provinces", regionController.SeedProvinces)
		routerGroup.POST("/region/seed-cities", regionController.SeedCities)
	}
}
