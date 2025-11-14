package controllers

import (
	"math/rand"
	"time"

	"abdanhafidz.com/go-clean-layered-architecture/models/dto"
	"abdanhafidz.com/go-clean-layered-architecture/services"
	"github.com/gin-gonic/gin"
)

type EmailVerificationController interface {
	Create(ctx *gin.Context)
	Validate(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type emailVerificationController struct {
	emailVerificationService services.EmailVerificationService
}

func NewEmailVerificationController(emailVerificationService services.EmailVerificationService) EmailVerificationController {
	return &emailVerificationController{emailVerificationService: emailVerificationService}
}

func (c *emailVerificationController) Create(ctx *gin.Context) {
	req := RequestJSON[dto.CreateEmailVerificationRequest](ctx)
	token := uint(rand.Intn(900000) + 100000)
	due := time.Now().Add(15 * time.Minute)
	res, err := c.emailVerificationService.CreateToken(ctx.Request.Context(), req.Email, token, due)
	ResponseJSON(ctx, req, res, err)
}

func (c *emailVerificationController) Validate(ctx *gin.Context) {
	req := RequestJSON[dto.ValidateVerifyEmailRequest](ctx)
	err := c.emailVerificationService.VerifyToken(ctx.Request.Context(), req.Email, req.Token)
	ResponseJSON[any](ctx, req, gin.H{"status": "ok"}, err)
}

func (c *emailVerificationController) Delete(ctx *gin.Context) {
	type delReq struct {
		Token uint `json:"token" binding:"required"`
	}
	req := RequestJSON[delReq](ctx)
	err := c.emailVerificationService.DeleteByToken(ctx.Request.Context(), req.Token)
	ResponseJSON[any](ctx, req, gin.H{"status": "ok"}, err)
}
