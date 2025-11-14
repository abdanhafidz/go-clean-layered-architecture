package repositories

import (
	"context"
	"time"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"gorm.io/gorm"
)

type ForgotPasswordRepository interface {
	Create(ctx context.Context, rec entity.ForgotPassword) (entity.ForgotPassword, error)
	GetByToken(ctx context.Context, token uint) (entity.ForgotPassword, error)
	MarkExpired(ctx context.Context, id interface{}) error
	DeleteByToken(ctx context.Context, token uint) error
	ExpireAllOverdue(ctx context.Context, now time.Time) (int64, error)
}

type forgotPasswordRepository struct {
	db *gorm.DB
}

func NewForgotPasswordRepository(db *gorm.DB) ForgotPasswordRepository {
	return &forgotPasswordRepository{db: db}
}

func (r *forgotPasswordRepository) Create(ctx context.Context, rec entity.ForgotPassword) (entity.ForgotPassword, error) {
	if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
		return entity.ForgotPassword{}, err
	}
	return rec, nil
}

func (r *forgotPasswordRepository) GetByToken(ctx context.Context, token uint) (entity.ForgotPassword, error) {
	var res entity.ForgotPassword
	if err := r.db.WithContext(ctx).Where("token = ? AND is_expired = ?", token, false).First(&res).Error; err != nil {
		return entity.ForgotPassword{}, err
	}
	return res, nil
}

func (r *forgotPasswordRepository) MarkExpired(ctx context.Context, id interface{}) error {
	return r.db.WithContext(ctx).
		Model(&entity.ForgotPassword{}).
		Where("id = ?", id).
		Update("is_expired", true).Error
}

func (r *forgotPasswordRepository) DeleteByToken(ctx context.Context, token uint) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&entity.ForgotPassword{}).Error
}

func (r *forgotPasswordRepository) ExpireAllOverdue(ctx context.Context, now time.Time) (int64, error) {
	tx := r.db.WithContext(ctx).
		Model(&entity.ForgotPassword{}).
		Where("is_expired = ? AND expired_at <= ?", false, now).
		Update("is_expired", true)
	return tx.RowsAffected, tx.Error
}
