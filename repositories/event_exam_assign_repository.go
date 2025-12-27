package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventExamAssignRepository interface {
	Create(ctx context.Context, m entity.EventExamAssign) error
	ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.EventExamAssign, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Check(ctx context.Context, eventId uuid.UUID, examId uuid.UUID) error
}

type eventExamAssignRepository struct{ db *gorm.DB }

func NewEventExamAssignRepository(db *gorm.DB) EventExamAssignRepository {
	return &eventExamAssignRepository{db}
}

func (r *eventExamAssignRepository) Check(ctx context.Context, eventId uuid.UUID, examId uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("event_id = ?", eventId).
		Where("exam_id = ?", examId).
		First(&entity.EventExamAssign{}).Error
}
func (r *eventExamAssignRepository) Create(ctx context.Context, m entity.EventExamAssign) error {
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *eventExamAssignRepository) ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.EventExamAssign, error) {
	var items []entity.EventExamAssign
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventId).
		Preload("Exam").
		Find(&items).Error
	return items, err
}

func (r *eventExamAssignRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.EventExamAssign{}).Error
}
