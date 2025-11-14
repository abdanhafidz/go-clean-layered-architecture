package repositories

import (
	"context"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, account entity.Account) (entity.Account, error)
	GetAccountById(ctx context.Context, accountId uuid.UUID) (entity.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (entity.Account, error)
	GetAccountByUsername(ctx context.Context, username string) (entity.Account, error)
	GetAllaccount(ctx context.Context) ([]entity.Account, error)
	UpdateAccount(ctx context.Context, account entity.Account) (entity.Account, error)
	SoftDeleteAccount(ctx context.Context, accountId uuid.UUID) error
	DeleteAccount(ctx context.Context, accountId uuid.UUID) error
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) CreateAccount(ctx context.Context, account entity.Account) (entity.Account, error) {
	if err := r.db.WithContext(ctx).Create(&account).Error; err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (r *accountRepository) GetAccountById(ctx context.Context, accountId uuid.UUID) (entity.Account, error) {
	var account entity.Account
	if err := r.db.WithContext(ctx).First(&account, "id = ?", accountId).Error; err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (r *accountRepository) GetAccountByEmail(ctx context.Context, email string) (entity.Account, error) {
	var account entity.Account
	if err := r.db.WithContext(ctx).First(&account, "email = ?", email).Error; err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (r *accountRepository) GetAccountByUsername(ctx context.Context, username string) (entity.Account, error) {
	var account entity.Account
	if err := r.db.WithContext(ctx).First(&account, "username = ?", username).Error; err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (r *accountRepository) GetAllaccount(ctx context.Context) ([]entity.Account, error) {
	var account []entity.Account
	if err := r.db.WithContext(ctx).Find(&account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

func (r *accountRepository) UpdateAccount(ctx context.Context, account entity.Account) (entity.Account, error) {
	if err := r.db.WithContext(ctx).Save(&account).Error; err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (r *accountRepository) SoftDeleteAccount(ctx context.Context, accountId uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Account{}, "id = ?", accountId).Error
}

func (r *accountRepository) DeleteAccount(ctx context.Context, accountId uuid.UUID) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&entity.Account{}, "id = ?", accountId).Error
}
