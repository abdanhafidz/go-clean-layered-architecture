package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemSetRepository interface {
	CreateProblemSet(ctx context.Context, problemSet entity.ProblemSet) (entity.ProblemSet, error)
	GetProblemSetByID(ctx context.Context, id uuid.UUID) (entity.ProblemSet, error)
	GetAllProblemSets(ctx context.Context) ([]entity.ProblemSet, error)
	UpdateProblemSet(ctx context.Context, problemSet entity.ProblemSet) (entity.ProblemSet, error)
	DeleteProblemSet(ctx context.Context, id uuid.UUID) error
	GetProblemSetsWithQuestions(ctx context.Context) ([]entity.ProblemSet, error)
	GetProblemSetWithQuestionsByID(ctx context.Context, id uuid.UUID) (entity.ProblemSet, error)
}

type problemSetRepository struct {
	db *gorm.DB
}

func NewProblemSetRepository(db *gorm.DB) ProblemSetRepository {
	return &problemSetRepository{db: db}
}

func (r *problemSetRepository) CreateProblemSet(ctx context.Context, problemSet entity.ProblemSet) (entity.ProblemSet, error) {
	if err := r.db.WithContext(ctx).Create(&problemSet).Error; err != nil {
		return entity.ProblemSet{}, err
	}
	return problemSet, nil
}

func (r *problemSetRepository) GetProblemSetByID(ctx context.Context, id uuid.UUID) (entity.ProblemSet, error) {
	var problemSet entity.ProblemSet
	if err := r.db.WithContext(ctx).First(&problemSet, "id = ?", id).Error; err != nil {
		return entity.ProblemSet{}, err
	}
	return problemSet, nil
}

func (r *problemSetRepository) GetAllProblemSets(ctx context.Context) ([]entity.ProblemSet, error) {
	var problemSets []entity.ProblemSet
	if err := r.db.WithContext(ctx).Find(&problemSets).Error; err != nil {
		return nil, err
	}
	return problemSets, nil
}

func (r *problemSetRepository) UpdateProblemSet(ctx context.Context, problemSet entity.ProblemSet) (entity.ProblemSet, error) {
	if err := r.db.WithContext(ctx).Save(&problemSet).Error; err != nil {
		return entity.ProblemSet{}, err
	}
	return problemSet, nil
}

func (r *problemSetRepository) DeleteProblemSet(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.ProblemSet{}, "id = ?", id).Error
}

// Get all ProblemSets with preloaded questions (JOIN)
func (r *problemSetRepository) GetProblemSetsWithQuestions(ctx context.Context) ([]entity.ProblemSet, error) {
	var problemSets []entity.ProblemSet
	if err := r.db.WithContext(ctx).
		Preload("Questions").
		Find(&problemSets).Error; err != nil {
		return nil, err
	}
	return problemSets, nil
}

// Get single ProblemSet (by ID) with Questions
func (r *problemSetRepository) GetProblemSetWithQuestionsByID(ctx context.Context, id uuid.UUID) (entity.ProblemSet, error) {
	var problemSet entity.ProblemSet
	if err := r.db.WithContext(ctx).
		Preload("Questions").
		First(&problemSet, "id = ?", id).Error; err != nil {
		return entity.ProblemSet{}, err
	}
	return problemSet, nil
}
