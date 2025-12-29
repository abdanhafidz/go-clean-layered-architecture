package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	"abdanhafidz.com/go-boilerplate/utils"
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
// GetAcademy godoc
// @Summary      Get Academy by Slug
// @Description  Retrieve academy details using its slug
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        academy_slug  path      string  true  "Academy Slug"
// @Success      200           {object}  dto.SuccessResponse[any]
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/v1/academy/{academy_slug} [get]
func (c *academyController) GetAcademy(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.academyService.GetAcademyResponse(ctx.Request.Context(), accountId, academySlug)
	ResponseJSON(ctx, gin.H{"academy_slug": academySlug}, res, err)
}

// GetAcademyDetail godoc
// @Summary      Get Academy Detail by ID
// @Description  Retrieve detailed academy information using its ID
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Academy ID"
// @Success      200  {object}  dto.SuccessResponse[entity.Academy]
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/v1/admin/academy/id/{id}/detail [get]
func (c *academyController) GetAcademyDetail(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	res, err := c.academyService.GetAcademyDetail(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, res, err)
}

// ListAcademy godoc
// @Summary      List Academies
// @Description  Retrieve a paginated list of academies with optional filters
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        limit           query     int     false  "Number of items per page"            default(10)
// @Param        page            query     int     false  "Page number"                         default(1)
// @Param        search          query     string  false  "Search term for academy title/name"
// @Param        sortBy          query     string  false  "Field to sort by"
// @Param        order           query     string  false  "Sort order (asc or desc)"
// @Param        registerStatus  query     int     false  "Filter by registration status"
// @Param        status          query     string  false  "Filter by academy status"
// @Success      200             {object}  dto.SuccessResponse[[]entity.Academy]
// @Failure	  400             {object}  dto.ErrorResponse
// @Router       /api/v1/academy [get]

func (c *academyController) ListAcademy(ctx *gin.Context) {
	accountId := ParseAccountId(ctx)

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	search := ctx.DefaultQuery("search", "")
	sortBy := ctx.DefaultQuery("sortBy", "")
	order := ctx.DefaultQuery("order", "")

	var registerStatus *int
	if val := ctx.Query("registerStatus"); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			registerStatus = &i
		}
	}
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

	var status *string
	if val := ctx.Query("status"); val != "" {
		if val == entity.StatusNotStarted || val == entity.StatusInProgress || val == entity.StatusFinished {
			status = &val
		}
	}

	offset := (page - 1) * limit
	p := entity.Pagination{Limit: limit, Offset: offset, Search: search, SortBy: sortBy, Order: order, RegisterStatus: registerStatus, Status: status}
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

