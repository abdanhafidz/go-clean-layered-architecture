package controllers

import (
	"strconv"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventController interface {
	List(ctx *gin.Context)
	DetailBySlug(ctx *gin.Context)
	Join(ctx *gin.Context)
	CreateEvent(ctx *gin.Context)
	UpdateEvent(ctx *gin.Context)
	DeleteEvent(ctx *gin.Context)
}

type eventController struct {
	eventService services.EventService
}

func NewEventController(eventService services.EventService) EventController {
	return &eventController{eventService: eventService}
}

func (c *eventController) List(ctx *gin.Context) {
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
	if limit < 1 {
		limit = 10
	} else if limit > 50 {
		limit = 50
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit
	p := entity.Pagination{Limit: limit, Offset: offset, Search: search, SortBy: sortBy, Order: order, RegisterStatus: registerStatus}
	accountId := ParseAccountId(ctx)
	list, total, err := c.eventService.List(ctx.Request.Context(), accountId, p)
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
		list, total, err = c.eventService.List(ctx.Request.Context(), accountId, p)
	}
	meta := gin.H{
		"totalItems":  total,
		"totalPages":  totalPages,
		"currentPage": page,
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

func (c *eventController) CreateEvent(ctx *gin.Context) {
	req := RequestJSON[dto.CreateEventRequest](ctx)
	res, err := c.eventService.CreateEvent(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *eventController) UpdateEvent(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}
	req := RequestJSON[dto.UpdateEventRequest](ctx)
	res, err := c.eventService.UpdateEvent(ctx.Request.Context(), id, req)
	ResponseJSON(ctx, req, res, err)
}

func (c *eventController) DeleteEvent(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ResponseJSON[any, any](ctx, nil, nil, http_error.BAD_REQUEST_ERROR)
		return
	}
	err = c.eventService.DeleteEvent(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}
