package controllers

import (
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
	uuid "github.com/google/uuid"
)

func ParseAccountId(ctx *gin.Context) uuid.UUID {
	gaccountId, _ := ctx.Get("account_id")
	accountId, err := utils.ToUUID(gaccountId)
	if err != nil {
		ResponseJSON(ctx, gin.H{"account_id": accountId}, uuid.UUID{}, http_error.INVALID_TOKEN)
		return uuid.UUID{}
	}
	return accountId
}

func ParseUUID(ctx *gin.Context, attrName string) uuid.UUID {
	uuidRaw, _ := ctx.Get(attrName)
	uuidParsed, err := utils.ToUUID(uuidRaw)

	if err != nil {
		ResponseJSON(ctx, gin.H{"id": uuidParsed}, uuid.UUID{}, http_error.INVALID_TOKEN)
		return uuid.UUID{}
	}
	return uuidParsed
}
func RequestJSON[TRequest any](ctx *gin.Context) TRequest {
	var request TRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.ResponseFAILED(ctx, request, http_error.BAD_REQUEST_ERROR)
		ctx.Abort()
		return request
	} else {
		return request
	}
}

func RequestForm[TRequest any](ctx *gin.Context) TRequest {
	var request TRequest
	if err := ctx.ShouldBind(&request); err != nil {
		utils.ResponseFAILED(ctx, request, http_error.BAD_REQUEST_ERROR)
		ctx.Abort()
		return request
	}
	return request
}

func ResponseJSON[TResponse any, TMetaData any](ctx *gin.Context, metaData TMetaData, res TResponse, err error) {
	utils.SendResponse(ctx, metaData, res, err)
}
