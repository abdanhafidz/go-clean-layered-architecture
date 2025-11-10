package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResultRepository interface {
	Create(ctx context.Context, r entity.Result) error
	Get(ctx context.Context, id uuid.UUID) (entity.Result, error)
	Update(ctx context.Context, r entity.Result) error
	ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Result, error)
}

type resultRepository struct{ db *gorm.DB }

func NewResultRepository(db *gorm.DB) ResultRepository {
	return &resultRepository{db}
}

func (r *resultRepository) Create(ctx context.Context, rs entity.Result) error {
	return r.db.WithContext(ctx).Create(&rs).Error
}

func (r *resultRepository) Get(ctx context.Context, id uuid.UUID) (entity.Result, error) {
	var rs entity.Result
	err := r.db.WithContext(ctx).
		Preload("Account").
		Preload("Event").
		Preload("ProblemSet").
		Preload("ExamEventAttempt").
		First(&rs, "id_result = ?", id).Error
	return rs, err
}

func (r *resultRepository) Update(ctx context.Context, rs entity.Result) error {
	return r.db.WithContext(ctx).
		Model(&entity.Result{}).
		Where("id_result = ?", rs.Id).
		Updates(rs).Error
}

func (r *resultRepository) ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Result, error) {
	var list []entity.Result
	err := r.db.WithContext(ctx).
		Where("id_event = ?", eventId).
		Preload("Account").
		Find(&list).Error
	return list, err
}
