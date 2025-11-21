package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionsRepository interface {
	Create(ctx context.Context, q entity.Questions) error
	Get(ctx context.Context, id uuid.UUID) (entity.Questions, error)
	Update(ctx context.Context, q entity.Questions) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByProblemSet(ctx context.Context, problemSetId uuid.UUID) ([]entity.Questions, error)
}

type questionsRepository struct{ db *gorm.DB }

func NewQuestionsRepository(db *gorm.DB) QuestionsRepository {
	return &questionsRepository{db}
}

func (r *questionsRepository) Create(ctx context.Context, q entity.Questions) error {
	return r.db.WithContext(ctx).Create(&q).Error
}

func (r *questionsRepository) Get(ctx context.Context, id uuid.UUID) (entity.Questions, error) {
	var q entity.Questions
	err := r.db.WithContext(ctx).First(&q, "id = ?", id).Error
	return q, err
}

func (r *questionsRepository) Update(ctx context.Context, q entity.Questions) error {
	return r.db.WithContext(ctx).
		Model(&entity.Questions{}).
		Where("id = ?", q.Id).
		Updates(q).Error
}

func (r *questionsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&entity.Questions{}).Error
}

func (r *questionsRepository) ListByProblemSet(ctx context.Context, problemSetId uuid.UUID) ([]entity.Questions, error) {
	var q []entity.Questions
	err := r.db.WithContext(ctx).
		Where("problem_set_id = ?", problemSetId).
		Order("id").
		Find(&q).Error
	return q, err
}
