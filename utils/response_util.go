package utils

import (
	"errors"

	"abdanhafidz.com/go-clean-layered-architecture/models/dto"
	http_error "abdanhafidz.com/go-clean-layered-architecture/models/error"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ResponseOK[Tdata any, TMetaData any](c *gin.Context, metaData TMetaData, data Tdata) {
	c.JSON(200, dto.SuccessResponse[Tdata]{
		Status:   "success",
		Data:     data,
		Message:  "Data retrieved Successfully!",
		MetaData: metaData,
	})
}

func ResponseFAILED[TMetaData any](c *gin.Context, metaData TMetaData, err error) {
	if errors.Is(err, http_error.BAD_REQUEST_ERROR) {
		c.JSON(400, dto.ErrorResponse{
			Status:   "error",
			Error:    err,
			Message:  "Invalid request format!",
			MetaData: metaData,
		})
		return
	} else if errors.Is(err, http_error.INTERNAL_SERVER_ERROR) {
		c.JSON(500, dto.ErrorResponse{
			Status:   "error",
			Error:    err,
			Message:  "Internal Server Error!",
			MetaData: metaData,
		})
		return
	} else if errors.Is(err, http_error.UNAUTHORIZED) {
		c.JSON(401, dto.ErrorResponse{
			Status:   "error",
			Error:    err,
			Message:  "Unauthorized, you don't have permission to access this service!",
			MetaData: metaData,
		})
		return
	} else if errors.Is(err, http_error.DATA_NOT_FOUND) || errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, dto.ErrorResponse{
			Status:   "error",
			Error:    err,
			Message:  "There is not data with given credential / given parameter!",
			MetaData: metaData,
		})
		return
	} else if errors.Is(err, http_error.TIMEOUT) {
		c.JSON(504, dto.ErrorResponse{
			Status:   "error",
			Error:    err,
			Message:  "Server took to long to respond!",
			MetaData: metaData,
		})
		return
	} else {
		c.JSON(405, dto.ErrorResponse{
			Status:   "error",
			Error:    err,
			Message:  err.Error(),
			MetaData: metaData,
		})
		return
	}

}

func SendResponse[Tdata any, TMetaData any](c *gin.Context, metaData TMetaData, data Tdata, err error) {
	if !c.IsAborted() {
		if err != nil {
			ResponseFAILED(c, metaData, err)
			c.Abort()
			return
		} else {
			ResponseOK(c, metaData, data)
			c.Abort()
			return
		}
	}

}
