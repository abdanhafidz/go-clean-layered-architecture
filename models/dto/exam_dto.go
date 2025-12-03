package dto

import (
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
	Answer   entity.ExamEventAnswer `json:"answer"`
	Question entity.Questions       `json:"question"`
}

type AnswerExamEventRequest struct {
	QuestionId uuid.UUID `json:"question_id" binding:"required"`
	Answer     []string  `json:"answer"`
}
