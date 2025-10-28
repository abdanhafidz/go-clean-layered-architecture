package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProblemSetAssignRepository interface {
	CreateProblemSetAssign(ctx context.Context, psa entity.ProblemSetAssign) (entity.ProblemSetAssign, error)
	GetProblemSetAssignByID(ctx context.Context, id uuid.UUID) (entity.ProblemSetAssign, error)
	GetProblemSetAssignByEventID(ctx context.Context, eventID uuid.UUID) ([]entity.ProblemSetAssign, error)
	GetAllProblemSetAssign(ctx context.Context) ([]entity.ProblemSetAssign, error)
	UpdateProblemSetAssign(ctx context.Context, psa entity.ProblemSetAssign) (entity.ProblemSetAssign, error)
	DeleteProblemSetAssign(ctx context.Context, id uuid.UUID) error
}

type problemSetAssignRepository struct {
	db *gorm.DB
}

func NewProblemSetAssignRepository(db *gorm.DB) ProblemSetAssignRepository {
	return &problemSetAssignRepository{db: db}
}

func (r *problemSetAssignRepository) CreateProblemSetAssign(ctx context.Context, psa entity.ProblemSetAssign) (entity.ProblemSetAssign, error) {
	if err := r.db.WithContext(ctx).Create(&psa).Error; err != nil {
		return entity.ProblemSetAssign{}, err
	}
	return psa, nil
}

func (r *problemSetAssignRepository) GetProblemSetAssignByID(ctx context.Context, id uuid.UUID) (entity.ProblemSetAssign, error) {
	var psa entity.ProblemSetAssign
	if err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("ProblemSet").
		First(&psa, "id = ?", id).Error; err != nil {
		return entity.ProblemSetAssign{}, err
	}
	return psa, nil
}

func (r *problemSetAssignRepository) GetProblemSetAssignByEventID(ctx context.Context, eventID uuid.UUID) ([]entity.ProblemSetAssign, error) {
	var psas []entity.ProblemSetAssign
	if err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("ProblemSet").
		Where("event_id = ?", eventID).
		Find(&psas).Error; err != nil {
		return nil, err
	}
	if len(psas) == 0 {
		// Optional logging or custom error
	}
	return psas, nil
}

func (r *problemSetAssignRepository) GetAllProblemSetAssign(ctx context.Context) ([]entity.ProblemSetAssign, error) {
	var psas []entity.ProblemSetAssign
	if err := r.db.WithContext(ctx).
		Preload("Event").
		Preload("ProblemSet").
		Find(&psas).Error; err != nil {
		return nil, err
	}
	return psas, nil
}

func (r *problemSetAssignRepository) UpdateProblemSetAssign(ctx context.Context, psa entity.ProblemSetAssign) (entity.ProblemSetAssign, error) {
	if err := r.db.WithContext(ctx).Save(&psa).Error; err != nil {
		return entity.ProblemSetAssign{}, err
	}
	return psa, nil
}

func (r *problemSetAssignRepository) DeleteProblemSetAssign(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.ProblemSetAssign{}, "id = ?", id).Error
}
