package services

import (
	"context"

	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/repositories"
	"google.golang.org/api/idtoken"
)

type ExternalAuthService interface {
	GoogleAuth(ctx context.Context, idToken string) (dto.AuthenticatedUser, error)
}

type externalAuthService struct {
	jwtService       JWTService
	accountRepo      repositories.AccountRepository
	externalAuthRepo repositories.ExternalAuthRepository
}

func NewExternalAuthService(jwtService JWTService, accountRepo repositories.AccountRepository, externalAuthRepo repositories.ExternalAuthRepository) ExternalAuthService {
	return &externalAuthService{
		jwtService:       jwtService,
		accountRepo:      accountRepo,
		externalAuthRepo: externalAuthRepo,
	}
}

func (s *externalAuthService) GoogleAuth(ctx context.Context, idToken string) (dto.AuthenticatedUser, error) {
	_, err := s.externalAuthRepo.GetByOauthId(ctx, idToken)

	if err != nil {
		return dto.AuthenticatedUser{}, err
	}

	payload, err := idtoken.Validate(context.Background(), idToken, "")

	if err != nil {
		return dto.AuthenticatedUser{}, err
	}

	email := payload.Claims["email"].(string)

	acc, err := s.accountRepo.GetAccountByEmail(ctx, email)

	if err != nil {
		return dto.AuthenticatedUser{}, err
	}

	token, _ := s.jwtService.GenerateToken(ctx, dto.JWTCustomClaims{
		AccountId: acc.Id,
	})

	return dto.AuthenticatedUser{Account: acc, Token: token}, err

}
