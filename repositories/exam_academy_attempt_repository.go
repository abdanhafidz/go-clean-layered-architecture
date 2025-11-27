package repositories

import (
    "context"

    entity "abdanhafidz.com/go-boilerplate/models/entity"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ExamAcademyAttemptRepository interface {
    Create(ctx context.Context, a *entity.ExamAcademyAttempt) error
    GetById(ctx context.Context, attemptId uuid.UUID) (entity.ExamAcademyAttempt, error)
    GetByExamAcademy(ctx context.Context, academyId uuid.UUID, examId uuid.UUID, accountId uuid.UUID) (entity.ExamAcademyAttempt, error)
    Update(ctx context.Context, a *entity.ExamAcademyAttempt) error
}

type examAcademyAttemptRepository struct{ db *gorm.DB }

func NewExamAcademyAttemptRepository(db *gorm.DB) ExamAcademyAttemptRepository {
    return &examAcademyAttemptRepository{db}
}

func (r *examAcademyAttemptRepository) Create(ctx context.Context, a *entity.ExamAcademyAttempt) error {
    return r.db.WithContext(ctx).Create(a).Error
}

func (r *examAcademyAttemptRepository) GetById(ctx context.Context, attemptId uuid.UUID) (entity.ExamAcademyAttempt, error) {
    var a entity.ExamAcademyAttempt
    err := r.db.WithContext(ctx).
        Preload("Answers").
        First(&a, "id = ?", attemptId).Error
    return a, err
}

func (r *examAcademyAttemptRepository) GetByExamAcademy(ctx context.Context, academyId uuid.UUID, examId uuid.UUID, accountId uuid.UUID) (entity.ExamAcademyAttempt, error) {
    var attempt entity.ExamAcademyAttempt
    err := r.db.WithContext(ctx).
        Preload("Answers").
        Where("academy_id = ?", academyId).
        Where("exam_id = ?", examId).
        Where("account_id = ?", accountId).
        First(&attempt).Error
    return attempt, err
}

func (r *examAcademyAttemptRepository) Update(ctx context.Context, a *entity.ExamAcademyAttempt) error {
    return r.db.WithContext(ctx).
        Model(&entity.ExamAcademyAttempt{}).
        Where("id = ?", a.Id).
        Updates(a).Error
}