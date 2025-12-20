package controllers

import (
    "abdanhafidz.com/go-boilerplate/models/dto"
    "abdanhafidz.com/go-boilerplate/services"
    "abdanhafidz.com/go-boilerplate/utils"
    "github.com/gin-gonic/gin"
)

type AcademyExamController interface {
    Attempt(ctx *gin.Context)
    Answer(ctx *gin.Context)
    Submit(ctx *gin.Context)
    List(ctx *gin.Context)
}

type academyExamController struct { academyExamService services.AcademyExamService }

func NewAcademyExamController(academyExamService services.AcademyExamService) AcademyExamController { return &academyExamController{ academyExamService: academyExamService } }

func (c *academyExamController) Attempt(ctx *gin.Context) {
    academySlug := ctx.Param("academy_slug")
    examSlug := ctx.Param("exam_slug")
    accountId := ParseAccountId(ctx)
    res, err := c.academyExamService.AttemptExamAcademy(ctx.Request.Context(), academySlug, examSlug, accountId)
    ResponseJSON(ctx, gin.H{"academy_slug": academySlug, "exam_slug": examSlug}, res, err)
}

func (c *academyExamController) Answer(ctx *gin.Context) {
    academySlug := ctx.Param("academy_slug")
    attemptId, _ := utils.ToUUID(ctx.Param("attempt_id"))
    req := RequestJSON[dto.AnswerExamEventRequest](ctx)
    res, err := c.academyExamService.AnswerExamAcademy(ctx.Request.Context(), academySlug, attemptId, req.QuestionId, req.Answer)
    ResponseJSON(ctx, gin.H{"cp_grader_result": res}, req, err)
}

func (c *academyExamController) Submit(ctx *gin.Context) {
    attemptId, _ := utils.ToUUID(ctx.Param("attempt_id"))
    res, err := c.academyExamService.SubmitExamAcademy(ctx.Request.Context(), attemptId)
    ResponseJSON(ctx, gin.H{}, res, err)
}

func (c *academyExamController) List(ctx *gin.Context) {
    academySlug := ctx.Param("academy_slug")
    accountId := ParseAccountId(ctx)
    res, err := c.academyExamService.ListExamByAcademy(ctx.Request.Context(), academySlug, accountId)
    ResponseJSON(ctx, gin.H{}, res, err)
}