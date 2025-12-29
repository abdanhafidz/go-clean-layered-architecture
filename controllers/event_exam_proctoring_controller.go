package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventExamProctoringController interface {
	CreateLog(c *gin.Context)
	ListLogs(c *gin.Context)
	GetLogById(c *gin.Context)
	UpdateLog(c *gin.Context)
	DeleteLog(c *gin.Context)
}

type eventExamProctoringController struct {
	service services.EventExamProctoringService
}

func NewEventExamProctoringController(service services.EventExamProctoringService) EventExamProctoringController {
	return &eventExamProctoringController{service: service}
}

// CreateLog godoc
// @Summary Create Proctoring Log
// @Description Create a new proctoring log entry with optional file attachment
// @Tags Event Exam Proctoring
// @Accept multipart/form-data
// @Produce json
// @Param id_event formData string true "Event ID"
// @Param id_exam formData string true "Exam ID"
// @Param id_account formData string true "Account ID"
// @Param violation_score formData int true "Violation Score"
// @Param violation_category formData string true "Violation Category"
// @Param file formData file false "Attachment File"
// @Success 200 {object} dto.SuccessResponse[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/events/{event_slug}/exam/{exam_slug}/proctoring [post]
func (ctrl *eventExamProctoringController) CreateLog(c *gin.Context) {
	req := RequestForm[dto.EventExamProctoringLogsRequest](c)

	eventSlug := c.Param("event_slug")
	examSlug := c.Param("exam_slug")
	accountId := ParseAccountId(c)

	file, _ := c.FormFile("file")

	err := ctrl.service.CreateLog(c.Request.Context(), eventSlug, examSlug, accountId, req, file)
	ResponseJSON(c, req, "OK", err)
}

// ListLogs godoc
// @Summary List Proctoring Logs
// @Description List proctoring logs by account, exam, or event
// @Tags Event Exam Proctoring
// @Accept json
// @Produce json
// @Param event_slug path string true "Event Slug"
// @Param exam_slug path string true "Exam Slug"
// @Param account_id query string false "Account ID"
// @Param exam_id query string false "Exam ID"
// @Param event_id query string false "Event ID"
// @Success 200 {object} dto.SuccessResponse[[]models.EventExamProctoringLogs]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/admin/proctoring/events/{event_slug}/exam/{exam_slug}/proctoring [get]
func (ctrl *eventExamProctoringController) ListLogs(c *gin.Context) {
	accountId := ParseUUIDFromQuery(c, "account_id")
	examId := ParseUUIDFromQuery(c, "exam_id")
	eventId := ParseUUIDFromQuery(c, "event_id")

	logs, err := ctrl.service.ListLogs(c.Request.Context(), accountId, examId, eventId)
	ResponseJSON(c, gin.H{}, logs, err)
}

// GetLogById godoc
// @Summary Get Proctoring Log By ID
// @Description Get details of a specific proctoring log
// @Tags Event Exam Proctoring
// @Accept json
// @Produce json
// @Param event_slug path string true "Event Slug"
// @Param exam_slug path string true "Exam Slug"
// @Param log_id path string true "Log ID"
// @Success 200 {object} dto.SuccessResponse[models.EventExamProctoringLogs]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/admin/proctoring/events/{event_slug}/exam/{exam_slug}/proctoring/{log_id} [get]
func (ctrl *eventExamProctoringController) GetLogById(c *gin.Context) {
	id := ParseUUID(c, "log_id")
	log, err := ctrl.service.GetLogById(c.Request.Context(), id)
	ResponseJSON(c, gin.H{}, log, err)
}

// UpdateLog godoc
// @Summary Update Proctoring Log
// @Description Update an existing proctoring log
// @Tags Event Exam Proctoring
// @Accept multipart/form-data
// @Produce json
// @Param event_slug path string true "Event Slug"
// @Param exam_slug path string true "Exam Slug"
// @Param log_id path string true "Log ID"
// @Param violation_score formData int false "Violation Score"
// @Param violation_category formData string false "Violation Category"
// @Param file formData file false "Attachment File"
// @Param id_account formData string true "Account ID (required for upload context)"
// @Success 200 {object} dto.SuccessResponse[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/admin/proctoring/events/{event_slug}/exam/{exam_slug}/proctoring/{log_id} [put]
func (ctrl *eventExamProctoringController) UpdateLog(c *gin.Context) {
	id := ParseUUID(c, "log_id")

	req := RequestForm[dto.EventExamProctoringLogsRequest](c)

	file, _ := c.FormFile("file")

	err := ctrl.service.UpdateLog(c.Request.Context(), id, req, file)
	ResponseJSON(c, req, "OK", err)
}

// DeleteLog godoc
// @Summary Delete Proctoring Log
// @Description Delete a proctoring log entry
// @Tags Event Exam Proctoring
// @Accept json
// @Produce json
// @Param event_slug path string true "Event Slug"
// @Param exam_slug path string true "Exam Slug"
// @Param log_id path string true "Log ID"
// @Success 200 {object} dto.SuccessResponse[string]
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/admin/proctoring/events/{event_slug}/exam/{exam_slug}/proctoring/{log_id} [delete]
func (ctrl *eventExamProctoringController) DeleteLog(c *gin.Context) {
	id := ParseUUID(c, "log_id")
	err := ctrl.service.DeleteLog(c.Request.Context(), id)
	ResponseJSON(c, gin.H{"id": id}, "OK", err)
}

// Helper to parse UUID from query param since ParseUUID didn't seem to do it
func ParseUUIDFromQuery(c *gin.Context, key string) uuid.UUID {
	val := c.Query(key)
	if val == "" {
		return uuid.Nil
	}
	id, _ := uuid.Parse(val)
	return id
}
