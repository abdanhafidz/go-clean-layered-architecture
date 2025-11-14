package dto

import entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"

type SignInRequest struct {
	EmailorUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateEmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required" `
	NewPassword string `json:"new_password" binding:"required" `
}

type ValidateVerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Token uint   `json:"token" binding:"required"`
}

type ExternalAuthRequest struct {
	OauthID       string `json:"oauth_id" binding:"required"`
	OauthProvider string `json:"oauth_provider" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ValidateForgotPasswordRequest struct {
	Token       uint   `json:"token" binding:"required"`
	NewPassword string `json:"new_password"`
}

type AuthenticatedUser struct {
	Account entity.Account `json:"account"`
	Token   string         `json:"token"`
}
