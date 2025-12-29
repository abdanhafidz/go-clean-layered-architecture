package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademyExamAssignRepository interface {
	Create(ctx context.Context, m entity.AcademyExamAssign) error
	ListByAcademy(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyExamAssign, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Check(ctx context.Context, academyId uuid.UUID, examId uuid.UUID) error
}

type academyExamAssignRepository struct{ db *gorm.DB }

func NewAcademyExamAssignRepository(db *gorm.DB) AcademyExamAssignRepository {
	return &academyExamAssignRepository{db}
}

func (r *academyExamAssignRepository) Check(ctx context.Context, academyId uuid.UUID, examId uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("academy_id = ?", academyId).
		Where("exam_id = ?", examId).
		First(&entity.AcademyExamAssign{}).Error
}

func (r *academyExamAssignRepository) Create(ctx context.Context, m entity.AcademyExamAssign) error {
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *academyExamAssignRepository) ListByAcademy(ctx context.Context, academyId uuid.UUID) ([]entity.AcademyExamAssign, error) {
	var items []entity.AcademyExamAssign
	err := r.db.WithContext(ctx).
		Where("academy_id = ?", academyId).
		Preload("Exam").
		Find(&items).Error
	return items, err
}

func (r *academyExamAssignRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.AcademyExamAssign{}).Error
}
