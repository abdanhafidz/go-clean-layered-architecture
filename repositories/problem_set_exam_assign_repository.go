package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemSetExamAssignRepository interface {
	Create(ctx context.Context, m entity.ProblemSetExamAssign) error
	GetByExam(ctx context.Context, examId uuid.UUID) (entity.ProblemSetExamAssign, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type problemSetExamAssignRepository struct{ db *gorm.DB }

func NewProblemSetExamAssignRepository(db *gorm.DB) ProblemSetExamAssignRepository {
	return &problemSetExamAssignRepository{db}
}

func (r *problemSetExamAssignRepository) Create(ctx context.Context, m entity.ProblemSetExamAssign) error {
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *problemSetExamAssignRepository) GetByExam(ctx context.Context, examId uuid.UUID) (entity.ProblemSetExamAssign, error) {
	var items entity.ProblemSetExamAssign
	err := r.db.WithContext(ctx).
		Where("exam_id = ?", examId).
		Preload("ProblemSet").
		First(&items).Error
	return items, err
}

func (r *problemSetExamAssignRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.ProblemSetExamAssign{}).Error
}
