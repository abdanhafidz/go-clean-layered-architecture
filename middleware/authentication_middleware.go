package middleware

import (
	"errors"
	"fmt"
	"strings"

	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	utils "abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
)

type AuthenticationMiddleware interface {
	VerifyAccount(ctx *gin.Context)
}
type authenticationMiddleware struct {
	jwtService services.JWTService
}

func NewAuthenticationMiddleware(jwtService services.JWTService) AuthenticationMiddleware {
	return &authenticationMiddleware{
		jwtService: jwtService,
	}
}
func (m *authenticationMiddleware) VerifyAccount(c *gin.Context) {

	authorizationBearer := c.Request.Header["Authorization"]

	if authorizationBearer != nil {
		token := strings.Split(authorizationBearer[0], " ")[1]
		claim, err := m.jwtService.ValidateToken(c.Request.Context(), token)

		if err != nil && errors.Is(err, http_error.INVALID_TOKEN) {
			utils.ResponseFAILED(c, claim, http_error.INVALID_TOKEN)
			c.Abort()
			return
		}
		fmt.Println("Claims:", claim)
		c.Set("account_id", claim.AccountId)
		c.Next()

	} else {
		utils.ResponseFAILED(c, "Empty Token", http_error.UNAUTHORIZED)
		return
	}

}
