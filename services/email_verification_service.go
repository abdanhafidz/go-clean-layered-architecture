package services

import (
	"context"
	"errors"
	"time"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"

	"gorm.io/gorm"
)

type EmailVerificationService interface {
	CreateToken(ctx context.Context, email string, token uint, due time.Time) (entity.EmailVerification, error)
	VerifyToken(ctx context.Context, email string, token uint) error
	DeleteByToken(ctx context.Context, token uint) error
}

type emailVerificationService struct {
	accountRepo            repositories.AccountRepository
	emailsVerificationRepo repositories.EmailVerificationRepository
}

func NewEmailVerificationService(accountRepo repositories.AccountRepository, emailsVerificationRepo repositories.EmailVerificationRepository) EmailVerificationService {
	return &emailVerificationService{accountRepo: accountRepo, emailsVerificationRepo: emailsVerificationRepo}
}

func (s *emailVerificationService) CreateToken(ctx context.Context, email string, token uint, due time.Time) (entity.EmailVerification, error) {
	acc, err := s.accountRepo.GetAccountByEmail(ctx, email)
	if err != nil {
		return entity.EmailVerification{}, err
	}
	if due.IsZero() {
		due = time.Now().Add(15 * time.Minute)
	}
	ev := entity.EmailVerification{AccountId: acc.Id, Token: token, IsExpired: false, CreatedAt: time.Now(), ExpiredAt: due}
	return s.emailsVerificationRepo.Create(ctx, ev)
}

func (s *emailVerificationService) VerifyToken(ctx context.Context, email string, token uint) error {
	acc, err := s.accountRepo.GetAccountByEmail(ctx, email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("account not found")
	}
	if err != nil {
		return err
	}

	ev, err := s.emailsVerificationRepo.GetByAccountAndToken(ctx, acc.Id, token)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return http_error.INVALID_TOKEN
	}
	if err != nil {
		return err
	}

	if ev.ExpiredAt.Before(time.Now()) {
		_ = s.emailsVerificationRepo.MarkExpired(ctx, ev.Id)
		_ = s.emailsVerificationRepo.DeleteByToken(ctx, ev.Token)
		return http_error.EXPIRED_TOKEN
	}

	acc.IsEmailVerified = true
	if _, err := s.accountRepo.UpdateAccount(ctx, acc); err != nil {
		return err
	}
	return s.emailsVerificationRepo.MarkExpired(ctx, ev.Id)
}

func (s *emailVerificationService) DeleteByToken(ctx context.Context, token uint) error {
	return s.emailsVerificationRepo.DeleteByToken(ctx, token)
}
