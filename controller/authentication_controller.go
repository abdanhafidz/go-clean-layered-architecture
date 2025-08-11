package controller

import (
	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type AuthenticationController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authenticationController struct {
	authenticationService services.AuthenticationService
}

func NewAuthenticationController(authenticationService services.AuthenticationService) AuthenticationController {
	return &authenticationController{
		authenticationService: authenticationService,
	}
}

func (c *authenticationController) Register(ctx *gin.Context) {
	req := RequestJSON[dto.RegisterRequest](ctx)
	res, err := c.authenticationService.Register(ctx.Request.Context(), req)
	ResponseJSON(ctx, req, res, err)
}

func (c *authenticationController) Login(ctx *gin.Context) {
	req := RequestJSON[dto.LoginRequest](ctx)
	res, err := c.authenticationService.Login(ctx.Request.Context(), req.Email, req.Password)
	ResponseJSON(ctx, req, res, err)
}
