package controllers

import (
	"strconv"

	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/repositories"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type EventController interface {
	List(ctx *gin.Context)
	DetailBySlug(ctx *gin.Context)
	Join(ctx *gin.Context)
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
	accountId := ParseAccountId(ctx)
	list, total, err := c.eventService.List(ctx.Request.Context(), accountId, p)
	meta := gin.H{
		"total_records": total,
		"page_size":     limit,
		"current_page":  (offset / limit) + 1,
	}
	ResponseJSON(ctx, meta, list, err)
}

func (c *eventController) DetailBySlug(ctx *gin.Context) {
	slug := ctx.Param("event_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.eventService.DetailBySlug(ctx.Request.Context(), slug, accountId)
	ResponseJSON(ctx, gin.H{"event_slug": slug, "id_account": accountId}, res, err)
}

func (c *eventController) Join(ctx *gin.Context) {
	req := RequestJSON[dto.JoinEventRequest](ctx)
	accountId := ParseAccountId(ctx)
	res, err := c.eventService.JoinByCode(ctx.Request.Context(), accountId, req.EventCode)
	ResponseJSON(ctx, req, res, err)
}
