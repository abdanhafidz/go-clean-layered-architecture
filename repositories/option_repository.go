package repositories

import (
	"context"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"gorm.io/gorm"
)

type OptionRepository interface {
	CreateOptionCategory(ctx context.Context, cat entity.OptionCategory) (entity.OptionCategory, error)
	CreateOptionValue(ctx context.Context, val entity.OptionValues) (entity.OptionValues, error)
	GetCategoryBySlug(ctx context.Context, slug string) (entity.OptionCategory, error)
	ListValuesByCategoryId(ctx context.Context, categoryId uint) ([]entity.OptionValues, error)
}

type optionRepository struct {
	db *gorm.DB
}

func NewOptionRepository(db *gorm.DB) OptionRepository {
	return &optionRepository{db: db}
}

func (r *optionRepository) CreateOptionCategory(ctx context.Context, cat entity.OptionCategory) (entity.OptionCategory, error) {
	if err := r.db.WithContext(ctx).Create(&cat).Error; err != nil {
		return entity.OptionCategory{}, err
	}
	return cat, nil
}

func (r *optionRepository) CreateOptionValue(ctx context.Context, val entity.OptionValues) (entity.OptionValues, error) {
	if err := r.db.WithContext(ctx).Create(&val).Error; err != nil {
		return entity.OptionValues{}, err
	}
	return val, nil
}

func (r *optionRepository) GetCategoryBySlug(ctx context.Context, slug string) (entity.OptionCategory, error) {
	var cat entity.OptionCategory
	if err := r.db.WithContext(ctx).First(&cat, "option_slug = ?", slug).Error; err != nil {
		return entity.OptionCategory{}, err
	}
	return cat, nil
}

func (r *optionRepository) ListValuesByCategoryId(ctx context.Context, categoryId uint) ([]entity.OptionValues, error) {
	var list []entity.OptionValues
	if err := r.db.WithContext(ctx).Where("option_category_id = ?", categoryId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
