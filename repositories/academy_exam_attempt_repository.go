package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademyExamAttemptRepository interface {
	Create(ctx context.Context, a *entity.AcademyExamAttempt) error
	GetById(ctx context.Context, attemptId uuid.UUID) (entity.AcademyExamAttempt, error)
	GetByAcademyExam(ctx context.Context, academyId uuid.UUID, examId uuid.UUID, accountId uuid.UUID) (entity.AcademyExamAttempt, error)
	Update(ctx context.Context, a *entity.AcademyExamAttempt) error
}

type academyExamAttemptRepository struct{ db *gorm.DB }

func NewAcademyExamAttemptRepository(db *gorm.DB) AcademyExamAttemptRepository {
	return &academyExamAttemptRepository{db}
}

func (r *academyExamAttemptRepository) Create(ctx context.Context, a *entity.AcademyExamAttempt) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *academyExamAttemptRepository) GetById(ctx context.Context, attemptId uuid.UUID) (entity.AcademyExamAttempt, error) {
	var a entity.AcademyExamAttempt
	err := r.db.WithContext(ctx).
		Preload("Answers").
		First(&a, "id = ?", attemptId).Error
	return a, err
}

func (r *academyExamAttemptRepository) GetByAcademyExam(ctx context.Context, academyId uuid.UUID, examId uuid.UUID, accountId uuid.UUID) (entity.AcademyExamAttempt, error) {
	var attempt entity.AcademyExamAttempt
	err := r.db.WithContext(ctx).
		Preload("Answers").
		Where("academy_id = ?", academyId).
		Where("exam_id = ?", examId).
		Where("account_id = ?", accountId).
		First(&attempt).Error
	return attempt, err
}

func (r *academyExamAttemptRepository) Update(ctx context.Context, a *entity.AcademyExamAttempt) error {
	return r.db.WithContext(ctx).
		Model(&entity.AcademyExamAttempt{}).
		Where("id = ?", a.Id).
		Updates(a).Error
}
