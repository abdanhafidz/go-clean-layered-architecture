package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventExamProctoringRepository interface {
	Create(ctx context.Context, log *entity.EventExamProctoringLogs) error
	List(ctx context.Context, accountId uuid.UUID, examId uuid.UUID, eventId uuid.UUID) ([]entity.EventExamProctoringLogs, error)
	GetById(ctx context.Context, id uuid.UUID) (entity.EventExamProctoringLogs, error)
	Update(ctx context.Context, log *entity.EventExamProctoringLogs) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type eventExamProctoringRepository struct {
	db *gorm.DB
}

func NewEventExamProctoringRepository(db *gorm.DB) EventExamProctoringRepository {
	return &eventExamProctoringRepository{db}
}

func (r *eventExamProctoringRepository) Create(ctx context.Context, log *entity.EventExamProctoringLogs) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *eventExamProctoringRepository) List(ctx context.Context, accountId uuid.UUID, examId uuid.UUID, eventId uuid.UUID) ([]entity.EventExamProctoringLogs, error) {
	var logs []entity.EventExamProctoringLogs
	query := r.db.WithContext(ctx)

	if accountId != uuid.Nil {
		query = query.Where("account_id = ?", accountId)
	}
	if examId != uuid.Nil {
		query = query.Where("exam_id = ?", examId)
	}
	if eventId != uuid.Nil {
		query = query.Where("event_id = ?", eventId)
	}

	err := query.Find(&logs).Error
	return logs, err
}

func (r *eventExamProctoringRepository) GetById(ctx context.Context, id uuid.UUID) (entity.EventExamProctoringLogs, error) {
	var log entity.EventExamProctoringLogs
	err := r.db.WithContext(ctx).First(&log, "id = ?", id).Error
	return log, err
}

func (r *eventExamProctoringRepository) Update(ctx context.Context, log *entity.EventExamProctoringLogs) error {
	return r.db.WithContext(ctx).Save(log).Error
}

func (r *eventExamProctoringRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.EventExamProctoringLogs{}, "id = ?", id).Error
}
