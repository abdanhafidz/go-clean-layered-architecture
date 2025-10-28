package services

import (
	"context"
	"errors"

	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/repositories"
	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type ExternalAuthService interface {
	GoogleAuth(ctx context.Context, idToken string) (dto.AuthenticatedUser, error)
}

type externalAuthService struct {
	jwtService       JWTService
	accountService   AccountService
	externalAuthRepo repositories.ExternalAuthRepository
}

func NewExternalAuthService(jwtService JWTService, accountService AccountService, externalAuthRepo repositories.ExternalAuthRepository) ExternalAuthService {
	return &externalAuthService{
		jwtService:       jwtService,
		accountService:   accountService,
		externalAuthRepo: externalAuthRepo,
	}
}

func (s *externalAuthService) GoogleAuth(ctx context.Context, idToken string) (dto.AuthenticatedUser, error) {

	var (
		acc    entity.Account
		errAcc error
	)
	_, err := s.externalAuthRepo.GetByOauthId(ctx, idToken)
	payload, _ := idtoken.Validate(context.Background(), idToken, "")

	name := payload.Claims["name"].(string)
	email := payload.Claims["email"].(string)
	password := payload.Claims["sub"].(string)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		acc, errAcc = s.accountService.Create(ctx, name, email, name, password)
	}

	if errAcc != nil {
		return dto.AuthenticatedUser{}, err
	}

	token, _ := s.jwtService.GenerateToken(ctx, dto.JWTCustomClaims{
		AccountId: acc.Id,
	})

	return dto.AuthenticatedUser{Account: acc, Token: token}, err

}
