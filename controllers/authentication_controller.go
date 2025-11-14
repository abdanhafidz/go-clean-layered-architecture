package controllers

import (
	"abdanhafidz.com/go-clean-layered-architecture/models/dto"
	"abdanhafidz.com/go-clean-layered-architecture/services"
	"github.com/gin-gonic/gin"
)

type AuthenticationController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	ExternalAuth(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
}

type authenticationController struct {
	accountService      services.AccountService
	externalAuthService services.ExternalAuthService
}

func NewAuthenticationController(accountService services.AccountService, externalAuthService services.ExternalAuthService) AuthenticationController {
	return &authenticationController{
		accountService:      accountService,
		externalAuthService: externalAuthService,
	}
}
func (c *authenticationController) SignUp(ctx *gin.Context) {
	req := RequestJSON[dto.SignUpRequest](ctx)
	res, err := c.accountService.Create(ctx.Request.Context(), req.Name, req.Email, req.Username, req.Password)
	ResponseJSON(ctx, req, res, err)
}

func (c *authenticationController) SignIn(ctx *gin.Context) {
	req := RequestJSON[dto.SignInRequest](ctx)
	res, err := c.accountService.Validate(ctx, req.EmailorUsername, req.Password)
	ResponseJSON(ctx, req, res, err)
}

func (c *authenticationController) ChangePassword(ctx *gin.Context) {
	req := RequestJSON[dto.ChangePasswordRequest](ctx)
	accountId := ParseAccountId(ctx)
	res, err := c.accountService.ChangePassword(ctx.Request.Context(), accountId, req.OldPassword, req.NewPassword)
	ResponseJSON(ctx, req, res, err)
}

func (c *authenticationController) ExternalAuth(ctx *gin.Context) {
	req := RequestJSON[dto.ExternalAuthRequest](ctx)
	var (
		res dto.AuthenticatedUser
		err error
	)

	switch req.OauthProvider {
	case "google":
		res, err = c.externalAuthService.GoogleAuth(ctx.Request.Context(), req.OauthID)
	default:
		ResponseJSON(ctx, req, dto.AuthenticatedUser{}, nil)
		return
	}

	ResponseJSON(ctx, req, res, err)
}
