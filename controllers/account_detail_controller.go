package controllers

import (
	"abdanhafidz.com/go-clean-layered-architecture/models/dto"
	entity "abdanhafidz.com/go-clean-layered-architecture/models/entity"
	"abdanhafidz.com/go-clean-layered-architecture/services"
	"github.com/gin-gonic/gin"
)

type AccountDetailController interface {
	GetDetail(ctx *gin.Context)
	UpdateDetail(ctx *gin.Context)
}

type accountDetailController struct {
	accountService services.AccountService
}

func NewAccountDetailController(accountService services.AccountService) AccountDetailController {
	return &accountDetailController{
		accountService: accountService,
	}
}

func (c *accountDetailController) GetDetail(ctx *gin.Context) {
	accountId := ParseAccountId(ctx)
	res, err := c.accountService.GetDetail(ctx.Request.Context(), accountId)
	ResponseJSON(ctx, gin.H{"accountId": accountId}, res, err)
}

func (c *accountDetailController) UpdateDetail(ctx *gin.Context) {
	req := RequestJSON[dto.UpdateAccountDetailRequest](ctx)

	accountId := ParseAccountId(ctx)

	details := entity.AccountDetail{
		AccountId:   accountId,
		FullName:    req.FullName,
		SchoolName:  req.SchoolName,
		Province:    req.Province,
		City:        req.City,
		Avatar:      req.Avatar,
		PhoneNumber: req.PhoneNumber,
	}
	res, err := c.accountService.UpdateDetail(ctx.Request.Context(), details)
	ResponseJSON(ctx, req, res, err)
}
