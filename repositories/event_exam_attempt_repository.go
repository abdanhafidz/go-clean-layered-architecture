package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventExamAttemptRepository interface {
	Create(ctx context.Context, a *entity.EventExamAttempt) error
	GetById(ctx context.Context, attemptId uuid.UUID) (entity.EventExamAttempt, error)
	GetByEventExam(ctx context.Context, eventId uuid.UUID, examId uuid.UUID, accountId uuid.UUID) (entity.EventExamAttempt, error)
	Update(ctx context.Context, a *entity.EventExamAttempt) error
}

type eventExamAttemptRepository struct{ db *gorm.DB }

func NewEventExamAttemptRepository(db *gorm.DB) EventExamAttemptRepository {
	return &eventExamAttemptRepository{db}
}

func (r *eventExamAttemptRepository) Create(ctx context.Context, a *entity.EventExamAttempt) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *eventExamAttemptRepository) GetById(ctx context.Context, attemptId uuid.UUID) (entity.EventExamAttempt, error) {
	var a entity.EventExamAttempt
	err := r.db.WithContext(ctx).
		Preload("Answers").
		First(&a, "id = ?", attemptId).Error
	return a, err
}

func (r *eventExamAttemptRepository) GetByEventExam(ctx context.Context, eventId uuid.UUID, examId uuid.UUID, accountId uuid.UUID) (entity.EventExamAttempt, error) {

	var attempt entity.EventExamAttempt

	err := r.db.WithContext(ctx).
		Preload("Answers").
		Where("event_id = ?", eventId).
		Where("exam_id = ?", examId).
		Where("account_id = ?", accountId).
		First(&attempt).Error

	return attempt, err
}

func (r *eventExamAttemptRepository) Update(ctx context.Context, a *entity.EventExamAttempt) error {
	return r.db.WithContext(ctx).
		Model(&entity.EventExamAttempt{}).
		Where("id = ?", a.Id).
		Updates(a).Error
}
