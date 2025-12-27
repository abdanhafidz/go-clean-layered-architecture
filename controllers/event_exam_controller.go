package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type EventExamController interface {
	Attempt(ctx *gin.Context)
	Answer(ctx *gin.Context)
	Submit(ctx *gin.Context)
	List(ctx *gin.Context)
}

type eventExamController struct {
	eventExamService services.EventExamService
}

func NewEventExamController(eventExamService services.EventExamService) EventExamController {
	return &eventExamController{
		eventExamService: eventExamService,
	}
}

// Exam Event Attempt godoc
// @Summary      Attempt Exam Event
// @Description  Start an attempt for a specific exam in an event
// @Tags         Exam Event
// @Accept       json
// @Produce      json
// @Param        event_slug  path      string  true  "Event Slug"
// @Param        exam_slug     path      string  true  "Exam Slug"
// @Success      200           {object}  dto.SuccessResponse[models.EventExamAttempt]
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/v1/events/{event_slug}/exam/{exam_slug}/attempt [get]
func (c *eventExamController) Attempt(ctx *gin.Context) {
	eventSlug := ctx.Param("event_slug")
	examSlug := ctx.Param("exam_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.eventExamService.AttemptEventExam(ctx.Request.Context(), eventSlug, examSlug, accountId)
	ResponseJSON(ctx, gin.H{"event_slug": eventSlug, "exam_slug": examSlug}, res, err)
}

// Answer Exam Event godoc
// @Summary      Answer Exam Event Question
// @Description  Submit an answer for a specific question in an exam attempt
// @Tags         Exam Event
// @Accept       json
// @Produce      json
// @Param        event_slug  path      string  true  "Event Slug"
// @Param        attempt_id    path      string  true  "Exam Attempt ID"
// @Param        request       body      dto.AnswerEventExamRequest  true  "Answer Exam Event Request"
// @Success      200           {object}  dto.SuccessResponse[dto.AnswerEventExamRequest]
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/v1/events/{event_slug}/exam/{attempt_id}/answer_question [post]
func (c *eventExamController) Answer(ctx *gin.Context) {
	eventSlug := ctx.Param("event_slug")
	attemptId := ParseUUID(ctx, "attempt_id")
	req := RequestJSON[dto.AnswerEventExamRequest](ctx)
	res, err := c.eventExamService.AnswerEventExam(ctx.Request.Context(), eventSlug, attemptId, req.QuestionId, req.Answer)
	ResponseJSON(ctx, gin.H{"cp_grader_result": res}, req, err)
}

// Submit Exam Event godoc
// @Summary      Submit Exam Event
// @Description  Submit the exam attempt for evaluation
// @Tags         Exam Event
// @Accept       json
// @Produce      json
// @Param        event_slug  path      string  true  "Event Slug"
// @Param        attempt_id    path      string  true  "Exam Attempt ID"
// @Success      200           {object}  dto.SuccessResponse[entity.Result]
// @Failure      400           {object}  dto.ErrorResponse
// @Router       /api/v1/events/{event_slug}/exam/{attempt_id}/submit [post]

func (c *eventExamController) Submit(ctx *gin.Context) {
	attemptId := ParseUUID(ctx, "attempt_id")
	res, err := c.eventExamService.SubmitEventExam(ctx.Request.Context(), attemptId)
	ResponseJSON(ctx, gin.H{}, res, err)
}

// List Exam by Event godoc
// @Summary      List Exams by Event
// @Description  Retrieve a list of exams associated with a specific event
// @Tags         Exam Event
// @Accept       json
// @Produce      json
// @Param        event_slug  path      string  true  "Event Slug"
// @Success      200         {object}   dto.SuccessResponse[[]models.Exam]
// @Failure      400         {object}  dto.ErrorResponse
// @Router       /api/v1/events/{event_slug}/exam [get]
func (c *eventExamController) List(ctx *gin.Context) {
	eventSlug := ctx.Param("event_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.eventExamService.ListExamByEvent(ctx.Request.Context(), eventSlug, accountId)
	ResponseJSON(ctx, gin.H{}, res, err)
}
