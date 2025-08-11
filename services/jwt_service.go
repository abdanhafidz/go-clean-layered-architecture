package services

import (
	"context"

	"abdanhafidz.com/go-boilerplate/models/dto"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

type JWTService interface {
	GenerateToken(ctx context.Context, payload dto.JWTCustomClaims) (token string, err error)
	ValidateToken(ctx context.Context, tokenStr string) (claim *dto.JWTCustomClaims, err error)
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
		"user_id": payload.IdUser,
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err_convertion := jwtToken.SignedString([]byte(s.secretKey))

	if err_convertion != nil {
		return "", http_error.INTERNAL_SERVER_ERROR
	}

	return token, nil
}

func (s *jwtService) ValidateToken(ctx context.Context, tokenStr string) (claim *dto.JWTCustomClaims, err error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", http_error.INTERNAL_SERVER_ERROR
		}
		return []byte(s.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, http_error.INVALID_TOKEN
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, http_error.INVALID_TOKEN
	}

	return &dto.JWTCustomClaims{
		IdUser: claims["user_id"].(uuid.UUID),
	}, nil
}
