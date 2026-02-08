package router

import (
	"abdanhafidz.com/go-boilerplate/provider"
	"github.com/gin-gonic/gin"
)

func PaymentCallbackRouter(r *gin.Engine, controller provider.ControllerProvider) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/payment/callback", controller.ProvidePaymentCallbackController().HandleCallback)
	}
}
