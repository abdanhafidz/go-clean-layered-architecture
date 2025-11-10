package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
)

type ExamController interface {
	Attempt(ctx *gin.Context)
	Answer(ctx *gin.Context)
	Submit(ctx *gin.Context)
}

type examController struct {
	examService services.ExamService
}

func NewExamController(examService services.ExamService) ExamController {
	return &examController{
		examService: examService,
	}
}

func (c *examController) Attempt(ctx *gin.Context) {
	eventSlug := ctx.Param("event_slug")
	examSlug := ctx.Param("exam_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.examService.AttemptExamEvent(ctx.Request.Context(), eventSlug, examSlug, accountId)
	ResponseJSON(ctx, gin.H{"event_slug": eventSlug, "exam_slug": examSlug}, res, err)
}

func (c *examController) Answer(ctx *gin.Context) {
	attemptId, _ := utils.ToUUID(ctx.Param("attempt_id"))
	req := RequestJSON[dto.AnswerExamEventRequest](ctx)
	res, err := c.examService.AnswerExamEvent(ctx.Request.Context(), attemptId, req.QuestionId, req.Answer)
	ResponseJSON(ctx, gin.H{"cp_grader_result": res}, req, err)
}

func (c *examController) Submit(ctx *gin.Context) {
	attemptId, _ := utils.ToUUID(ctx.Param("attempt_id"))
	res, err := c.examService.SubmitExamEvent(ctx.Request.Context(), attemptId)
	ResponseJSON(ctx, gin.H{}, res, err)
}
