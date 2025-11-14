package controllers

import (
	"abdanhafidz.com/go-clean-layered-architecture/models/dto"
	"abdanhafidz.com/go-clean-layered-architecture/services"
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

func (c *optionController) CreateBulk(ctx *gin.Context) {
	var payload []dto.OptionsRequest
	payload = RequestJSON[[]dto.OptionsRequest](ctx)
	err := c.optionService.CreateBulk(ctx.Request.Context(), payload)
	ResponseJSON[any](ctx, payload, gin.H{"status": "ok"}, err)
}

func (c *optionController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	res, err := c.optionService.GetBySlug(ctx.Request.Context(), slug)
	ResponseJSON(ctx, gin.H{"slug": slug}, res, err)
}
