package controllers

import (
	"strconv"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type ExamController interface {
	CreateExam(ctx *gin.Context)
	UpdateExam(ctx *gin.Context)
	DeleteExam(ctx *gin.Context)
	ListExam(ctx *gin.Context)
	GetExamDetail(ctx *gin.Context)
	AssignExamToEvent(ctx *gin.Context)
	AssignExamToAcademy(ctx *gin.Context)
}

type examController struct {
	examService services.ExamService
}

func NewExamController(examService services.ExamService) ExamController {
	return &examController{examService}
}

// CreateExam godoc
// @Summary      Create Exam
// @Description  Create a new exam with configuration and proctoring settings
// @Tags         Exam
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateExamRequest  true  "Create Exam Request"
// @Success      200      {object}  dto.SuccessResponse[entity.Exam]
// @Failure      400      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/exam [post]
func (c *examController) CreateExam(ctx *gin.Context) {
	req := RequestJSON[dto.CreateExamRequest](ctx)
	res, err := c.examService.CreateExam(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

// UpdateExam godoc
// @Summary      Update Exam
// @Description  Update an existing exam
// @Tags         Exam
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Exam ID"
// @Param        request  body      dto.CreateExamRequest  true  "Update Exam Request"
// @Success      200      {object}  dto.SuccessResponse[entity.Exam]
// @Failure      400      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/exam/{id} [put]
func (c *examController) UpdateExam(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	req := RequestJSON[dto.CreateExamRequest](ctx)
	res, err := c.examService.UpdateExam(ctx.Request.Context(), id, req)
	ResponseJSON(ctx, req, res, err)
}

// DeleteExam godoc
// @Summary      Delete Exam
// @Description  Delete an existing exam
// @Tags         Exam
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Exam ID"
// @Success      200  {object}  dto.SuccessResponse[any]
// @Failure      400  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/exam/{id} [delete]
func (c *examController) DeleteExam(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	err := c.examService.DeleteExam(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}

// ListExam godoc
// @Summary      List Exams
// @Description  Retrieve a paginated list of exams
// @Tags         Exam
// @Accept       json
// @Produce      json
// @Param        limit   query     int     false  "Number of items per page"            default(10)
// @Param        page    query     int     false  "Page number"                         default(1)
// @Param        search  query     string  false  "Search term for exam title"
// @Param        sortBy  query     string  false  "Field to sort by"
// @Param        order   query     string  false  "Sort order (asc or desc)"
// @Success      200     {object}  dto.SuccessResponse[[]entity.Exam]
// @Failure      400     {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/exam [get]
func (c *examController) ListExam(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	search := ctx.DefaultQuery("search", "")
	sortBy := ctx.DefaultQuery("sortBy", "")
	order := ctx.DefaultQuery("order", "")

	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}

	if page < 1 {
		page = 1
	}

	offset := (page - 1) * limit
	p := entity.Pagination{Limit: limit, Offset: offset, Search: search, SortBy: sortBy, Order: order}

	list, total, err := c.examService.GetExamList(ctx.Request.Context(), p)

	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, err)
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if total == 0 {
		totalPages = 1
	}

	meta := gin.H{
		"totalItems":  total,
		"totalPages":  totalPages,
		"currentPage": page,
	}

	ResponseJSON(ctx, meta, list, err)
}

// GetExamDetail godoc
// @Summary      Get Exam Detail
// @Description  Retrieve detailed exam information
// @Tags         Exam
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Exam ID"
// @Success      200  {object}  dto.SuccessResponse[entity.Exam]
// @Failure      400  {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/exam/{id} [get]
func (c *examController) GetExamDetail(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	res, err := c.examService.GetExamDetail(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, res, err)
}

// AssignExamToEvent godoc
// @Summary      Assign Exam to Event
// @Description  Assign an exam to an event
// @Tags         Exam
// @Accept       json
// @Produce      json
// @Param        exam_id   path      string  true  "Exam ID"
// @Param        event_id  path      string  true  "Event ID"
// @Success      200       {object}  dto.SuccessResponse[any]
// @Failure      400       {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/exam/{exam_id}/event/{event_id} [post]
func (c *examController) AssignExamToEvent(ctx *gin.Context) {
	examId := ParseUUID(ctx, "exam_id")
	eventId := ParseUUID(ctx, "event_id")

	err := c.examService.AssignExamToEvent(ctx.Request.Context(), examId, eventId)
	ResponseJSON(ctx, gin.H{"exam_id": examId, "event_id": eventId}, gin.H{"assigned": true}, err)
}

// AssignExamToAcademy godoc
// @Summary      Assign Exam to Academy
// @Description  Assign an exam to an academy
// @Tags         Exam
// @Accept       json
// @Produce      json
// @Param        exam_id     path      string  true  "Exam ID"
// @Param        academy_id  path      string  true  "Academy ID"
// @Success      200         {object}  dto.SuccessResponse[any]
// @Failure      400         {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/exam/{exam_id}/academy/{academy_id} [post]
func (c *examController) AssignExamToAcademy(ctx *gin.Context) {
	examId := ParseUUID(ctx, "exam_id")
	academyId := ParseUUID(ctx, "academy_id")

	err := c.examService.AssignExamToAcademy(ctx.Request.Context(), examId, academyId)
	ResponseJSON(ctx, gin.H{"exam_id": examId, "academy_id": academyId}, gin.H{"assigned": true}, err)
}
