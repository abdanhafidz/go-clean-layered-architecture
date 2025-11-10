package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemSetRepository interface {
	Create(ctx context.Context, ps entity.ProblemSet) error
	Get(ctx context.Context, id uuid.UUID) (entity.ProblemSet, error)
	Update(ctx context.Context, ps entity.ProblemSet) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]entity.ProblemSet, error)
}

type problemSetRepository struct {
	db *gorm.DB
}

func NewProblemSetRepository(db *gorm.DB) ProblemSetRepository {
	return &problemSetRepository{db: db}
}

func (r *problemSetRepository) Create(ctx context.Context, ps entity.ProblemSet) error {
	return r.db.WithContext(ctx).Create(&ps).Error
}

func (r *problemSetRepository) Get(ctx context.Context, id uuid.UUID) (entity.ProblemSet, error) {
	var ps entity.ProblemSet
	err := r.db.WithContext(ctx).
		First(&ps, "id_problem_set = ?", id).Error
	return ps, err
}

func (r *problemSetRepository) List(ctx context.Context) ([]entity.ProblemSet, error) {
	var list []entity.ProblemSet
	err := r.db.WithContext(ctx).
		Order("title").
		Find(&list).Error
	return list, err
}

func (r *problemSetRepository) Update(ctx context.Context, ps entity.ProblemSet) error {
	return r.db.WithContext(ctx).
		Model(&entity.ProblemSet{}).
		Where("id_problem_set = ?", ps.Id).
		Updates(ps).Error
}

func (r *problemSetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("id_problem_set = ?", id).
		Delete(&entity.ProblemSet{}).Error
}
