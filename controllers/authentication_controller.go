package controllers

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
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

// SignUp godoc
// @Summary      User Registration
// @Description  Register a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SignUpRequest  true  "Sign Up Request"
// @Success      200      {object}  dto.SuccessResponse[entity.Account]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/authentication/register [post]
func (c *authenticationController) SignUp(ctx *gin.Context) {
	req := RequestJSON[dto.SignUpRequest](ctx)
	res, err := c.accountService.Create(ctx.Request.Context(), req.Name, req.Email, req.Username, req.Password)
	ResponseJSON(ctx, req, res, err)
}

// SignIn godoc
// @Summary      User Login
// @Description  Authenticate user and obtain access token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      dto.SignInRequest  true  "Sign In Request"
// @Success      200      {object}  dto.SuccessResponse[dto.AuthenticatedUser]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/authentication/login [post]
func (c *authenticationController) SignIn(ctx *gin.Context) {
	req := RequestJSON[dto.SignInRequest](ctx)
	res, err := c.accountService.Validate(ctx, req.EmailorUsername, req.Password)
	ResponseJSON(ctx, req, res, err)
}

// ChangePassword godoc
// @Summary      Change User Password
// @Description  Change the password of the authenticated user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ChangePasswordRequest  true  "Change Password Request"
// @Success      200      {object}  dto.SuccessResponse[dto.AuthenticatedUser]
// @Failure      400      {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/authentication/change-password [put]

func (c *authenticationController) ChangePassword(ctx *gin.Context) {
	req := RequestJSON[dto.ChangePasswordRequest](ctx)
	accountId := ParseAccountId(ctx)
	res, err := c.accountService.ChangePassword(ctx.Request.Context(), accountId, req.OldPassword, req.NewPassword)
	ResponseJSON(ctx, req, res, err)
}

// ExternalAuth godoc
// @Summary      External Authentication
// @Description  Authenticate user using external OAuth provider
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ExternalAuthRequest  true  "External Auth Request"
// @Success      200      {object}  dto.SuccessResponse[dto.AuthenticatedUser]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/authentication/external-login [post]
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

// UpdateUserRole godoc
// @Summary      Update User Role
// @Description  Update the role of a user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        accountId  path      string                     true  "Account ID"
// @Param        request    body      dto.UpdateUserRoleRequest  true  "Update User Role Request"
// @Success      200        {object}  dto.SuccessResponse[entity.Account]
// @Failure      400        {object}  dto.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/admin/authentication/{account_id}/assign [put]

func (c *authenticationController) UpdateUserRole(ctx *gin.Context) {
	req := RequestJSON[dto.UpdateUserRoleRequest](ctx)
	accountId := ParseUUID(ctx, "accountId")
	acc, err := c.accountService.GetById(ctx.Request.Context(), accountId)
	if err != nil {
		ResponseJSON(ctx, req, entity.Account{}, err)
		return
	}

	accountToUpdate := acc
	accountToUpdate.Role = req.Role

	res, err := c.accountService.Update(ctx.Request.Context(), accountToUpdate)
	ResponseJSON(ctx, req, res, err)
}
