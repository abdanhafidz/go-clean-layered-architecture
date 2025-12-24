package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type AcademyExamController interface {
	Attempt(ctx *gin.Context)
	Answer(ctx *gin.Context)
	Submit(ctx *gin.Context)
	List(ctx *gin.Context)
}

type academyExamController struct{ academyExamService services.AcademyExamService }

func NewAcademyExamController(academyExamService services.AcademyExamService) AcademyExamController {
	return &academyExamController{academyExamService: academyExamService}
}

// Attempt godoc
// @Summary      Attempt Academy Exam
// @Description  Start an attempt for a specific exam in an academy
// @Tags         Academy Exam
// @Accept       json
// @Produce      json
// @Param        academy_slug  path      string  true  "Academy Slug"
// @Param        exam_slug     path      string  true  "Exam Slug"
// @Success      200           {object}  dto.SuccessResponse[models.ExamAcademyAttempt]
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/v1/academy/{academy_slug}/exam/{exam_slug}/attempt [get]
func (c *academyExamController) Attempt(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	examSlug := ctx.Param("exam_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.academyExamService.AttemptExamAcademy(ctx.Request.Context(), academySlug, examSlug, accountId)
	ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "exam_slug": examSlug}, res, err)
}

// Answer godoc
// @Summary      Answer Academy Exam Question
// @Description  Submit an answer for a specific question in an exam attempt
// @Tags         Academy Exam
// @Accept       json
// @Produce      json
// @Param        academy_slug  path      string  true  "Academy Slug"
// @Param        attempt_id    path      string  true  "Exam Attempt ID"
// @Param        request       body      dto.AnswerExamEventRequest  true  "Answer Exam Event Request"
// @Success      200           {object}  dto.SuccessResponse[any]
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/v1/academy/{academy_slug}/exam/{attempt_id}/answer_question [post]

func (c *academyExamController) Answer(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	attemptId := ParseUUID(ctx, "attempt_id")
	req := RequestJSON[dto.AnswerExamEventRequest](ctx)
	res, err := c.academyExamService.AnswerExamAcademy(ctx.Request.Context(), academySlug, attemptId, req.QuestionId, req.Answer)
	ResponseJSON(ctx, gin.H{"cp_grader_result": res}, req, err)
}

// Submit godoc
// @Summary      Submit Academy Exam
// @Description  Submit the exam attempt for evaluation
// @Tags         Academy Exam
// @Accept       json
// @Produce      json
// @Param        academy_slug  path      string  true  "Academy Slug"
// @Param        attempt_id    path      string  true  "Exam Attempt ID"
// @Success      200           {object}  dto.SuccessResponse[entity.ExamAcademyResult]
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/v1/academy/{academy_slug}/exam/{attempt_id}/submit [post]

func (c *academyExamController) Submit(ctx *gin.Context) {
	attemptId := ParseUUID(ctx, "attempt_id")
	res, err := c.academyExamService.SubmitExamAcademy(ctx.Request.Context(), attemptId)
	ResponseJSON(ctx, gin.H{}, res, err)
}

// List godoc
// @Summary      List Academy Exams
// @Description  Retrieve a list of exams available in a specific academy
// @Tags         Academy Exam
// @Accept       json
// @Produce      json
// @Param        academy_slug  path      string  true  "Academy Slug"
// @Success      200           {object}   dto.SuccessResponse[[]entity.Exam]
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/v1/academy/{academy_slug}/exam [get]

func (c *academyExamController) List(ctx *gin.Context) {
	academySlug := ctx.Param("academy_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.academyExamService.ListExamByAcademy(ctx.Request.Context(), academySlug, accountId)
	ResponseJSON(ctx, gin.H{}, res, err)
}
