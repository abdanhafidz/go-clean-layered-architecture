package repositories

import (
    "context"

    entity "abdanhafidz.com/go-boilerplate/models/entity"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type AcademyResultRepository interface {
    Create(ctx context.Context, r *entity.AcademyExamResult) error
    GetById(ctx context.Context, id uuid.UUID) (entity.AcademyExamResult, error)
    Update(ctx context.Context, r *entity.AcademyExamResult) error
}

type academyResultRepository struct { db *gorm.DB }

func NewAcademyResultRepository(db *gorm.DB) AcademyResultRepository {
    return &academyResultRepository{ db }
}

func (r *academyResultRepository) Create(ctx context.Context, rec *entity.AcademyExamResult) error {
    return r.db.WithContext(ctx).Create(rec).Error
}

func (r *academyResultRepository) GetById(ctx context.Context, id uuid.UUID) (entity.AcademyExamResult, error) {
    var rec entity.AcademyExamResult
    err := r.db.WithContext(ctx).
        First(&rec, "id = ?", id).Error
    return rec, err
}

func (r *academyResultRepository) Update(ctx context.Context, rec *entity.AcademyExamResult) error {
    return r.db.WithContext(ctx).
        Model(&entity.AcademyExamResult{}).
        Where("id = ?", rec.Id).
        Updates(rec).Error
}