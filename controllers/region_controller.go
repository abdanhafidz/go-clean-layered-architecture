package controllers

import (
	"github.com/gin-gonic/gin"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"abdanhafidz.com/go-clean-layered-architecture/services"
)

type RegionController interface {
	SeedProvinces(ctx *gin.Context)
	SeedCities(ctx *gin.Context)
	ListProvinces(ctx *gin.Context)
	ListCitiesByProvince(ctx *gin.Context)
}

type regionController struct {
	regionService services.RegionService
}

func NewRegionController(regionService services.RegionService) RegionController {
	return &regionController{regionService: regionService}
}

func (c *regionController) SeedProvinces(ctx *gin.Context) {
	req := RequestJSON[[]entity.RegionProvince](ctx)
	err := c.regionService.SeedProvinces(ctx.Request.Context(), req)
	ResponseJSON[any](ctx, gin.H{"count": len(req)}, gin.H{"status": "ok"}, err)
}

func (c *regionController) SeedCities(ctx *gin.Context) {
	req := RequestJSON[[]entity.RegionCity](ctx)
	err := c.regionService.SeedCities(ctx.Request.Context(), req)
	ResponseJSON[any](ctx, gin.H{"count": len(req)}, gin.H{"status": "ok"}, err)
}

func (c *regionController) ListProvinces(ctx *gin.Context) {
	res, err := c.regionService.ListProvinces(ctx.Request.Context())
	ResponseJSON(ctx, gin.H{}, res, err)
}

func (c *regionController) ListCitiesByProvince(ctx *gin.Context) {
	type q struct {
		ProvinceID uint `json:"province_id" binding:"required"`
	}
	req := RequestJSON[q](ctx)
	res, err := c.regionService.ListCitiesByProvince(ctx.Request.Context(), req.ProvinceID)
	ResponseJSON(ctx, req, res, err)
}
