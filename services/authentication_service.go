package services

import (
	"context"
	"errors"

	"abdanhafidz.com/go-boilerplate/models/dto"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (res dto.AuthenticatedResponse, err error)
	Login(ctx context.Context, email string, password string) (res dto.AuthenticatedResponse, err error)
}

type authenticationService struct {
	accountRepository repositories.AccountRepository
	jwtService        JWTService
}

func NewAuthenticationService(accountRepository repositories.AccountRepository, jwtService JWTService) AuthenticationService {
	return &authenticationService{
		accountRepository: accountRepository,
		jwtService:        jwtService,
	}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func verifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("invalid password")
	}
	return nil
}

func (s *authenticationService) Register(ctx context.Context, req dto.RegisterRequest) (res dto.AuthenticatedResponse, err error) {

	hashed_password, err_hash := hashPassword(req.Password)

	if err_hash != nil {
		return res, http_error.INTERNAL_SERVER_ERROR
	}

	account, err_account := s.accountRepository.CreateAccount(ctx, req.Email, req.Name, hashed_password)

	if err_repo := RepoError(err_account); err_repo != nil {
		return res, err_repo
	}

	return s.Login(ctx, account.Email, account.Password)

}

func (s *authenticationService) Login(ctx context.Context, email string, password string) (res dto.AuthenticatedResponse, err error) {
	account, err_account := s.accountRepository.GetAccountByEmail(ctx, email)

	if err_repo := RepoError(err_account); err_repo != nil {
		return res, err_repo
	}

	if err_password := verifyPassword(account.Password, password); err_password != nil {
		return res, http_error.WRONG_PASSWORD
	}

	auth_token, err_jwt := s.jwtService.GenerateToken(ctx, dto.JWTCustomClaims{
		IdUser: account.Id,
	})

	if err_jwt != nil {
		return res, http_error.INVALID_TOKEN
	}

	res = dto.AuthenticatedResponse{
		Account: account,
		Token:   auth_token,
	}

	return res, nil
}
