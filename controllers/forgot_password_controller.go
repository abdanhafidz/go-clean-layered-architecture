package controllers

import (
	"math/rand"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
	"github.com/gin-gonic/gin"
)

type ForgotPasswordController interface {
	Request(ctx *gin.Context)
	Reset(ctx *gin.Context)
}

type forgotPasswordController struct {
	forgotPasswordService services.ForgotPasswordService
}

func NewForgotPasswordController(forgotPasswordService services.ForgotPasswordService) ForgotPasswordController {
	return &forgotPasswordController{forgotPasswordService: forgotPasswordService}
}

// Request Forgot Password godoc
// @Summary      Request Password Reset
// @Description  Generate a password reset token and send it to the specified email address
// @Tags         Forgot Password
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ForgotPasswordRequest  true  "Forgot Password Request"
// @Success      200      {object}  dto.SuccessResponse[models.ForgotPassword]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/authentication/forgot-password [post]

func (c *forgotPasswordController) Request(ctx *gin.Context) {
	req := RequestJSON[dto.ForgotPasswordRequest](ctx)
	token := uint(rand.Intn(900000) + 100000)
	due := time.Now().Add(15 * time.Minute)
	res, err := c.forgotPasswordService.Request(ctx.Request.Context(), req.Email, token, due)
	ResponseJSON(ctx, req, res, err)
}

// Reset Forgot Password godoc
// @Summary      Reset Password
// @Description  Reset the user's password using the provided token
// @Tags         Forgot Password
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ResetPasswordRequest  true  "Reset Password Request"
// @Success      200      {object}  dto.SuccessResponse[any]
// @Failure      400      {object}  dto.ErrorResponse
func (c *forgotPasswordController) Reset(ctx *gin.Context) {
	req := RequestJSON[dto.ResetPasswordRequest](ctx)
	err := c.forgotPasswordService.Reset(ctx.Request.Context(), req.Token, req.NewPassword)
	ResponseJSON[any](ctx, req, gin.H{"status": "ok"}, err)
}
