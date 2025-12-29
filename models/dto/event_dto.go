package dto

import (
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
)

type EventDetailResponse struct {
	Data           *entity.Events
	RegisterStatus int                            `json:"register_status" `
	EventPayment   entity.EventPaymentTransaction `json:"event_payment,omitempty"`
}

type JoinEventRequest struct {
	EventCode string `json:"event_code" binding:"required"`
}

type EventStatus struct {
	IsHasNotStarted bool
	IsOnGoing       bool
	IsFinished      bool
}

type CreateEventRequest struct {
	Title      string `json:"title" binding:"required"`
	StartEvent string `json:"start_event" binding:"required"`
	EndEvent   string `json:"end_event" binding:"required"`
	Overview   string `json:"overview" binding:"required"`
	ImgBanner  string `json:"img_banner" binding:"required"`
	EventCode  string `json:"event_code" binding:"required"`
	IsPublic   bool   `json:"is_public"`
}

type UpdateEventRequest struct {
	Title      string `json:"title"`
	StartEvent string `json:"start_event"`
	EndEvent   string `json:"end_event"`
	Overview   string `json:"overview"`
	ImgBanner  string `json:"img_banner"`
	IsPublic   *bool  `json:"is_public"`
}

type EventExamProctoringLogsRequest struct {
	ViolationScore    uint   `json:"violation_score,omitempty" form:"violation_score"`
	ViolationCategory string `json:"violation_category,omitempty" form:"violation_category"`
	Attachement       string `json:"attachement,omitempty" form:"attachement"`
}

type EventExamListResponse struct {
	entity.Exam
	Score float32 `json:"score,omitempty"`
}

type ExamScore struct {
	ExamId    uuid.UUID `json:"exam_id"`
	ExamTitle string    `json:"exam_title"`
	Score     float32   `json:"score"`
}

type ScoreboardItem struct {
	Username          string      `json:"username"`
	FullName          string      `json:"full_name,omitempty"`
	TotalExamDuration string      `json:"total_exam_duration"` // Formatted string or float? "username, total exam duration". Let's use string for display or seconds. Requested "column is ... total exam duration".
	TotalDurationInt  int64       `json:"-"`                   // For sorting
	Scores            []ExamScore `json:"scores"`
	TotalScore        float32     `json:"total_score"`
	AverageScore      float32     `json:"average_score"`
}
