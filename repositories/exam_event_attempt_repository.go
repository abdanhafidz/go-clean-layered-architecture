package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamEventAttemptRepository interface {
	Create(ctx context.Context, a entity.ExamEventAttempt) error
	GetById(ctx context.Context, attemptId uuid.UUID) (entity.ExamEventAttempt, error)
	GetByExamEvent(ctx context.Context, eventId uuid.UUID, examId uuid.UUID, accountId uuid.UUID) (entity.ExamEventAttempt, error)
	Update(ctx context.Context, a entity.ExamEventAttempt) error
}

type examEventAttemptRepository struct{ db *gorm.DB }

func NewExamEventAttemptRepository(db *gorm.DB) ExamEventAttemptRepository {
	return &examEventAttemptRepository{db}
}

func (r *examEventAttemptRepository) Create(ctx context.Context, a entity.ExamEventAttempt) error {
	return r.db.WithContext(ctx).Create(&a).Error
}

func (r *examEventAttemptRepository) GetById(ctx context.Context, attemptId uuid.UUID) (entity.ExamEventAttempt, error) {
	var a entity.ExamEventAttempt
	err := r.db.WithContext(ctx).
		Preload("Questions").
		First(&a, "id_attempt = ?", attemptId).Error
	return a, err
}

func (r *examEventAttemptRepository) GetByExamEvent(ctx context.Context, eventId uuid.UUID, examId uuid.UUID, accountId uuid.UUID) (entity.ExamEventAttempt, error) {

	var attempt entity.ExamEventAttempt

	err := r.db.WithContext(ctx).
		Preload("Questions").
		Where("id_event = ?", eventId).
		Where("id_exam = ?", examId).
		Where("id_account = ?", accountId).
		First(&attempt).Error

	return attempt, err
}

func (r *examEventAttemptRepository) Update(ctx context.Context, a entity.ExamEventAttempt) error {
	return r.db.WithContext(ctx).
		Model(&entity.ExamEventAttempt{}).
		Where("id_attempt = ?", a.Id).
		Updates(a).Error
}
