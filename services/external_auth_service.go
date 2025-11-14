package services

import (
	"context"
	"errors"

	"abdanhafidz.com/go-clean-layered-architecture/models/dto"
	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"abdanhafidz.com/go-clean-layered-architecture/repositories"
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
		acc        entity.Account
		errAcc     error
		errExtAuth error
		name       string
		email      string
		password   string
	)

	payload, errTok := idtoken.Validate(context.Background(), idToken, "")

	if errTok != nil {
		return dto.AuthenticatedUser{}, errTok
	}
	claims := payload.Claims

	if v, ok := claims["email"].(string); ok {
		email = v
	}

	if v, ok := claims["name"].(string); ok {
		name = v
	} else {
		if v, ok := claims["given_name"].(string); ok {
			name = v
		}
	}

	if v, ok := claims["sub"].(string); ok {
		password = v
	}

	acc, err := s.accountService.GetByEmail(ctx, email)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		acc, errAcc = s.accountService.Create(ctx, name, email, name, password)
		acc.IsEmailVerified = true
		s.accountService.Update(ctx, acc)
		_, errExtAuth = s.externalAuthRepo.Create(ctx, entity.ExternalAuth{
			OauthID:       idToken,
			OauthProvider: "google",
			AccountId:     acc.Id,
		})
	}

	if errAcc != nil {
		return dto.AuthenticatedUser{}, errAcc
	}

	if errExtAuth != nil {
		return dto.AuthenticatedUser{}, errExtAuth
	}
	token, _ := s.jwtService.GenerateToken(ctx, dto.JWTCustomClaims{
		AccountId: acc.Id.String(),
	})

	err = errors.Join(errAcc, errExtAuth)
	return dto.AuthenticatedUser{Account: acc, Token: token}, err

}
