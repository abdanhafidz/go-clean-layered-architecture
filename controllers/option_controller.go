package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type OptionController interface {
	CreateBulk(ctx *gin.Context)
	GetBySlug(ctx *gin.Context)
}

type optionController struct {
	optionService services.OptionService
}

func NewOptionController(optionService services.OptionService) OptionController {
	return &optionController{optionService: optionService}
}

// Create Bulk Options godoc
// @Summary      Create Bulk Options
// @Description  Create multiple options in bulk
// @Tags         Option
// @Accept       json
// @Produce      json
// @Param        request  body      []dto.OptionsRequest  true  "Bulk Options Request"
// @Success      200      {object}  dto.SuccessResponse[any]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/options/create [post]
func (c *optionController) CreateBulk(ctx *gin.Context) {
	var payload []dto.OptionsRequest
	payload = RequestJSON[[]dto.OptionsRequest](ctx)
	err := c.optionService.CreateBulk(ctx.Request.Context(), payload)
	ResponseJSON[any](ctx, payload, gin.H{"status": "ok"}, err)
}

// Get Option By Slug godoc
// @Summary      Get Option by Slug
// @Description  Retrieve option details using its slug
// @Tags         Option
// @Accept       json
// @Produce      json
// @Param        slug  path      string  true  "Option Slug"
// @Success      200   {object}  dto.SuccessResponse[models.Options]
// @Failure      400   {object}  dto.ErrorResponse
// @Router       /api/v1/options/list/{slug} [get]
func (c *optionController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	res, err := c.optionService.GetBySlug(ctx.Request.Context(), slug)
	ResponseJSON(ctx, gin.H{"slug": slug}, res, err)
}
