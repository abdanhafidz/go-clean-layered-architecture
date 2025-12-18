package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthenticationController interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	ExternalAuth(ctx *gin.Context)
	ChangePassword(ctx *gin.Context)
	UpdateUserRole(ctx *gin.Context)
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

func (c *authenticationController) UpdateUserRole(ctx *gin.Context) {
	req := RequestJSON[dto.UpdateUserRoleRequest](ctx)
	targetAccountId, err := uuid.Parse(ctx.Param("account_id"))
	if err != nil {
		ResponseJSON(ctx, req, entity.Account{}, http_error.BAD_REQUEST_ERROR)
		return
	}

	acc, err := c.accountService.GetById(ctx.Request.Context(), targetAccountId)
	if err != nil {
		ResponseJSON(ctx, req, entity.Account{}, err)
		return
	}

	accountToUpdate := acc
	accountToUpdate.Role = req.Role

	res, err := c.accountService.Update(ctx.Request.Context(), accountToUpdate)
	ResponseJSON(ctx, req, res, err)
}
