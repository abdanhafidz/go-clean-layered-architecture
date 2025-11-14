package services

import (
	"context"
	"errors"
	"strings"
	"time"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	http_error "abdanhafidz.com/go-clean-layered-architecture/models/error"
	"abdanhafidz.com/go-clean-layered-architecture/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ForgotPasswordService interface {
	Request(ctx context.Context, email string, token uint, due time.Time) (entity.ForgotPassword, error)
	Reset(ctx context.Context, token uint, newPassword string) error
}

type forgotPasswordService struct {
	jwtService         JWTService
	accountRepo        repositories.AccountRepository
	forgotPasswordRepo repositories.ForgotPasswordRepository
}

func NewForgotPasswordService(jwtService JWTService, accountRepo repositories.AccountRepository, forgotPasswordRepo repositories.ForgotPasswordRepository) ForgotPasswordService {
	return &forgotPasswordService{
		jwtService:         jwtService,
		accountRepo:        accountRepo,
		forgotPasswordRepo: forgotPasswordRepo}
}

func (s *forgotPasswordService) Request(ctx context.Context, email string, token uint, due time.Time) (entity.ForgotPassword, error) {
	acc, err := s.accountRepo.GetAccountByEmail(ctx, email)
	if err != nil {
		return entity.ForgotPassword{}, err
	}
	if due.IsZero() {
		due = time.Now().Add(15 * time.Minute)
	}
	rec := entity.ForgotPassword{AccountId: acc.Id, Token: token, IsExpired: false, CreatedAt: time.Now(), ExpiredAt: due}
	return s.forgotPasswordRepo.Create(ctx, rec)
}

func (s *forgotPasswordService) Reset(ctx context.Context, token uint, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" {
		return http_error.BAD_REQUEST_ERROR
	}

	rec, err := s.forgotPasswordRepo.GetByToken(ctx, token)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http_error.INVALID_OTP
	}

	if err != nil {
		return err
	}

	if rec.ExpiredAt.Before(time.Now()) {
		_ = s.forgotPasswordRepo.MarkExpired(ctx, rec.Id)
		return http_error.EXPIRED_TOKEN
	}

	acc, err := s.accountRepo.GetAccountById(ctx, rec.AccountId)
	if err != nil {
		return err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)

	if err != nil {
		return err
	}

	acc.Password = string(bytes)

	if _, err := s.accountRepo.UpdateAccount(ctx, acc); err != nil {
		return err
	}

	return s.forgotPasswordRepo.MarkExpired(ctx, rec.Id)
}
