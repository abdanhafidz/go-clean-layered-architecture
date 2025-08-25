package controller

import (
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
)

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
