package dto

import (
	"time"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
)

type UserExamStatus struct {
	IsNotAttempt bool
	IsOnAttempt  bool
	IsSubmitted  bool
	IsTimeOut    bool
}

type AnswerWithQuestion struct {
	Answer   entity.EventExamAnswer `json:"answer"`
	Question entity.Questions       `json:"question"`
}

type AnswerEventExamRequest struct {
	QuestionId uuid.UUID `json:"question_id" binding:"required"`
	Answer     []string  `json:"answer"`
}

type CreateExamRequest struct {
	Slug               string        `json:"slug,omitempty"`
	Title              string        `json:"title,omitempty"`
	Description        string        `json:"description,omitempty"`
	Duration           time.Duration `json:"duration,omitempty"`
	Randomize          uint          `json:"randomize,omitempty"`
	AllowRetake        bool          `json:"allow_retake,omitempty"`
	AllowReview        bool          `json:"allow_review,omitempty"`
	EnableTimer        bool          `json:"enable_timer,omitempty"`
	EnableWebCam       bool          `json:"enable_webcam,omitempty"`
	EnableVAD          bool          `json:"enable_vad,omitempty"`
	EnableTabBlock     bool          `json:"enable_tab_block,omitempty"`
	RequiredFullScreen bool          `json:"enable_full_screen,omitempty"`
	EnableEyeTracking  bool          `json:"enable_eye_tracking,omitempty"`
	DisableCopyPaste   bool          `json:"disable_copy_paste,omitempty"`
	EnableExamBrowser  bool          `json:"enable_exam_browser,omitempty"`
}