// CreateAcademy godoc
// @Summary      Create Academy
// @Description  Create a new academy
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateAcademyRequest  true  "Create Academy Request"
// @Success      200      {object}  dto.SuccessResponse[entity.Academy]
// @Failure      400      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy [post]
func (c *academyController) CreateAcademy(ctx *gin.Context) {
	req := RequestJSON[dto.CreateAcademyRequest](ctx)
	res, err := c.academyService.CreateAcademy(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

// UpdateAcademy godoc
// @Summary      Update Academy
// @Description  Update an existing academy
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "Academy ID"
// @Param        request  body      dto.UpdateAcademyRequest  true  "Update Academy Request"
// @Success      200      {object}  dto.SuccessResponse[entity.Academy]
// @Failure      400      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/id/{id} [put]
func (c *academyController) UpdateAcademy(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	req := RequestJSON[dto.UpdateAcademyRequest](ctx)
	res, err := c.academyService.UpdateAcademy(ctx.Request.Context(), id, req)
	ResponseJSON(ctx, req, res, err)
}

// DeleteAcademy godoc
// @Summary      Delete Academy
// @Description  Delete an existing academy
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Academy ID"
// @Success      200  {object}  dto.SuccessResponse[any]
// @Failure      400  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/id/{id} [delete]
func (c *academyController) DeleteAcademy(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	err := c.academyService.DeleteAcademy(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

// ================= MATERIAL =================
// GetMaterial godoc
// @Summary      Get Material by Slug
// @Description  Retrieve material details using its slug
// @Tags         Material
// @Accept       json
// @Produce      json
// @Param        academy_slug   path      string  true  "Academy Slug"
// @Param        material_slug  path      string  true  "Material Slug"
// @Success      200            {object}  dto.SuccessResponse[dto.MaterialDetailResponse]
// @Failure      400            {object}  dto.ErrorResponse
// @Router       /api/v1/academy/{academy_slug}/{material_slug} [get]
func (c *academyController) GetMaterial(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	materialSlug := ctx.Param("material_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.academyService.GetMaterialResponse(ctx.Request.Context(), accountId, academySlug, materialSlug)
	ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "material_slug": materialSlug}, res, err)
}

// CreateMaterial godoc
// @Summary      Create Material
// @Description  Create a new material for an academy
// @Tags         Material
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateMaterialRequest  true  "Create Material Request"
// @Success      200      {object}  dto.SuccessResponse[entity.AcademyMaterial]
// @Failure      400      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/materials [post]
func (c *academyController) CreateMaterial(ctx *gin.Context) {
	req := RequestJSON[dto.CreateMaterialRequest](ctx)
	res, err := c.academyService.CreateMaterial(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

// DeleteMaterial godoc
// @Summary      Delete Material
// @Description  Delete an existing material
// @Tags         Material
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Material ID"
// @Success      200  {object}  dto.SuccessResponse[any]
// @Failure      400  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/materials/{id} [delete]

func (c *academyController) DeleteMaterial(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	err := c.academyService.DeleteMaterial(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

// ================= CONTENT =================
// GetContent godoc
// @Summary      Get Content by Order
// @Description  Retrieve content details using its order within a material
// @Tags         Content
// @Accept       json
// @Produce      json
// @Param        academy_slug   path      string  true  "Academy Slug"
// @Param        material_slug  path      string  true  "Material Slug"
// @Param        order          path      int     true  "Content Order"
// @Success      200            {object}  dto.SuccessResponse[entity.AcademyContent]
// @Failure      400            {object}  dto.ErrorResponse
// @Router       /api/v1/academy/{academy_slug}/{material_slug}/{order} [get]

func (c *academyController) GetContent(ctx *gin.Context) {
	accountId := ParseAccountId(ctx)
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

// CreateContent godoc
// @Summary      Create Content
// @Description  Create a new content for a material
// @Tags         Content
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateContentRequest  true  "Create Content Request"
// @Success      200      {object}  dto.SuccessResponse[entity.AcademyContent]
// @Failure      400      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/contents [post]
func (c *academyController) CreateContent(ctx *gin.Context) {
	req := RequestJSON[dto.CreateContentRequest](ctx)
	res, err := c.academyService.CreateContent(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

// DeleteContent godoc
// @Summary      Delete Content
// @Description  Delete an existing content
// @Tags         Content
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Content ID"
// @Success      200  {object}  dto.SuccessResponse[any]
// @Failure      400  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/contents/{id} [delete]
func (c *academyController) DeleteContent(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	err := c.academyService.DeleteContent(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

// ================= PROGRESS =================
// UpdateContentProgress godoc
// @Summary      Update Content Progress
// @Description  Update the progress of a content within a material
// @Tags         Progress
// @Accept       json
// @Produce      json
// @Param        academy_slug   path      string  true  "Academy Slug"
// @Param        material_slug  path      string  true  "Material Slug"
// @Param        order          path      int     true  "Content Order"
// @Success      200            {object}  dto.SuccessResponse[any]
// @Failure      400            {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/academy/{academy_slug}/{material_slug}/{order} [post]
func (c *academyController) UpdateContentProgress(ctx *gin.Context) {
	accountId := ParseAccountId(ctx)
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

// JoinAcademyByCode godoc
// @Summary      Join Academy by Code
// @Description  Join an academy using a unique code
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        request  body      dto.JoinAcademyByCodeRequest  true  "Join Academy by Code Request"
// @Success      200      {object}  dto.SuccessResponse[dto.AcademyMiniDetailResponse]
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      402      {object}  dto.SuccessResponse[entity.AcademyPaymentTransaction]
// @Security     BearerAuth
// @Router       /api/v1/academy/join [post]
func (c *academyController) JoinAcademyByCode(ctx *gin.Context) {
	req := RequestJSON[dto.JoinAcademyByCodeRequest](ctx)

	if err := utils.ValidateCode(req.Code); err != nil {
		ResponseJSON[any, any](ctx, req, nil, http_error.INVALID_CODE)
		return
	}

	accountId := ParseAccountId(ctx)
	res, err := c.academyService.JoinByCode(ctx.Request.Context(), accountId, req.Code)

	if errors.Is(err, http_error.PAYMENT_REQUIRED) {
		ResponseJSON(ctx, res.Payment, res, err)
		return
	}

	ResponseJSON(ctx, req, res, err)
}

// AssignAccountToAcademy godoc
// @Summary      Assign Account to Academy
// @Description  Assign an account to an academy
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AssignRequest  true  "Assign Account to Academy Request"
// @Success      200      {object}  dto.SuccessResponse[entity.AcademyAssign]
// @Failure      400      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/assign [post]
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

// UnassignAccountFromAcademy godoc
// @Summary      Unassign Account from Academy
// @Description  Unassign an account from an academy
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Assignment ID"
// @Success      200  {object}  dto.SuccessResponse[any]
// @Failure      400  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/assign/{id} [delete]
func (c *academyController) UnassignAccountFromAcademy(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	err := c.academyService.UnassignAccountFromAcademy(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

// ListAssignmentsByAcademy godoc
// @Summary      List Assignments by Academy
// @Description  Retrieve a list of assignments for a specific academy
// @Tags         Academy
// @Accept       json
// @Produce      json
// @Param        academy_id  path	  string  true  "Academy ID"
// @Success      200         {object}  dto.SuccessResponse[[]entity.AcademyAssign]
// @Failure      400         {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/academy/assign/{academy_id} [get]
func (c *academyController) ListAssignmentsByAcademy(ctx *gin.Context) {
	academyId := ParseUUID(ctx, "academy_id")
	res, err := c.academyService.ListAssignmentsByAcademy(ctx.Request.Context(), academyId)
	ResponseJSON(ctx, gin.H{"academy_id": academyId}, res, err)
}
