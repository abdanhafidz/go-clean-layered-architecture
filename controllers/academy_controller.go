package controllers

import (
	"fmt"
	"net/http"
	"strconv"

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

	GetMaterial(ctx *gin.Context)
	GetContent(ctx *gin.Context)
	UpdateContentProgress(ctx *gin.Context)
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
	academySlug := ctx.Param("academy_slug")
	accountIdVal := ctx.Value("account_id")
	accountId, _ := accountIdVal.(string)
	res, err := c.academyService.GetAcademy(ctx.Request.Context(), accountId, academySlug)

	ResponseJSON(ctx, gin.H{"academy_slug": academySlug}, res, err)
}

func (c *academyController) GetAcademyDetail(ctx *gin.Context) {
	id, _ := uuid.Parse(ctx.Param("id"))
	res, err := c.academyService.GetAcademyDetail(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, res, err)
}

func (c *academyController) ListAcademies(ctx *gin.Context) {
	accountIdVal := ctx.Value("account_id")
	accountId, _ := accountIdVal.(string)
	fmt.Println("Account ID in ListAcademies:", accountId)
	res, err := c.academyService.ListAcademies(ctx.Request.Context(), accountId)
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

// MATERIAL

func (c *academyController) CreateMaterial(ctx *gin.Context) {
	req := RequestJSON[dto.CreateMaterialRequest](ctx)

	res, err := c.academyService.CreateMaterial(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) GetMaterial(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	materialSlug := ctx.Param("material_slug")
	res, err := c.academyService.GetMaterial(ctx.Request.Context(), academySlug, materialSlug)
	ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "material_slug": materialSlug}, res, err)
}

// CONTENT

func (c *academyController) CreateContent(ctx *gin.Context) {
	req := RequestJSON[dto.CreateContentRequest](ctx)
	res, err := c.academyService.CreateContent(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) GetContent(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	materialSlug := ctx.Param("material_slug")
	orderString := ctx.Param("order")

	orderID64, err := strconv.ParseUint(orderString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'order' parameter. Must be a positive integer."})
		return
	}
	order := uint(orderID64)

	res, err := c.academyService.GetContent(ctx.Request.Context(), academySlug, materialSlug, order)
	ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "material_slug": materialSlug, "content_slug": order}, res, err)
}

func (c *academyController) UpdateContentProgress(ctx *gin.Context) {
	accountIdVal := ctx.Value("account_id")
	academySlug := ctx.Param("academy_slug")
	materialSlug := ctx.Param("material_slug")
	orderString := ctx.Param("order")
	accountId, _ := accountIdVal.(string)

	orderID64, err := strconv.ParseUint(orderString, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'order' parameter. Must be a positive integer."})
		return
	}
	order := uint(orderID64)

	contentProgress, materialProgress, academyProgress, err := c.academyService.UpdateContentProgress(ctx, accountId, academySlug, materialSlug, order)
	res := gin.H{
		"content_progress":  contentProgress,
		"material_progress": materialProgress,
		"academy_progress":  academyProgress,
	}

	ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "material_slug": materialSlug, "content_slug": order}, res, err)
}
