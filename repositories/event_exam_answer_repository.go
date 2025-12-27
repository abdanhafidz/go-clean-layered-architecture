package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventExamAnswerRepository interface {
	Create(ctx context.Context, ans *entity.EventExamAnswer) error
	Update(ctx context.Context, ans *entity.EventExamAnswer) error
	GetByAttemptAndQuestion(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID) (entity.EventExamAnswer, error)
	ListByAttempt(ctx context.Context, attemptId uuid.UUID) ([]entity.EventExamAnswer, error)
	DeleteByAttempt(ctx context.Context, attemptId uuid.UUID) error

	// decomposed result (answer + question)
}

type eventExamAnswerRepository struct {
	db *gorm.DB
}

func NewEventExamAnswerRepository(db *gorm.DB) EventExamAnswerRepository {
	return &eventExamAnswerRepository{db: db}
}

func (r *eventExamAnswerRepository) Create(ctx context.Context, ans *entity.EventExamAnswer) error {
	return r.db.WithContext(ctx).Create(ans).Error
}

func (r *eventExamAnswerRepository) Update(ctx context.Context, ans *entity.EventExamAnswer) error {
	return r.db.WithContext(ctx).
		Where("attempt_id = ? AND question_id = ?", ans.AttemptId, ans.QuestionId).
		Updates(ans).Error
}

func (r *eventExamAnswerRepository) GetByAttemptAndQuestion(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID) (entity.EventExamAnswer, error) {
	var ans entity.EventExamAnswer
	err := r.db.WithContext(ctx).
		Where("attempt_id = ? AND question_id = ?", attemptId, questionId).
		First(&ans).Error
	return ans, err
}

func (r *eventExamAnswerRepository) ListByAttempt(ctx context.Context, attemptId uuid.UUID) ([]entity.EventExamAnswer, error) {
	var answers []entity.EventExamAnswer
	err := r.db.WithContext(ctx).
		Where("attempt_id = ?", attemptId).
		Find(&answers).Error
	return answers, err
}

func (r *eventExamAnswerRepository) DeleteByAttempt(ctx context.Context, attemptId uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("attempt_id = ?", attemptId).
		Delete(&entity.EventExamAnswer{}).Error
}
