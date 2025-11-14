package dto

import (
	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
)

type AccountDetailResponse struct {
	Account entity.Account       `json:"account"`
	Details entity.AccountDetail `json:"details"`
}

type UpdateAccountDetailRequest struct {
	FullName    *string `json:"full_name"`
	SchoolName  *string `json:"school_name"`
	Province    *string `json:"province"`
	City        *string `json:"city"`
	Avatar      *string `json:"avatar"`
	PhoneNumber *string `json:"phone_number"`
}
