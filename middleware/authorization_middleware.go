package middleware

import (
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	utils "abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthorizationMiddleware interface {
	AuthorizeUserToEvent(ctx *gin.Context)
}
type authorizationMiddleware struct {
	eventService services.EventService
}

func NewAuthorizationMiddleware(eventService services.EventService) AuthorizationMiddleware {
	return &authorizationMiddleware{
		eventService: eventService,
	}
}

func (m *authorizationMiddleware) AuthorizeUserToEvent(c *gin.Context) {

	eventSlug := c.Param("slug")
	accountId, exists := c.Get("account_id")
	if !exists {
		utils.ResponseFAILED(c, eventSlug, http_error.DATA_NOT_FOUND)
		c.Abort()
		return
	}

	err := m.eventService.AuthorizeUserToEvent(c.Request.Context(), eventSlug, accountId.(uuid.UUID))

	if err != nil {
		utils.ResponseFAILED(c, eventSlug, err)
		c.Abort()
		return
	}

	c.Next()

}
