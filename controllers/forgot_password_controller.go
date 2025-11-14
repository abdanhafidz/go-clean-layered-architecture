package controllers

import (
	"math/rand"
	"time"

	"abdanhafidz.com/go-clean-layered-architecture/models/dto"
	"abdanhafidz.com/go-clean-layered-architecture/services"
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

func (c *forgotPasswordController) Request(ctx *gin.Context) {
	req := RequestJSON[dto.ForgotPasswordRequest](ctx)
	token := uint(rand.Intn(900000) + 100000)
	due := time.Now().Add(15 * time.Minute)
	res, err := c.forgotPasswordService.Request(ctx.Request.Context(), req.Email, token, due)
	ResponseJSON(ctx, req, res, err)
}

func (c *forgotPasswordController) Reset(ctx *gin.Context) {
	type resetReq struct {
		Token       uint   `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	req := RequestJSON[resetReq](ctx)
	err := c.forgotPasswordService.Reset(ctx.Request.Context(), req.Token, req.NewPassword)
	ResponseJSON[any](ctx, req, gin.H{"status": "ok"}, err)
}
