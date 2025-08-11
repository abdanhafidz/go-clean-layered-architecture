package controller

import (
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
)

func RequestJSON[TRequest any](ctx *gin.Context) TRequest {
	request := new(TRequest)
	if err := ctx.ShouldBindBodyWithJSON(request); err != nil {
		utils.ResponseFAILED(ctx, request, http_error.BAD_REQUEST_ERROR)
		return *request
	} else {
		return *request
	}
}

func ResponseJSON[TResponse any, TMetaData any](ctx *gin.Context, metaData TMetaData, res TResponse, err error) {
	utils.SendResponse(ctx, metaData, res, err)
}
