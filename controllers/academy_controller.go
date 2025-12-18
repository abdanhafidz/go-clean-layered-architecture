package controllers

import (
	"net/http"
	"strconv"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AcademyController interface {
	// Academy
	CreateAcademy(ctx *gin.Context)
	GetAcademy(ctx *gin.Context)
	GetAcademyDetail(ctx *gin.Context)
	ListAcademy(ctx *gin.Context)
	UpdateAcademy(ctx *gin.Context)
	DeleteAcademy(ctx *gin.Context)
	JoinAcademyByCode(ctx *gin.Context)
	AssignAccountToAcademy(ctx *gin.Context)
	UnassignAccountFromAcademy(ctx *gin.Context)
	ListAssignmentsByAcademy(ctx *gin.Context)

	// Material
	GetMaterial(ctx *gin.Context)
	CreateMaterial(ctx *gin.Context)
	DeleteMaterial(ctx *gin.Context)

	// Content
	CreateContent(ctx *gin.Context)
	GetContent(ctx *gin.Context)
	DeleteContent(ctx *gin.Context)

	// Progress
	UpdateContentProgress(ctx *gin.Context)
}

type academyController struct {
	academyService services.AcademyService
}

func NewAcademyController(academyService services.AcademyService) AcademyController {
	return &academyController{academyService}
}

// ================= ACADEMY =================

func (c *academyController) GetAcademy(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	accountIdStr := ctx.GetString("account_id")
	accountId, err := uuid.Parse(accountIdStr)
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.UNAUTHORIZED)
		return
	}

	res, err := c.academyService.GetAcademyResponse(ctx.Request.Context(), accountId, academySlug)
	ResponseJSON(ctx, gin.H{"academy_slug": academySlug}, res, err)
}

func (c *academyController) GetAcademyDetail(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}

	res, err := c.academyService.GetAcademyDetail(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, res, err)
}

func (c *academyController) ListAcademy(ctx *gin.Context) {
	accountIdStr := ctx.GetString("account_id")
	accountId, err := uuid.Parse(accountIdStr)
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.UNAUTHORIZED)
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	search := ctx.DefaultQuery("search", "")
	sortBy := ctx.DefaultQuery("sortBy", "")
	order := ctx.DefaultQuery("order", "")
	isModified := false

	if limit < 1 {
		limit = 10
		isModified = true
	} else if limit > 50 {
		limit = 50
		isModified = true
	}

	if page < 1 {
		page = 1
		isModified = true
	}

	offset := (page - 1) * limit
	p := entity.Pagination{Limit: limit, Offset: offset, Search: search, SortBy: sortBy, Order: order}
	list, total, err := c.academyService.ListAcademy(ctx.Request.Context(), accountId, p)

	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, err)
		return
	}

	var totalPages int
	if total == 0 {
		totalPages = 1
	} else {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	if page > totalPages {
		page = totalPages
		offset = (page - 1) * limit
		p.Offset = offset
		list, total, err = c.academyService.ListAcademy(ctx.Request.Context(), accountId, p)
		isModified = true
	}

	meta := gin.H{
		"totalItems":  total,
		"totalPages":  totalPages,
		"currentPage": page,
	}
	if isModified {
		ctx.Status(http.StatusAccepted)
	}

	ResponseJSON(ctx, meta, list, err)
}

func (c *academyController) CreateAcademy(ctx *gin.Context) {
	req := RequestJSON[dto.CreateAcademyRequest](ctx)
	res, err := c.academyService.CreateAcademy(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) UpdateAcademy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}

	req := RequestJSON[dto.UpdateAcademyRequest](ctx)
	res, err := c.academyService.UpdateAcademy(ctx.Request.Context(), id, req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) DeleteAcademy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}

	err = c.academyService.DeleteAcademy(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

// ================= MATERIAL =================

func (c *academyController) GetMaterial(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	materialSlug := ctx.Param("material_slug")

	accountIdStr := ctx.GetString("account_id")
	accountId, err := uuid.Parse(accountIdStr)
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.UNAUTHORIZED)
		return
	}

	res, err := c.academyService.GetMaterialResponse(ctx.Request.Context(), accountId, academySlug, materialSlug)
	ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "material_slug": materialSlug}, res, err)
}

func (c *academyController) CreateMaterial(ctx *gin.Context) {
	req := RequestJSON[dto.CreateMaterialRequest](ctx)
	res, err := c.academyService.CreateMaterial(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) DeleteMaterial(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}

	err = c.academyService.DeleteMaterial(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

// ================= CONTENT =================

func (c *academyController) GetContent(ctx *gin.Context) {
	accountIdStr := ctx.GetString("account_id")
	accountId, err := uuid.Parse(accountIdStr)
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.UNAUTHORIZED)
		return
	}

	academySlug := ctx.Param("academy_slug")
	materialSlug := ctx.Param("material_slug")

	orderID64, err := strconv.ParseUint(ctx.Param("order"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'order' parameter. Must be a positive integer."})
		return
	}
	order := uint(orderID64)

	res, err := c.academyService.GetContent(ctx.Request.Context(), accountId, academySlug, materialSlug, order)
	ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "material_slug": materialSlug, "content_order": order}, res, err)
}

func (c *academyController) CreateContent(ctx *gin.Context) {
	req := RequestJSON[dto.CreateContentRequest](ctx)
	res, err := c.academyService.CreateContent(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) DeleteContent(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}

	err = c.academyService.DeleteContent(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

// ================= PROGRESS =================

func (c *academyController) UpdateContentProgress(ctx *gin.Context) {
	accountIdStr := ctx.GetString("account_id")
	accountId, err := uuid.Parse(accountIdStr)
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.UNAUTHORIZED)
		return
	}

	academySlug := ctx.Param("academy_slug")
	materialSlug := ctx.Param("material_slug")

	orderID64, err := strconv.ParseUint(ctx.Param("order"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'order' parameter. Must be a positive integer."})
		return
	}
	order := uint(orderID64)

	contentProgress, materialProgress, academyProgress, err := c.academyService.UpdateContentProgress(ctx.Request.Context(), accountId, academySlug, materialSlug, order)

	res := gin.H{
		"content_progress":  contentProgress,
		"material_progress": materialProgress,
		"academy_progress":  academyProgress,
	}

	ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "material_slug": materialSlug, "content_order": order}, res, err)
}

func (c *academyController) JoinAcademyByCode(ctx *gin.Context) {
	req := RequestJSON[dto.JoinAcademyByCodeRequest](ctx)
	accountId := ParseAccountId(ctx)
	res, err := c.academyService.JoinByCode(ctx.Request.Context(), accountId, req.Code)
	ResponseJSON(ctx, req, res, err)
}
func (c *academyController) AssignAccountToAcademy(ctx *gin.Context) {
	req := RequestJSON[dto.AssignRequest](ctx)
	academyId, errA := uuid.Parse(req.AcademyId)
	accountId, errB := uuid.Parse(req.AccountId)
	if errA != nil || errB != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}
	res, err := c.academyService.AssignAccountToAcademy(ctx.Request.Context(), academyId, accountId)
	ResponseJSON(ctx, req, res, err)
}

func (c *academyController) UnassignAccountFromAcademy(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}
	err = c.academyService.UnassignAccountFromAcademy(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

func (c *academyController) ListAssignmentsByAcademy(ctx *gin.Context) {
	academyId, err := uuid.Parse(ctx.Param("academy_id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}
	res, err := c.academyService.ListAssignmentsByAcademy(ctx.Request.Context(), academyId)
	ResponseJSON(ctx, gin.H{"academy_id": academyId}, res, err)
}
