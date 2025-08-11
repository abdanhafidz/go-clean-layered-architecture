package repositories

import (
	"context"

	"abdanhafidz.com/go-boilerplate/models/entity"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, email string, name string, password string) (res entity.Account, err error)
	GetAccountByEmail(ctx context.Context, email string) (res entity.Account, err error)
}
type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) GetAccountByEmail(ctx context.Context, email string) (res entity.Account, err error) {
	tx := r.db.WithContext(ctx).Where("email = ?", email).First(&res)
	err = tx.Error
	return res, err
}
func (r *accountRepository) CreateAccount(ctx context.Context, email string, name string, password string) (res entity.Account, err error) {

	res = entity.Account{
		Email:    email,
		Name:     name,
		Password: password,
	}

	tx := r.db.WithContext(ctx).Create(&res)

	err = tx.Error

	return res, err
}
