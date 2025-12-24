package controllers

import (
	"errors"
	"strconv"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
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

// Event List godoc
// @Summary      List Events
// @Description  Retrieve a paginated list of events with optional filters
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        limit           query     int     false  "Number of items per page"            default(10)
// @Param        page            query     int     false  "Page number"                         default(1)
// @Param        search          query     string  false  "Search term for event titles"
// @Param        sortBy          query     string  false  "Field to sort by (e.g., 'date', 'title')"
// @Param        order           query     string  false  "Sort order ('asc' or 'desc')"
// @Param        registerStatus  query     int     false  "Filter by registration status"
// @Param        status          query     string  false  "Filter by event status (upcoming, ongoing, ended)"
// @Success	  	 200  {object}  dto.SuccessResponse[[]entity.Events]
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/v1/events [get]
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

	var status *string

	if val := ctx.Query("status"); val != "" {
		if val == entity.EventStatusUpcoming || val == entity.EventStatusOngoing || val == entity.EventStatusEnded {
			status = &val
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
	p := entity.Pagination{Limit: limit, Offset: offset, Search: search, SortBy: sortBy, Order: order, RegisterStatus: registerStatus, Status: status}

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

// Event Detail By Slug godoc
// @Summary      Get Event Detail by Slug
// @Description  Retrieve detailed information about a specific event using its slug
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        event_slug  path      string  true  "Event Slug"
// @Success      200         {object}  dto.SuccessResponse[dto.EventDetailResponse]
// @Failure      400         {object}  dto.ErrorResponse
// @Router       /api/v1/events/{event_slug} [get]
func (c *eventController) DetailBySlug(ctx *gin.Context) {
	slug := ctx.Param("event_slug")
	accountId := ParseAccountId(ctx)
	res, err := c.eventService.DetailBySlug(ctx.Request.Context(), slug, accountId)
	ResponseJSON(ctx, gin.H{"event_slug": slug, "id_account": accountId}, res, err)
}

// Join Event godoc
// @Summary      Join Event
// @Description  Register the authenticated user for an event using an event code
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        request  body      dto.JoinEventRequest  true  "Join Event Request"
// @Success      200      {object}  dto.SuccessResponse[dto.EventDetailResponse]
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      402      {object}  dto.SuccessResponse[entity.EventPaymentTransaction]
// @Router       /api/v1/events/register-event [post]
func (c *eventController) Join(ctx *gin.Context) {
	req := RequestJSON[dto.JoinEventRequest](ctx)

	if err := utils.ValidateCode(req.EventCode); err != nil {
		ResponseJSON(ctx, req, gin.H{}, err)
		return
	}

	accountId := ParseAccountId(ctx)
	res, err := c.eventService.JoinByCode(ctx.Request.Context(), accountId, req.EventCode)

	if errors.Is(err, http_error.PAYMENT_REQUIRED) {
		ResponseJSON(ctx, res.EventPayment, res, err)
		return
	}

	ResponseJSON(ctx, req, res, err)
}

// Create Event godoc
// @Summary      Create Event
// @Description  Create a new event with the provided details
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateEventRequest  true  "Create Event Request"
// @Success      200      {object}  dto.SuccessResponse[dto.EventDetailResponse]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/admin/events [post]
func (c *eventController) CreateEvent(ctx *gin.Context) {
	req := RequestJSON[dto.CreateEventRequest](ctx)
	res, err := c.eventService.CreateEvent(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

// Update Event godoc
// @Summary      Update Event
// @Description  Update an existing event with the provided details
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Event ID"
// @Param        request  body      dto.UpdateEventRequest  true  "Update Event Request"
// @Success      200      {object}  dto.SuccessResponse[entity.Events]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/admin/events/{id} [put]
func (c *eventController) UpdateEvent(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	req := RequestJSON[dto.UpdateEventRequest](ctx)
	res, err := c.eventService.UpdateEvent(ctx.Request.Context(), id, req)
	ResponseJSON(ctx, req, res, err)
}

// Delete Event godoc
// @Summary      Delete Event
// @Description  Delete an existing event by its ID
// @Tags         Event
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Event ID"
// @Success      200      {object}  dto.SuccessResponse[map[string]bool]
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/v1/admin/events/{id} [delete]
func (c *eventController) DeleteEvent(ctx *gin.Context) {
	id := ParseUUID(ctx, "id")
	err := c.eventService.DeleteEvent(ctx.Request.Context(), id)
	ResponseJSON(ctx, gin.H{"id": id}, gin.H{"deleted": true}, err)
}
