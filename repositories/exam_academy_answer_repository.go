package repositories

import (
    "context"

    entity "abdanhafidz.com/go-boilerplate/models/entity"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ExamAcademyAnswerRepository interface {
    Create(ctx context.Context, ans *entity.ExamAcademyAnswer) error
    Update(ctx context.Context, ans *entity.ExamAcademyAnswer) error
    GetByAttemptAndQuestion(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID) (entity.ExamAcademyAnswer, error)
    ListByAttempt(ctx context.Context, attemptId uuid.UUID) ([]entity.ExamAcademyAnswer, error)
    DeleteByAttempt(ctx context.Context, attemptId uuid.UUID) error
}

type examAcademyAnswerRepository struct { db *gorm.DB }

func NewExamAcademyAnswerRepository(db *gorm.DB) ExamAcademyAnswerRepository {
    return &examAcademyAnswerRepository{ db: db }
}

func (r *examAcademyAnswerRepository) Create(ctx context.Context, ans *entity.ExamAcademyAnswer) error {
    return r.db.WithContext(ctx).Create(ans).Error
}

func (r *examAcademyAnswerRepository) Update(ctx context.Context, ans *entity.ExamAcademyAnswer) error {
    return r.db.WithContext(ctx).
        Where("attempt_id = ? AND question_id = ?", ans.AttemptId, ans.QuestionId).
        Updates(ans).Error
}

func (r *examAcademyAnswerRepository) GetByAttemptAndQuestion(ctx context.Context, attemptId uuid.UUID, questionId uuid.UUID) (entity.ExamAcademyAnswer, error) {
    var ans entity.ExamAcademyAnswer
    err := r.db.WithContext(ctx).
        Where("attempt_id = ? AND question_id = ?", attemptId, questionId).
        First(&ans).Error
    return ans, err
}

func (r *examAcademyAnswerRepository) ListByAttempt(ctx context.Context, attemptId uuid.UUID) ([]entity.ExamAcademyAnswer, error) {
    var answers []entity.ExamAcademyAnswer
    err := r.db.WithContext(ctx).
        Where("attempt_id = ?", attemptId).
        Find(&answers).Error
    return answers, err
}

func (r *examAcademyAnswerRepository) DeleteByAttempt(ctx context.Context, attemptId uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("attempt_id = ?", attemptId).
        Delete(&entity.ExamAcademyAnswer{}).Error
}