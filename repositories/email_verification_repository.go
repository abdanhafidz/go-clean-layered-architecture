package repositories

import (
	"context"
	"time"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmailVerificationRepository interface {
	Create(ctx context.Context, verification entity.EmailVerification) (entity.EmailVerification, error)
	GetByAccountAndToken(ctx context.Context, accountID uuid.UUID, token uint) (entity.EmailVerification, error)
	MarkExpired(ctx context.Context, id uuid.UUID) error
	DeleteByToken(ctx context.Context, token uint) error
	GetActiveByAccount(ctx context.Context, accountID uuid.UUID) ([]entity.EmailVerification, error)
	ExpireAllOverdue(ctx context.Context, now time.Time) (int64, error)
}

type emailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(ctx context.Context, verification entity.EmailVerification) (entity.EmailVerification, error) {
	if err := r.db.WithContext(ctx).Create(&verification).Error; err != nil {
		return entity.EmailVerification{}, err
	}
	return verification, nil
}

func (r *emailVerificationRepository) GetByAccountAndToken(ctx context.Context, accountID uuid.UUID, token uint) (entity.EmailVerification, error) {
	var ev entity.EmailVerification
	if err := r.db.WithContext(ctx).
		Where("account_id = ? AND token = ? AND is_expired = ?", accountID, token, false).
		First(&ev).Error; err != nil {
		return entity.EmailVerification{}, err
	}
	return ev, nil
}

func (r *emailVerificationRepository) MarkExpired(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&entity.EmailVerification{}).
		Where("id = ?", id).
		Update("is_expired", true).Error
}

func (r *emailVerificationRepository) DeleteByToken(ctx context.Context, token uint) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&entity.EmailVerification{}).Error
}

func (r *emailVerificationRepository) GetActiveByAccount(ctx context.Context, accountID uuid.UUID) ([]entity.EmailVerification, error) {
	var list []entity.EmailVerification
	if err := r.db.WithContext(ctx).
		Where("account_id = ? AND is_expired = ?", accountID, false).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *emailVerificationRepository) ExpireAllOverdue(ctx context.Context, now time.Time) (int64, error) {
	tx := r.db.WithContext(ctx).
		Model(&entity.EmailVerification{}).
		Where("is_expired = ? AND expired_at <= ?", false, now).
		Update("is_expired", true)
	return tx.RowsAffected, tx.Error
}
