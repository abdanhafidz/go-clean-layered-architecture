package dto

import "abdanhafidz.com/go-boilerplate/models/entity"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AuthenticatedResponse struct {
	Account entity.Account `json:"account"`
	Token   string         `json:"authorization_token:"`
}
