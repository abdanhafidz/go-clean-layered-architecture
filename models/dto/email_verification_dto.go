package dto

type DeleteEmailVerificationRequest struct {
	Token uint `json:"token" binding:"required"`
}
