package controllers

import (
	"log"

	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/services"
	"abdanhafidz.com/go-boilerplate/utils"
	"github.com/gin-gonic/gin"
)

type PaymentCallbackController interface {
	HandleCallback(ctx *gin.Context)
}

type paymentCallbackController struct {
	paymentService services.PaymentService
	eventService   services.EventService
	academyService services.AcademyService
}

func NewPaymentCallbackController(
	paymentService services.PaymentService,
	eventService services.EventService,
	academyService services.AcademyService,
) PaymentCallbackController {
	return &paymentCallbackController{
		paymentService: paymentService,
		eventService:   eventService,
		academyService: academyService,
	}
}

// Handle Payment Callback godoc
// @Summary      Handle Xendit Payment Callback
// @Description  Receive and process payment status updates from Xendit
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        request  body      map[string]interface{}  true  "Xendit Callback Payload"
// @Success      200      {object}  dto.SuccessResponse[string]
// @Failure      400      {object}  dto.ErrorResponse
// @Router       /api/v1/payment/callback [post]
func (c *paymentCallbackController) HandleCallback(ctx *gin.Context) {
	// Xendit sends JSON payload
	// Basic structure for Invoice Callback:
	// { "id": "...", "external_id": "...", "status": "PAID", ... }
	var callbackData map[string]interface{}
	if err := ctx.ShouldBindJSON(&callbackData); err != nil {
		utils.ResponseFAILED(ctx, gin.H(nil), http_error.BAD_REQUEST_ERROR)
		return
	}

	log.Printf("Payment Callback Received: %+v", callbackData)

	status, ok := callbackData["status"].(string)
	if !ok {
		// Not a status update or unknown format
		ResponseJSON(ctx, gin.H(nil), "Ignored: No status", nil)
		return
	}

	invoiceId, _ := callbackData["id"].(string)

	if status == "PAID" || status == "SETTLED" {
		// Handle Event Payment
		// We need a method in Service to handle "ConfirmPayment" by InvoiceID
		// But existing services don't have it.
		// Let's add it to PaymentService? Yes.
		err := c.paymentService.ConfirmPayment(ctx.Request.Context(), invoiceId)
		if err != nil {
			log.Printf("Payment Confirmation Failed: %v", err)
			// Don't return error to Xendit if logic failed, but maybe we should?
			// Xendit expects 200 OK.
		}
	} else if status == "EXPIRED" {
		c.paymentService.ExpirePayment(ctx.Request.Context(), invoiceId)
	}

	// Always return 200 to Xendit
	ResponseJSON(ctx, gin.H(nil), "Callback Received", nil)
}
