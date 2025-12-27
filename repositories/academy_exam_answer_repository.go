package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademyExamAnswerRepository interface {
	Create(ctx context.Context, ans *entity.AcademyExamAnswer) error
	Update(ctx context.Context, ans *entity.AcademyExamAnswer) error
	GetByAttemptAndQuestion(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID) (entity.AcademyExamAnswer, error)
	ListByAttempt(ctx context.Context, attemptId uuid.UUID) ([]entity.AcademyExamAnswer, error)
	DeleteByAttempt(ctx context.Context, attemptId uuid.UUID) error
}

type academyExamAnswerRepository struct{ db *gorm.DB }

func NewAcademyExamAnswerRepository(db *gorm.DB) AcademyExamAnswerRepository {
	return &academyExamAnswerRepository{db: db}
}

func (r *academyExamAnswerRepository) Create(ctx context.Context, ans *entity.AcademyExamAnswer) error {
	return r.db.WithContext(ctx).Create(ans).Error
}

func (r *academyExamAnswerRepository) Update(ctx context.Context, ans *entity.AcademyExamAnswer) error {
	return r.db.WithContext(ctx).
		Where("attempt_id = ? AND question_id = ?", ans.AttemptId, ans.QuestionId).
		Updates(ans).Error
}

func (r *academyExamAnswerRepository) GetByAttemptAndQuestion(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID) (entity.AcademyExamAnswer, error) {
	var ans entity.AcademyExamAnswer
	err := r.db.WithContext(ctx).
		Where("attempt_id = ? AND question_id = ?", attemptId, questionId).
		First(&ans).Error
	return ans, err
}

func (r *academyExamAnswerRepository) ListByAttempt(ctx context.Context, attemptId uuid.UUID) ([]entity.AcademyExamAnswer, error) {
	var answers []entity.AcademyExamAnswer
	err := r.db.WithContext(ctx).
		Where("attempt_id = ?", attemptId).
		Find(&answers).Error
	return answers, err
}

func (r *academyExamAnswerRepository) DeleteByAttempt(ctx context.Context, attemptId uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("attempt_id = ?", attemptId).
		Delete(&entity.AcademyExamAnswer{}).Error
}
