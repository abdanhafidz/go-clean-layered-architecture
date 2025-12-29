package controllers

import (
	"math/rand"
	"time"

	"abdanhafidz.com/go-boilerplate/models/dto"
	"abdanhafidz.com/go-boilerplate/services"
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

// Create Email Verification godoc
// @Summary      Create Email Verification Token
// @Description  Generate a verification token and send it to the specified email address
// @Tags         Email Verification
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateEmailVerificationRequest  true  "Create Email Verification Request"
// @Success      200      {object}  dto.SuccessResponse[models.EmailVerification]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/email/create-verification [post]
func (c *emailVerificationController) Create(ctx *gin.Context) {
	req := RequestJSON[dto.CreateEmailVerificationRequest](ctx)
	token := uint(rand.Intn(900000) + 100000)
	due := time.Now().Add(15 * time.Minute)
	res, err := c.emailVerificationService.CreateToken(ctx.Request.Context(), req.Email, token, due)
	ResponseJSON(ctx, req, res, err)
}

// Validate Email Verification godoc
// @Summary      Validate Email Verification Token
// @Description  Validate the provided verification token for the specified email address
// @Tags         Email Verification
// @Accept       json
// @Produce      json
// @Param        request  body      dto.ValidateVerifyEmailRequest  true  "Validate Verify Email Request"
// @Success      200      {object}  dto.SuccessResponse[any]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/email/verify [post]
func (c *emailVerificationController) Validate(ctx *gin.Context) {
	req := RequestJSON[dto.ValidateVerifyEmailRequest](ctx)
	err := c.emailVerificationService.VerifyToken(ctx.Request.Context(), req.Email, req.Token)
	ResponseJSON[any](ctx, req, gin.H{"status": "ok"}, err)
}

// Delete Email Verification godoc
// @Summary      Delete Email Verification Token
// @Description  Delete the verification token after successful validation
// @Tags         Email Verification
// @Accept       json
// @Produce      json
// @Param        request  body      dto.DeleteEmailVerificationRequest  true  "Delete Email Verification Request"
// @Success      200      {object}  dto.SuccessResponse[any]
// @Failure      400      {object}  dto.ErrorResponse
func (c *emailVerificationController) Delete(ctx *gin.Context) {

	req := RequestJSON[dto.DeleteEmailVerificationRequest](ctx)
	err := c.emailVerificationService.DeleteByToken(ctx.Request.Context(), req.Token)
	ResponseJSON[any](ctx, req, gin.H{"status": "ok"}, err)
}
