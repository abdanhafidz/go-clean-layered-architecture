package repositories

import (
	"context"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExternalAuthRepository interface {
	Create(ctx context.Context, oauth entity.ExternalAuth) (entity.ExternalAuth, error)
	GetByAccountId(ctx context.Context, accountId uuid.UUID) ([]entity.ExternalAuth, error)
	GetByOauthId(ctx context.Context, oauthId string) (entity.ExternalAuth, error)
	DeleteById(ctx context.Context, id uuid.UUID) error
}

type externalAuthRepository struct {
	db *gorm.DB
}

func NewExternalAuthRepository(db *gorm.DB) ExternalAuthRepository {
	return &externalAuthRepository{db: db}
}

func (r *externalAuthRepository) Create(ctx context.Context, oauth entity.ExternalAuth) (entity.ExternalAuth, error) {
	if err := r.db.WithContext(ctx).Create(&oauth).Error; err != nil {
		return entity.ExternalAuth{}, err
	}
	return oauth, nil
}

func (r *externalAuthRepository) GetByAccountId(ctx context.Context, accountId uuid.UUID) ([]entity.ExternalAuth, error) {
	var list []entity.ExternalAuth
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *externalAuthRepository) GetByOauthId(ctx context.Context, oauthId string) (entity.ExternalAuth, error) {
	var res entity.ExternalAuth
	if err := r.db.WithContext(ctx).First(&res, "oauth_id = ?", oauthId).Error; err != nil {
		return entity.ExternalAuth{}, err
	}
	return res, nil
}

func (r *externalAuthRepository) DeleteById(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ExternalAuth{}, "id = ?", id).Error
}
