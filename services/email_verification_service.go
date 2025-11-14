package services

import (
	"context"
	"errors"
	"time"

	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	http_error "abdanhafidz.com/go-clean-layered-architecture/models/error"
	"abdanhafidz.com/go-clean-layered-architecture/repositories"

	"gorm.io/gorm"
)

type EmailVerificationService interface {
	CreateToken(ctx context.Context, email string, token uint, due time.Time) (entity.EmailVerification, error)
	VerifyToken(ctx context.Context, email string, token uint) error
	DeleteByToken(ctx context.Context, token uint) error
}

type emailVerificationService struct {
	accountService        AccountService
	emailVerificationRepo repositories.EmailVerificationRepository
}

func NewEmailVerificationService(accountService AccountService, emailVerificationRepo repositories.EmailVerificationRepository) EmailVerificationService {
	return &emailVerificationService{accountService: accountService, emailVerificationRepo: emailVerificationRepo}
}

func (s *emailVerificationService) CreateToken(ctx context.Context, email string, token uint, due time.Time) (entity.EmailVerification, error) {
	acc, err := s.accountService.GetByEmail(ctx, email)
	if err != nil {
		return entity.EmailVerification{}, err
	}
	if due.IsZero() {
		due = time.Now().Add(15 * time.Minute)
	}
	ev := entity.EmailVerification{AccountId: acc.Id, Token: token, IsExpired: false, CreatedAt: time.Now(), ExpiredAt: due}
	return s.emailVerificationRepo.Create(ctx, ev)
}

func (s *emailVerificationService) VerifyToken(ctx context.Context, email string, token uint) error {
	acc, err := s.accountService.GetByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("account not found")
	}
	if err != nil {
		return err
	}

	ev, err := s.emailVerificationRepo.GetByAccountAndToken(ctx, acc.Id, token)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http_error.INVALID_OTP
	}
	if err != nil {
		return err
	}

	if ev.ExpiredAt.Before(time.Now()) {
		_ = s.emailVerificationRepo.MarkExpired(ctx, ev.Id)
		_ = s.emailVerificationRepo.DeleteByToken(ctx, ev.Token)
		return http_error.EXPIRED_TOKEN
	}

	acc.IsEmailVerified = true
	if _, err := s.accountService.Update(ctx, acc); err != nil {
		return err
	}
	return s.emailVerificationRepo.MarkExpired(ctx, ev.Id)
}

func (s *emailVerificationService) DeleteByToken(ctx context.Context, token uint) error {
	return s.emailVerificationRepo.DeleteByToken(ctx, token)
}
