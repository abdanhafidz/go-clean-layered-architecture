package repositories

import (
    "context"

    entity "abdanhafidz.com/go-boilerplate/models/entity"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ExamAcademyAssignRepository interface {
    Create(ctx context.Context, m entity.ExamAcademyAssign) error
    ListByAcademy(ctx context.Context, academyId uuid.UUID) ([]entity.ExamAcademyAssign, error)
    Delete(ctx context.Context, id uuid.UUID) error
    Check(ctx context.Context, academyId uuid.UUID, examId uuid.UUID) error
}

type examAcademyAssignRepository struct{ db *gorm.DB }

func NewExamAcademyAssignRepository(db *gorm.DB) ExamAcademyAssignRepository {
    return &examAcademyAssignRepository{db}
}

func (r *examAcademyAssignRepository) Check(ctx context.Context, academyId uuid.UUID, examId uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("academy_id = ?", academyId).
        Where("exam_id = ?", examId).
        First(&entity.ExamAcademyAssign{}).Error
}

func (r *examAcademyAssignRepository) Create(ctx context.Context, m entity.ExamAcademyAssign) error {
    return r.db.WithContext(ctx).Create(&m).Error
}

func (r *examAcademyAssignRepository) ListByAcademy(ctx context.Context, academyId uuid.UUID) ([]entity.ExamAcademyAssign, error) {
    var items []entity.ExamAcademyAssign
    err := r.db.WithContext(ctx).
        Where("academy_id = ?", academyId).
        Preload("Exam").
        Find(&items).Error
    return items, err
}

func (r *examAcademyAssignRepository) Delete(ctx context.Context, id uuid.UUID) error {
    return r.db.WithContext(ctx).
        Where("id = ?", id).
        Delete(&entity.ExamAcademyAssign{}).Error
}