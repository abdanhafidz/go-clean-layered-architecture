package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamEventAnswerRepository interface {
	Create(ctx context.Context, ans entity.ExamEventAnswer) error
	Update(ctx context.Context, ans entity.ExamEventAnswer) error
	GetByAttemptAndQuestion(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID) (entity.ExamEventAnswer, error)
	ListByAttempt(ctx context.Context, attemptId uuid.UUID) ([]entity.ExamEventAnswer, error)
	DeleteByAttempt(ctx context.Context, attemptId uuid.UUID) error

	// decomposed result (answer + question)
}

type examEventAnswerRepository struct {
	db *gorm.DB
}

func NewExamEventAnswerRepository(db *gorm.DB) ExamEventAnswerRepository {
	return &examEventAnswerRepository{db: db}
}

func (r *examEventAnswerRepository) Create(ctx context.Context, ans entity.ExamEventAnswer) error {
	return r.db.WithContext(ctx).Create(&ans).Error
}

func (r *examEventAnswerRepository) Update(ctx context.Context, ans entity.ExamEventAnswer) error {
	return r.db.WithContext(ctx).
		Where("id_attempt = ? AND id_question = ?", ans.AttemptId, ans.QuestionId).
		Updates(ans).Error
}

func (r *examEventAnswerRepository) GetByAttemptAndQuestion(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID) (entity.ExamEventAnswer, error) {
	var ans entity.ExamEventAnswer
	err := r.db.WithContext(ctx).
		Where("id_attempt = ? AND id_question = ?", attemptId, questionId).
		First(&ans).Error
	return ans, err
}

func (r *examEventAnswerRepository) ListByAttempt(ctx context.Context, attemptId uuid.UUID) ([]entity.ExamEventAnswer, error) {
	var answers []entity.ExamEventAnswer
	err := r.db.WithContext(ctx).
		Where("id_attempt = ?", attemptId).
		Find(&answers).Error
	return answers, err
}

func (r *examEventAnswerRepository) DeleteByAttempt(ctx context.Context, attemptId uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id_attempt = ?", attemptId).
		Delete(&entity.ExamEventAnswer{}).Error
}
