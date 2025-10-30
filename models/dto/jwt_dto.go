package dto

import (
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/google/uuid"
)

type JWTCustomClaims struct {
	AccountId string `json:"account_id" binding:"required"`
	jwt.RegisteredClaims
}

type AccountData struct {
	AccountId uuid.UUID `json:"account_id" binding:"required"`
}
