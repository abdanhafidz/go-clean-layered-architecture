package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/xendit/xendit-go/v7"
)

type PaymentService interface {
	PaySomething(ctx context.Context, accountId uuid.UUID, amount float64)
	ConfirmPayment(ctx context.Context, paymentId string) error
	CancelPayment(ctx context.Context, paymentId string) error
	ExpirePayment(ctx context.Context, paymentId string) error
}

type paymentService struct {
	xenditClient *xendit.APIClient
}

func NewPaymentService(xenditClient *xendit.APIClient) PaymentService {
	return &paymentService{
		xenditClient: xenditClient,
	}
}

func (s *paymentService) PaySomething(ctx context.Context, accountId uuid.UUID, amount float64) {

}

func (s *paymentService) ConfirmPayment(ctx context.Context, paymentId string) error {
	return nil
}

func (s *paymentService) CancelPayment(ctx context.Context, paymentId string) error {
	return nil
}

func (s *paymentService) ExpirePayment(ctx context.Context, paymentId string) error {
	return nil
}
