package repositories

import (
	"context"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FCMRepository interface {
	CreateOrUpdate(ctx context.Context, accountId uuid.UUID, token string) (entity.FCM, error)
	GetByAccountId(ctx context.Context, accountId uuid.UUID) (entity.FCM, error)
	DeleteByAccountId(ctx context.Context, accountId uuid.UUID) error
}

type fcmRepository struct {
	db *gorm.DB
}

func NewFCMRepository(db *gorm.DB) FCMRepository {
	return &fcmRepository{db: db}
}

func (r *fcmRepository) CreateOrUpdate(ctx context.Context, accountId uuid.UUID, token string) (entity.FCM, error) {
	var rec entity.FCM
	err := r.db.WithContext(ctx).First(&rec, "account_id = ?", accountId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			rec = entity.FCM{AccountId: accountId, FCMToken: token}
			if err := r.db.WithContext(ctx).Create(&rec).Error; err != nil {
				return entity.FCM{}, err
			}
			return rec, nil
		}
		return entity.FCM{}, err
	}
	// update existing
	if err := r.db.WithContext(ctx).Model(&rec).Update("fcm_token", token).Error; err != nil {
		return entity.FCM{}, err
	}
	return rec, nil
}

func (r *fcmRepository) GetByAccountId(ctx context.Context, accountId uuid.UUID) (entity.FCM, error) {
	var rec entity.FCM
	if err := r.db.WithContext(ctx).First(&rec, "account_id = ?", accountId).Error; err != nil {
		return entity.FCM{}, err
	}
	return rec, nil
}

func (r *fcmRepository) DeleteByAccountId(ctx context.Context, accountId uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.FCM{}, "account_id = ?", accountId).Error
}
