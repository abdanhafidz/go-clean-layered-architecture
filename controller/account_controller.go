package controller

import (
	"fmt"

	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type AccountController interface {
	CreateAccount(ctx *gin.Context)
	VerifyAccount(ctx *gin.Context)
	ResetPIN(ctx *gin.Context)
	BlockAccount(ctx *gin.Context)
}

type accountController struct {
	accountService services.AccountService
}

func NewAccountController(accountService services.AccountService) AccountController {
	return &accountController{
		accountService: accountService,
	}
}

func (c *accountController) CreateAccount(ctx *gin.Context) {
	req := RequestJSON[dto.AccountRequest](ctx)
	fmt.Println(req)
	res, err := c.accountService.CreateAccount(ctx.Request.Context(), req.Name, req.Dateofbirth)
	ResponseJSON(ctx, req, res, err)
}
func (c *accountController) VerifyAccount(ctx *gin.Context) {
	req := RequestJSON[dto.VerifyAccountRequest](ctx)
	fmt.Println(req)
	res, err := c.accountService.VerifyAccount(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *accountController) ResetPIN(ctx *gin.Context) {
	req := RequestJSON[dto.VerifyAccountRequest](ctx)
	res, err := c.accountService.ResetPIN(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *accountController) BlockAccount(ctx *gin.Context) {
	req := RequestJSON[dto.VerifyAccountRequest](ctx)
	res, err := c.accountService.BlockAccount(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}
