package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamEventAssignRepository interface {
	Create(ctx context.Context, m entity.ExamEventAssign) error
	ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.ExamEventAssign, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type examEventAssignRepository struct{ db *gorm.DB }

func NewExamEventAssignRepository(db *gorm.DB) ExamEventAssignRepository {
	return &examEventAssignRepository{db}
}

func (r *examEventAssignRepository) Create(ctx context.Context, m entity.ExamEventAssign) error {
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *examEventAssignRepository) ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.ExamEventAssign, error) {
	var items []entity.ExamEventAssign
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventId).
		Preload("Exam").
		Find(&items).Error
	return items, err
}

func (r *examEventAssignRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.ExamEventAssign{}).Error
}
