package repositories

import (
	"context"
	"strconv"
	"time"

	"abdanhafidz.com/go-boilerplate/models/entity"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, name string, date_of_birth time.Time, account_number uint) (res entity.Account, err error)
	GetAccount(ctx context.Context, name string, date_of_birth time.Time, digits_account_number string) (res entity.Account, err error)
	ResetPIN(ctx context.Context, account_id uuid.UUID, newPin int) (res entity.Account, err error)
	UpdateAccountStatus(ctx context.Context, account_id uuid.UUID, status bool) (res entity.Account, err error)
}
type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

func (r *accountRepository) CreateAccount(ctx context.Context, name string, date_of_birth time.Time, account_number uint) (res entity.Account, err error) {

	res = entity.Account{
		Name:          name,
		DateOfBirth:   date_of_birth,
		AccountNumber: account_number,
		IsActive:      true,
	}

	tx := r.db.WithContext(ctx).Create(&res)

	err = tx.Error

	return res, err
}

func (r *accountRepository) GetAccount(ctx context.Context, name string, date_of_birth time.Time, digits_account_number string) (res entity.Account, err error) {
	// Cari account dengan filter nama & tanggal lahir
	tx := r.db.WithContext(ctx).Where("name = ? AND date_of_birth = ?", name, date_of_birth).Find(&res)

	if tx.Error != nil {
		return res, tx.Error
	}

	if tx.RowsAffected == 0 {
		return res, gorm.ErrRecordNotFound
	}

	if strconv.Itoa(int(res.AccountNumber%1000)) != digits_account_number {
		return res, gorm.ErrInvalidValue
	}
	return res, tx.Error
}

func (r *accountRepository) ResetPIN(ctx context.Context, account_id uuid.UUID, newPin int) (res entity.Account, err error) {
	tx := r.db.WithContext(ctx).Model(&entity.Account{Id: account_id}).Update("pin", newPin).First(&res)

	if tx.Error != nil {
		return res, tx.Error
	}

	if tx.RowsAffected == 0 {
		return res, gorm.ErrRecordNotFound
	}

	return res, tx.Error
}

func (r *accountRepository) UpdateAccountStatus(ctx context.Context, account_id uuid.UUID, status bool) (res entity.Account, err error) {
	tx := r.db.WithContext(ctx).Model(&entity.Account{Id: account_id}).Update("is_active", status).First(&res)

	if tx.Error != nil {
		return res, tx.Error
	}

	if tx.RowsAffected == 0 {
		return res, gorm.ErrRecordNotFound
	}

	return res, tx.Error
}
