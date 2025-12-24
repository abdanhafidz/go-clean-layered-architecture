package controllers

import (
	"github.com/gin-gonic/gin"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/services"
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

// Seed Provinces godoc
// @Summary      Seed Provinces
// @Description  Seed multiple provinces into the system
// @Tags         Region
// @Accept       json
// @Produce      json
// @Param        request  body      []entity.RegionProvince  true  "Seed Provinces Request"
// @Success      200      {object}  dto.SuccessResponse[any]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/options/region/seed-provinces [post]
func (c *regionController) SeedProvinces(ctx *gin.Context) {
	req := RequestJSON[[]entity.RegionProvince](ctx)
	err := c.regionService.SeedProvinces(ctx.Request.Context(), req)
	x := dto.SuccessResponse[any]{}
	x.Data = nil
	ResponseJSON(ctx, gin.H{"count": len(req)}, gin.H{"status": "ok"}, err)

}

// Seed Cities godoc
// @Summary      Seed Cities
// @Description  Seed multiple cities into the system
// @Tags         Region
// @Accept       json
// @Produce      json
// @Param        request  body      []entity.RegionCity  true  "Seed Cities Request"
// @Success      200      {object}  dto.SuccessResponse[[]entity.RegionCity]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/options/region/seed-cities [post]
func (c *regionController) SeedCities(ctx *gin.Context) {
	req := RequestJSON[[]entity.RegionCity](ctx)
	err := c.regionService.SeedCities(ctx.Request.Context(), req)
	ResponseJSON[any](ctx, gin.H{"count": len(req)}, gin.H{"status": "ok"}, err)
}

// List Provinces godoc
// @Summary      List Provinces
// @Description  Retrieve a list of all provinces
// @Tags         Region
// @Accept       json
// @Produce      json
// @Success      200  {object}   dto.SuccessResponse[[]entity.RegionProvince]
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/v1/options/region/provinces [get]
func (c *regionController) ListProvinces(ctx *gin.Context) {
	res, err := c.regionService.ListProvinces(ctx.Request.Context())
	ResponseJSON(ctx, gin.H{}, res, err)
}

// List Cities By Province godoc
// @Summary      List Cities by Province
// @Description  Retrieve a list of cities within a specific province
// @Tags         Region
// @Accept       json
// @Produce      json
// @Param        request  body      object{province_id=uint}  true  "Province ID"
// @Success      200      {object}   dto.SuccessResponse[[]entity.RegionCity]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/options/region/cities [get]
func (c *regionController) ListCitiesByProvince(ctx *gin.Context) {
	type q struct {
		ProvinceID uint `json:"province_id" binding:"required"`
	}
	req := RequestJSON[q](ctx)
	res, err := c.regionService.ListCitiesByProvince(ctx.Request.Context(), req.ProvinceID)
	ResponseJSON(ctx, req, res, err)
}
