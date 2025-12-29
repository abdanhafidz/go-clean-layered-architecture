package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResultRepository interface {
	Create(ctx context.Context, r *entity.Result) error
	GetById(ctx context.Context, id uuid.UUID) (entity.Result, error)
	Update(ctx context.Context, r *entity.Result) error
	ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Result, error)
	GetByAttemptId(ctx context.Context, attemptId uuid.UUID) (entity.Result, error)
	ListByEventAndAccount(ctx context.Context, eventId uuid.UUID, accountId uuid.UUID) ([]entity.Result, error)
}

type resultRepository struct{ db *gorm.DB }

func NewResultRepository(db *gorm.DB) ResultRepository {
	return &resultRepository{db}
}

func (r *resultRepository) Create(ctx context.Context, rs *entity.Result) error {
	return r.db.WithContext(ctx).Create(rs).Error
}
func (r *resultRepository) GetByAttemptId(ctx context.Context, attemptId uuid.UUID) (entity.Result, error) {
	var rs entity.Result
	err := r.db.WithContext(ctx).
		Preload("EventExamAttempt").
		First(&rs, "attempt_id = ?", attemptId).Error
	return rs, err
}
func (r *resultRepository) GetById(ctx context.Context, id uuid.UUID) (entity.Result, error) {
	var rs entity.Result
	err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Event").
		Preload("ProblemSet").
		Preload("EventExamAttempt").
		First(&rs, "id = ?", id).Error
	return rs, err
}

func (r *resultRepository) Update(ctx context.Context, rs *entity.Result) error {
	return r.db.WithContext(ctx).
		Model(&entity.Result{}).
		Where("id = ?", rs.Id).
		Updates(rs).Error
}

func (r *resultRepository) ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Result, error) {
	var list []entity.Result
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventId).
		Preload("EventExamAttempt").
		Preload("EventExamAttempt.Account").
		Preload("EventExamAttempt.Exam").
		Find(&list).Error
	return list, err
}

func (r *resultRepository) ListByEventAndAccount(ctx context.Context, eventId uuid.UUID, accountId uuid.UUID) ([]entity.Result, error) {
	var list []entity.Result
	err := r.db.WithContext(ctx).
		Table("result").
		Joins("JOIN exam_event_attempt ON exam_event_attempt.id = result.attempt_id").
		Where("exam_event_attempt.event_id = ?", eventId).
		Where("exam_event_attempt.account_id = ?", accountId).
		Preload("EventExamAttempt").
		Find(&list).Error
	return list, err
}
