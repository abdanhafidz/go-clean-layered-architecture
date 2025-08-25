package dto

import "time"

type AccountRequest struct {
	Dateofbirth time.Time `json:"dob" binding:"required"`
	Name        string    `json:"name" binding:"required"`
}

type VerifyAccountRequest struct {
	Dateofbirth time.Time `json:"dob" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	LastDigit   string    `json:"account_3_digits" binding:"required"`
}

type AccountInfoRequest struct {
	VerifyAccountRequest
}

type AccountInfoResponse struct {
	Name          string `json:"name" binding:"required"`
	AccountNumber uint   `gorm:"uniqueIndex" json:"account_number"`
	Balance       uint   `gorm:"type:number" json:"balance"`
}
