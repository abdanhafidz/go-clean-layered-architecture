package controllers

import (
	"strconv"

	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventController interface {
	List(ctx *gin.Context)
	DetailBySlug(ctx *gin.Context)
	Join(ctx *gin.Context)
	QuizListByEvent(ctx *gin.Context)
}

type eventController struct {
	eventService services.EventService
}

func NewEventController(eventService services.EventService) EventController {
	return &eventController{eventService: eventService}
}

func (c *eventController) List(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	p := repositories.Pagination{Limit: limit, Offset: offset}
	list, total, err := c.eventService.List(ctx.Request.Context(), p)
	meta := gin.H{
		"total_records": total,
		"page_size":     limit,
		"current_page":  (offset / limit) + 1,
	}
	ResponseJSON(ctx, meta, list, err)
}

func (c *eventController) DetailBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	accStr := ctx.Query("account_id")
	accountId, _ := uuid.Parse(accStr)
	res, err := c.eventService.DetailBySlug(ctx.Request.Context(), slug, accountId)
	ResponseJSON(ctx, gin.H{"slug": slug, "id_user": accountId}, res, err)
}

func (c *eventController) Join(ctx *gin.Context) {
	req := RequestJSON[dto.JoinEventRequest](ctx)
	accStr := ctx.Query("account_id")
	accountId, _ := uuid.Parse(accStr)
	res, err := c.eventService.JoinByCode(ctx.Request.Context(), accountId, req.EventCode)
	ResponseJSON(ctx, req, res, err)
}

func (c *eventController) QuizListByEvent(ctx *gin.Context) {
	slug := ctx.Param("slug")
	res, err := c.eventService.QuizListByEvent(ctx.Request.Context(), slug)
	ResponseJSON(ctx, gin.H{"slug": slug}, res, err)
}
