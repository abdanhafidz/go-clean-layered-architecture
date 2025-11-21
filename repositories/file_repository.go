package repositories

import (
    "context"
    "gorm.io/gorm"

    entity "abdanhafidz.com/go-boilerplate/models/entity"
)

type FileRepository interface {
    Create(ctx context.Context, file *entity.File) error
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