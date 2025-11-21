package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
)

type FileRepository interface {
	Create(ctx context.Context, file *entity.File) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.File, error)
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *entity.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.File, error) {
	var file entity.File
	result := r.db.WithContext(ctx).First(&file, "id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, http_error.NOT_FOUND_ERROR
		}
		return nil, result.Error
	}

	return &file, nil
}