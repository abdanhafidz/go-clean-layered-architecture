package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"gorm.io/gorm"
)

type AccountService interface {
	CreateAccount(ctx context.Context, name string, date_of_birth time.Time) (res entity.Account, err error)
	VerifyAccount(ctx context.Context, req dto.VerifyAccountRequest) (res entity.Account, err error)
	ResetPIN(ctx context.Context, req dto.VerifyAccountRequest) (res entity.Account, err error)
	BlockAccount(ctx context.Context, req dto.VerifyAccountRequest) (res entity.Account, err error)
}

type accountService struct {
	accountRepo repositories.AccountRepository
}

func NewAccountService(accountRepo repositories.AccountRepository) AccountService {
	return &accountService{
		accountRepo: accountRepo,
	}
}
func (s *accountService) CreateAccount(ctx context.Context, name string, date_of_birth time.Time) (res entity.Account, err error) {
	account_number := uint(rand.Intn(90_000_000) + 10_000_000)
	if name == "" || date_of_birth.IsZero() {
		return res, http_error.BAD_REQUEST_ERROR
	}
	res, err = s.accountRepo.CreateAccount(ctx, name, date_of_birth, account_number)
	if err != nil {
		return res, err
	}
	return res, err
}

func (s *accountService) VerifyAccount(ctx context.Context, req dto.VerifyAccountRequest) (res entity.Account, err error) {
	res, err = s.accountRepo.GetAccount(ctx, req.Name, req.Dateofbirth, req.LastDigit)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, http_error.DATA_NOT_FOUND
	}

	if errors.Is(err, gorm.ErrInvalidValue) {
		return res, http_error.INVALID_ACCOUNT_DIGITS
	}

	return res, err
}

func (s *accountService) ResetPIN(ctx context.Context, req dto.VerifyAccountRequest) (res entity.Account, err error) {
	res_verify, err_verify := s.VerifyAccount(ctx, req)

	if err_verify != nil {
		return res, http_error.UNAUTHORIZED
	}

	newPIN := int(rand.Intn(90_000) + 10_000)

	res, err = s.accountRepo.ResetPIN(ctx, res_verify.Id, newPIN)

	return res, err
}

func (s *accountService) BlockAccount(ctx context.Context, req dto.VerifyAccountRequest) (res entity.Account, err error) {
	res_verify, err_verify := s.VerifyAccount(ctx, req)

	if err_verify != nil {
		return res, http_error.UNAUTHORIZED
	}
	fmt.Println(res_verify)
	res, err = s.accountRepo.UpdateAccountStatus(ctx, res_verify.Id, false)

	return res, err
}
