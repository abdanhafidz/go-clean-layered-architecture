package services

import (
	"context"
	"fmt"

	"abdanhafidz.com/go-clean-layered-architecture/models/dto"
	http_error "abdanhafidz.com/go-clean-layered-architecture/models/error"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type JWTService interface {
	GenerateToken(ctx context.Context, payload dto.JWTCustomClaims) (token string, err error)
	ValidateToken(ctx context.Context, tokenStr string) (claim *dto.JWTCustomClaims, err error)
	VerifyPassword(ctx context.Context, hashedPassword string, password string) error
}

type jwtService struct {
	secretKey string
}

func NewJWTService(secretKey string) JWTService {
	return &jwtService{
		secretKey: secretKey,
	}
}

func (s *jwtService) GenerateToken(ctx context.Context, payload dto.JWTCustomClaims) (token string, err error) {

	claims := jwt.MapClaims{
		"account_id": payload.AccountId,
	}

	fmt.Println(s.secretKey)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err_convertion := jwtToken.SignedString([]byte(s.secretKey))

	if err_convertion != nil {
		return "", http_error.INTERNAL_SERVER_ERROR
	}

	return token, nil
}
func (s *jwtService) VerifyPassword(ctx context.Context, hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return http_error.WRONG_PASSWORD
	}
	return nil
}

func (s *jwtService) ValidateToken(ctx context.Context, tokenStr string) (claim *dto.JWTCustomClaims, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", http_error.INTERNAL_SERVER_ERROR
		}
		return []byte(s.secretKey), nil
	})

	fmt.Println("Token", token)
	fmt.Println("secretKey", s.secretKey)

	if err != nil || !token.Valid {
		return nil, http_error.INVALID_TOKEN
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, http_error.INTERNAL_SERVER_ERROR
	}
	account_id, ok := claims["account_id"].(string)
	if !ok {
		return nil, http_error.INTERNAL_SERVER_ERROR
	}
	return &dto.JWTCustomClaims{
		AccountId: account_id,
	}, nil
}
