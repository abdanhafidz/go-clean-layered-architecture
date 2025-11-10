package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AcademyController interface {
	CreateAcademy(ctx *gin.Context)
	GetAcademy(ctx *gin.Context)
	GetAcademyDetail(ctx *gin.Context)
	ListAcademies(ctx *gin.Context)
	UpdateAcademy(ctx *gin.Context)
	DeleteAcademy(ctx *gin.Context)

	CreateMaterial(ctx *gin.Context)
	CreateContent(ctx *gin.Context)
}

type academyController struct {
	academyService services.AcademyService
}

func NewAcademyController(academyService services.AcademyService) AcademyController {
	return &academyController{academyService}
}

func (c *academyController) CreateAcademy(ctx *gin.Context) {
	req := RequestJSON[dto.CreateAcademyRequest](ctx)
	res, err := c.academyService.CreateAcademy(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) GetAcademy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	res, err := c.academyService.GetAcademy(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, res, err)
}

func (c *academyController) GetAcademyDetail(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	res, err := c.academyService.GetAcademyDetail(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, res, err)
}

func (c *academyController) ListAcademies(ctx *gin.Context) {
	res, err := c.academyService.ListAcademies(ctx.Request.Context())
	ResponseJSON(ctx, gin.H{}, res, err)
}

func (c *academyController) UpdateAcademy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	req := RequestJSON[dto.UpdateAcademyRequest](ctx)
	res, err := c.academyService.UpdateAcademy(ctx.Request.Context(), id, req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) DeleteAcademy(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	err := c.academyService.DeleteAcademy(ctx.Request.Context(), id)

	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

func (c *academyController) CreateMaterial(ctx *gin.Context) {
	req := RequestJSON[dto.CreateMaterialRequest](ctx)

	res, err := c.academyService.CreateMaterial(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) CreateContent(ctx *gin.Context) {
	req := RequestJSON[dto.CreateContentRequest](ctx)

	res, err := c.academyService.CreateContent(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}
