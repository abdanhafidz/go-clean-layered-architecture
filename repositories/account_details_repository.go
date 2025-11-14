package repositories

import (
	"context"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountDetailRepository interface {
	CreateAccountDetail(ctx context.Context, details entity.AccountDetail) (entity.AccountDetail, error)
	GetAccountDetailById(ctx context.Context, id uuid.UUID) (entity.AccountDetail, error)
	GetAccountDetailByAccountId(ctx context.Context, accountId uuid.UUID) (entity.AccountDetail, error)
	GetAllAccountDetail(ctx context.Context) ([]entity.AccountDetail, error)
	UpdateAccountDetail(ctx context.Context, details entity.AccountDetail) (entity.AccountDetail, error)
	SoftDeleteAccountDetail(ctx context.Context, id uuid.UUID) error
	DeleteAccountDetail(ctx context.Context, id uuid.UUID) error
}

type accountDetailRepository struct {
	db *gorm.DB
}

func NewAccountDetailRepository(db *gorm.DB) AccountDetailRepository {
	return &accountDetailRepository{db: db}
}

func (r *accountDetailRepository) CreateAccountDetail(ctx context.Context, details entity.AccountDetail) (entity.AccountDetail, error) {
	if err := r.db.WithContext(ctx).Create(&details).Error; err != nil {
		return entity.AccountDetail{}, err
	}
	return details, nil
}

func (r *accountDetailRepository) GetAccountDetailById(ctx context.Context, id uuid.UUID) (entity.AccountDetail, error) {
	var details entity.AccountDetail
	if err := r.db.WithContext(ctx).Preload("Account").First(&details, "id = ?", id).Error; err != nil {
		return entity.AccountDetail{}, err
	}
	return details, nil
}

func (r *accountDetailRepository) GetAccountDetailByAccountId(ctx context.Context, accountId uuid.UUID) (entity.AccountDetail, error) {
	var details entity.AccountDetail
	if err := r.db.WithContext(ctx).Preload("Account").First(&details, "account_id = ?", accountId).Error; err != nil {
		return entity.AccountDetail{}, err
	}
	return details, nil
}

func (r *accountDetailRepository) GetAllAccountDetail(ctx context.Context) ([]entity.AccountDetail, error) {
	var list []entity.AccountDetail
	if err := r.db.WithContext(ctx).Preload("Account").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *accountDetailRepository) UpdateAccountDetail(ctx context.Context, details entity.AccountDetail) (entity.AccountDetail, error) {
	var existing entity.AccountDetail
	if err := r.db.WithContext(ctx).First(&existing, "account_id = ?", details.AccountId).Error; err != nil {
		return entity.AccountDetail{}, err
	}
	if err := r.db.WithContext(ctx).Model(&existing).Updates(details).Error; err != nil {
		return entity.AccountDetail{}, err
	}
	return existing, nil
}

func (r *accountDetailRepository) SoftDeleteAccountDetail(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.AccountDetail{}, "id = ?", id).Error
}

func (r *accountDetailRepository) DeleteAccountDetail(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.AccountDetail{}, "id = ?", id).Error
}
