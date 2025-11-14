package controllers

import (
	http_error "abdanhafidz.com/go-clean-layered-architecture/models/error"
	"abdanhafidz.com/go-clean-layered-architecture/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func ResponseJSON[TResponse any, TMetaData any](ctx *gin.Context, metaData TMetaData, res TResponse, err error) {
	utils.SendResponse(ctx, metaData, res, err)
}
