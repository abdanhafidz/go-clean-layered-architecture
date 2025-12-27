package services

import (
	"context"
	"errors"
	"time"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	http_error "abdanhafidz.com/go-boilerplate/models/error"
	"abdanhafidz.com/go-boilerplate/repositories"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
	"gorm.io/gorm"
)

type PaymentService interface {
	PayEvent(ctx context.Context, accountId uuid.UUID, eventId uuid.UUID, amount float64) (entity.EventPaymentTransaction, error)
	PayAcademy(ctx context.Context, accountId uuid.UUID, academySlug uuid.UUID, amount float64) (entity.AcademyPaymentTransaction, error)
}

type paymentService struct {
	xenditClient       *xendit.APIClient
	eventPaymentRepo   repositories.EventPaymentRepository
	academyPaymentRepo repositories.AcademyPaymentRepository
}

func NewPaymentService(xenditClient *xendit.APIClient, eventPaymentRepo repositories.EventPaymentRepository, academyPaymentRepo repositories.AcademyPaymentRepository) PaymentService {
	return &paymentService{
		xenditClient:       xenditClient,
		eventPaymentRepo:   eventPaymentRepo,
		academyPaymentRepo: academyPaymentRepo,
	}
}

func (s *paymentService) PayEvent(ctx context.Context, accountId uuid.UUID, eventId uuid.UUID, amount float64) (entity.EventPaymentTransaction, error) {

	lastPayment, err := s.eventPaymentRepo.GetPaymentByEventAndAccount(ctx, eventId, accountId)

	if errors.Is(err, gorm.ErrRecordNotFound) || lastPayment.Status == entity.PaymentStatusExpired {
		externalId := "event-" + uuid.NewString()
		expiredAt := time.Now().Add(24 * time.Hour)

		InvoiceReq := *invoice.NewCreateInvoiceRequest(externalId, amount)

		InvoiceReq.SetDescription("Payment Event ID: " + eventId.String() + "For Account ID: " + accountId.String())
		InvoiceReq.SetCurrency("IDR")
		InvoiceReq.SetInvoiceDuration(float32(24 * 60 * 60)) // 24 jam

		invoiceResp, _, err := s.xenditClient.
			InvoiceApi.
			CreateInvoice(ctx).
			CreateInvoiceRequest(InvoiceReq).
			Execute()

		if err != nil {
			return entity.EventPaymentTransaction{}, err
		}

		lastPayment, errCreating := s.eventPaymentRepo.CreatePayment(ctx, entity.EventPaymentTransaction{
			EventId:       eventId,
			AccountId:     accountId,
			ExternalId:    externalId,
			InvoiceId:     invoiceResp.GetId(),
			InvoiceUrl:    invoiceResp.GetInvoiceUrl(),
			Amount:        amount,
			TransactionAt: time.Now(),
			ExpiredAt:     expiredAt,
			Status:        entity.PaymentStatusPending,
		})

		if errCreating != nil {
			return lastPayment, errCreating
		}

		return lastPayment, nil

	} else {
		invoiceData, _, err := s.xenditClient.InvoiceApi.GetInvoiceById(ctx, lastPayment.InvoiceId).Execute()

		if err != nil {
			return entity.EventPaymentTransaction{}, err
		}

		if invoiceData.Status == "PAID" {
			lastPayment.Status = entity.PaymentStatusPaid

			lastPayment, err := s.eventPaymentRepo.UpdatePayment(ctx, lastPayment)

			if err != nil {
				return entity.EventPaymentTransaction{}, err
			}

			return lastPayment, err

		} else if invoiceData.Status == "EXPIRED" {
			lastPayment.Status = entity.PaymentStatusExpired
			_, err := s.eventPaymentRepo.UpdatePayment(ctx, lastPayment)
			if err != nil {
				return entity.EventPaymentTransaction{}, err
			}

			return s.PayEvent(ctx, accountId, eventId, amount)

		} else if invoiceData.Status == "FAILED" {
			lastPayment.Status = entity.PaymentStatusFailed
			lastPayment, err := s.eventPaymentRepo.UpdatePayment(ctx, lastPayment)
			if err != nil {
				return entity.EventPaymentTransaction{}, err
			}

			return lastPayment, http_error.PAYMENT_FAILED
		}

		return lastPayment, nil
	}

}

func (s *paymentService) PayAcademy(ctx context.Context, accountId uuid.UUID, academyId uuid.UUID, amount float64) (entity.AcademyPaymentTransaction, error) {

	lastPayment, err := s.academyPaymentRepo.GetPaymentByAcademyAndAccount(ctx, accountId, academyId)

	if errors.Is(err, gorm.ErrRecordNotFound) || lastPayment.Status == entity.PaymentStatusExpired {
		externalId := "academy-" + uuid.NewString()
		expiredAt := time.Now().Add(24 * time.Hour)

		invoiceReq := *invoice.NewCreateInvoiceRequest(externalId, amount)

		invoiceReq.SetDescription("Academy ID: " + academyId.String() + "For Account ID: " + accountId.String())
		invoiceReq.SetCurrency("IDR")
		invoiceReq.SetInvoiceDuration(float32(24 * 60 * 60)) // 24 jam

		invoiceResp, _, err := s.xenditClient.
			InvoiceApi.
			CreateInvoice(ctx).
			CreateInvoiceRequest(invoiceReq).
			Execute()

		if err != nil {
			return entity.AcademyPaymentTransaction{}, err
		}

		lastPayment, errCreating := s.academyPaymentRepo.CreatePayment(ctx,
			entity.AcademyPaymentTransaction{
				AcademyId:     academyId,
				AccountId:     accountId,
				ExternalId:    externalId,
				InvoiceId:     invoiceResp.GetId(),
				InvoiceUrl:    invoiceResp.GetInvoiceUrl(),
				Amount:        amount,
				TransactionAt: time.Now(),
				ExpiredAt:     expiredAt,
				Status:        entity.PaymentStatusPending,
			},
		)

		if errCreating != nil {
			return lastPayment, errCreating
		}

		return lastPayment, nil

	} else {
		invoiceData, _, err := s.xenditClient.InvoiceApi.GetInvoiceById(ctx, lastPayment.InvoiceId).Execute()

		if err != nil {
			return entity.AcademyPaymentTransaction{}, err
		}

		if invoiceData.Status == "PAID" {
			lastPayment.Status = entity.PaymentStatusPaid

			lastPayment, err := s.academyPaymentRepo.UpdatePayment(ctx, lastPayment)

			if err != nil {
				return entity.AcademyPaymentTransaction{}, err
			}

			return lastPayment, err

		} else if invoiceData.Status == "EXPIRED" {
			lastPayment.Status = entity.PaymentStatusExpired
			_, err := s.academyPaymentRepo.UpdatePayment(ctx, lastPayment)
			if err != nil {
				return entity.AcademyPaymentTransaction{}, err
			}

			return s.PayAcademy(ctx, accountId, academyId, amount)

		} else if invoiceData.Status == "FAILED" {
			lastPayment.Status = entity.PaymentStatusFailed
			lastPayment, err := s.academyPaymentRepo.UpdatePayment(ctx, lastPayment)
			if err != nil {
				return entity.AcademyPaymentTransaction{}, err
			}

			return lastPayment, http_error.PAYMENT_FAILED
		}

		return lastPayment, nil
	}

}
