package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/services"
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

// GetDetail godoc
// @Summary      Get Account Detail
// @Description  Retrieve detailed information about the authenticated user's account
// @Tags         Account Detail
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.SuccessResponse[dto.AccountDetailResponse]
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/v1/account/me [get]

func (c *accountDetailController) GetDetail(ctx *gin.Context) {
	accountId := ParseAccountId(ctx)
	res, err := c.accountService.GetDetail(ctx.Request.Context(), accountId)
	ResponseJSON(ctx, gin.H{"accountId": accountId}, res, err)
}

// UpdateDetail godoc
// @Summary      Update Account Detail
// @Description  Update detailed information about the authenticated user's account
// @Tags         Account Detail
// @Accept       json
// @Produce      json
// @Param        request  body      dto.UpdateAccountDetailRequest  true  "Update Account Detail Request"
// @Success      200      {object}  dto.SuccessResponse[dto.AccountDetailResponse]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/account/me [put]
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
