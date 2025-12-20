package middleware

import (
	"strings"

	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	utils "abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthorizationMiddleware interface {
	AuthorizeUserToEvent(ctx *gin.Context)
	VerifyAdmin(ctx *gin.Context)
	VerifySuperAdmin(ctx *gin.Context)
}
type authorizationMiddleware struct {
	eventService services.EventService
	jwtService   services.JWTService
}

func NewAuthorizationMiddleware(eventService services.EventService, jwtService services.JWTService) AuthorizationMiddleware {
	return &authorizationMiddleware{
		eventService: eventService,
		jwtService:   jwtService,
	}
}

func (m *authorizationMiddleware) AuthorizeUserToEvent(c *gin.Context) {
	eventSlug := c.Param("slug")
	accountId, exists := c.Get("account_id")
	if !exists {
		utils.ResponseFAILED(c, eventSlug, http_error.NOT_FOUND_ERROR)
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

func (m *authorizationMiddleware) VerifyAdmin(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.ResponseFAILED(c, "Authorization header missing", http_error.UNAUTHORIZED)
		c.Abort()
		return
	}

	tokenString := strings.Split(authHeader, " ")[1]
	claims, err := m.jwtService.ValidateToken(c.Request.Context(), tokenString)
	if err != nil {
		utils.ResponseFAILED(c, "Invalid token", http_error.UNAUTHORIZED)
		c.Abort()
		return
	}

	if claims.Role != "admin" && claims.Role != "super_admin" {
		utils.ResponseFAILED(c, "Forbidden: Admin access required", http_error.FORBIDDEN_ERROR)
		c.Abort()
		return
	}

	c.Set("role", claims.Role)
	c.Set("account_id", claims.AccountId)
	c.Next()
}

func (m *authorizationMiddleware) VerifySuperAdmin(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.ResponseFAILED(c, "Authorization header missing", http_error.UNAUTHORIZED)
		c.Abort()
		return
	}

	tokenString := strings.Split(authHeader, " ")[1]
	claims, err := m.jwtService.ValidateToken(c.Request.Context(), tokenString)
	if err != nil {
		utils.ResponseFAILED(c, "Invalid token", http_error.UNAUTHORIZED)
		c.Abort()
		return
	}

	if claims.Role != "super_admin" {
		utils.ResponseFAILED(c, "Forbidden: Superadmin access required", http_error.FORBIDDEN_ERROR)
		c.Abort()
		return
	}

	c.Set("role", claims.Role)
	c.Set("account_id", claims.AccountId)
	c.Next()
}
