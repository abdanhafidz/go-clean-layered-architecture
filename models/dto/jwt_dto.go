package dto

import (
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

type JWTCustomClaims struct {
	IdUser uuid.UUID `json:"user_id" binding:"required"`
	jwt.RegisteredClaims
}

type AccountData struct {
	IdUser uuid.UUID `json:"user_id" binding:"required"`
}
