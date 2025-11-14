package services

import (
	"context"
	"errors"
	"regexp"
	"strings"

	dto "abdanhafidz.com/go-clean-layered-architecture/models/dto"
	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	http_error "abdanhafidz.com/go-clean-layered-architecture/models/error"
	"abdanhafidz.com/go-clean-layered-architecture/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AccountService interface {
	GetByEmail(ctx context.Context, email string) (entity.Account, error)
	Create(ctx context.Context, name string, email string, username string, password string) (entity.Account, error)
	Update(ctx context.Context, account entity.Account) (entity.Account, error)
	Validate(ctx context.Context, emailorusername string, password string) (dto.AuthenticatedUser, error)
	ChangePassword(ctx context.Context, accountId uuid.UUID, oldPassword string, newPassword string) (dto.AuthenticatedUser, error)
	GetDetail(ctx context.Context, accountId uuid.UUID) (dto.AccountDetailResponse, error)
	CreateEmptyDetail(ctx context.Context, accountId uuid.UUID) (dto.AccountDetailResponse, error)
	UpdateDetail(ctx context.Context, details entity.AccountDetail) (dto.AccountDetailResponse, error)
}

type accountService struct {
	jwtService        JWTService
	accountRepo       repositories.AccountRepository
	accountDetailRepo repositories.AccountDetailRepository
}

func NewAccountService(jwtService JWTService, accountRepo repositories.AccountRepository, accountDetailRepo repositories.AccountDetailRepository) AccountService {
	return &accountService{
		jwtService:        jwtService,
		accountRepo:       accountRepo,
		accountDetailRepo: accountDetailRepo,
	}
}

func (s *accountService) GetByEmail(ctx context.Context, email string) (entity.Account, error) {
	return s.accountRepo.GetAccountByEmail(ctx, email)
}
func (s *accountService) Create(ctx context.Context, name string, email string, username string, password string) (entity.Account, error) {
	if email == "" || username == "" || password == "" {
		return entity.Account{}, http_error.BAD_REQUEST_ERROR
	}

	if _, err := s.accountRepo.GetAccountByEmail(ctx, email); err == nil {
		return entity.Account{}, http_error.EMAIL_ALREADY_EXISTS
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Account{}, err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return entity.Account{}, err
	}

	acc := entity.Account{Email: email, Username: username, Password: string(bytes), Role: "user"}
	created, err := s.accountRepo.CreateAccount(ctx, acc)

	if err != nil {
		return entity.Account{}, err
	}

	_, err = s.CreateEmptyDetail(ctx, created.Id)

	return created, nil

}

func (s *accountService) Update(ctx context.Context, account entity.Account) (entity.Account, error) {
	return s.accountRepo.UpdateAccount(ctx, account)
}
func (s *accountService) Validate(ctx context.Context, emailorusername string, password string) (dto.AuthenticatedUser, error) {
	acc, err := s.accountRepo.GetAccountByEmail(ctx, emailorusername)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		acc, err = s.accountRepo.GetAccountByUsername(ctx, emailorusername)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.AuthenticatedUser{}, errors.New("account not found")
		}
	}
	if err != nil {
		return dto.AuthenticatedUser{}, err
	}

	if err := s.jwtService.VerifyPassword(ctx, acc.Password, password); err != nil {
		return dto.AuthenticatedUser{}, errors.New("invalid credentials")
	}

	token, err := s.jwtService.GenerateToken(ctx, dto.JWTCustomClaims{
		AccountId: acc.Id.String(),
	})

	if err != nil {
		return dto.AuthenticatedUser{}, err
	}

	return dto.AuthenticatedUser{Account: acc, Token: token}, nil
}

func (s *accountService) ChangePassword(ctx context.Context, accountId uuid.UUID, oldPassword string, newPassword string) (dto.AuthenticatedUser, error) {
	acc, err := s.accountRepo.GetAccountById(ctx, accountId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.AuthenticatedUser{}, errors.New("account not found")
	}
	if err != nil {
		return dto.AuthenticatedUser{}, err
	}

	if err := s.jwtService.VerifyPassword(ctx, acc.Password, oldPassword); err != nil {
		return dto.AuthenticatedUser{}, errors.New("incorrect old password!")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)

	if err != nil {
		return dto.AuthenticatedUser{}, err
	}

	acc.Password = string(bytes)
	acc, err = s.accountRepo.UpdateAccount(ctx, acc)
	if err != nil {
		return dto.AuthenticatedUser{}, err
	}
	return dto.AuthenticatedUser{Account: acc}, nil
}

func sanitizePhone(input string) string {
	p := strings.TrimSpace(input)
	p = strings.ReplaceAll(p, " ", "")
	re := regexp.MustCompile(`[^0-9+]`)
	p = re.ReplaceAllString(p, "")
	if strings.HasPrefix(p, "0") {
		p = "+62" + p[1:]
	}
	if !strings.HasPrefix(p, "+62") && !strings.HasPrefix(p, "+") {
		p = "+" + p
	}
	return p
}

func (s *accountService) GetDetail(ctx context.Context, accountId uuid.UUID) (dto.AccountDetailResponse, error) {
	d, err := s.accountDetailRepo.GetAccountDetailByAccountId(ctx, accountId)
	if err != nil {
		return dto.AccountDetailResponse{}, err
	}
	acc, err := s.accountRepo.GetAccountById(ctx, accountId)
	if err != nil {
		return dto.AccountDetailResponse{}, err
	}
	acc.Password = "SECRET"
	return dto.AccountDetailResponse{Account: acc, Details: entity.AccountDetail(d)}, nil
}

func (s *accountService) CreateEmptyDetail(ctx context.Context, accountID uuid.UUID) (dto.AccountDetailResponse, error) {
	d, err := s.accountDetailRepo.CreateAccountDetail(ctx, entity.AccountDetail{AccountId: accountID})
	if err != nil {
		return dto.AccountDetailResponse{}, err
	}
	acc, err := s.accountRepo.GetAccountById(ctx, accountID)
	if err != nil {
		return dto.AccountDetailResponse{}, err
	}
	acc.IsDetailCompleted = false
	_, _ = s.accountRepo.UpdateAccount(ctx, acc)
	acc.Password = "SECRET"
	return dto.AccountDetailResponse{Account: acc, Details: entity.AccountDetail(d)}, nil
}

func (s *accountService) UpdateDetail(ctx context.Context, details entity.AccountDetail) (dto.AccountDetailResponse, error) {

	if details.PhoneNumber != nil {
		v := sanitizePhone(*details.PhoneNumber)
		details.PhoneNumber = &v
	}

	d, err := s.accountDetailRepo.UpdateAccountDetail(ctx, details)
	if err != nil {
		return dto.AccountDetailResponse{}, err
	}

	acc, err := s.accountRepo.GetAccountById(ctx, details.AccountId)
	if err != nil {
		return dto.AccountDetailResponse{}, err
	}
	acc.IsDetailCompleted = (d.FullName != nil && d.PhoneNumber != nil && d.SchoolName != nil && d.Province != nil && d.City != nil)
	_, _ = s.accountRepo.UpdateAccount(ctx, acc)
	return dto.AccountDetailResponse{Account: acc, Details: entity.AccountDetail(d)}, nil
}
