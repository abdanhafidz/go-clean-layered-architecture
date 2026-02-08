package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventPaymentRepository interface {
	CreatePayment(ctx context.Context, eventPayment entity.EventPaymentTransaction) (entity.EventPaymentTransaction, error)
	GetPaymentById(ctx context.Context, id uuid.UUID) (entity.EventPaymentTransaction, error)
	UpdatePayment(ctx context.Context, eventPayment entity.EventPaymentTransaction) (entity.EventPaymentTransaction, error)
	GetPaymentByEventAndAccount(ctx context.Context, eventId uuid.UUID, accountId uuid.UUID) (entity.EventPaymentTransaction, error)
	GetByInvoiceId(ctx context.Context, invoiceId string) (entity.EventPaymentTransaction, error)
}

type eventPaymentRepository struct {
	db *gorm.DB
}

func NewEventPaymentRepository(db *gorm.DB) EventPaymentRepository {
	return &eventPaymentRepository{db: db}
}
func (r *eventPaymentRepository) CreatePayment(ctx context.Context, eventPayment entity.EventPaymentTransaction) (entity.EventPaymentTransaction, error) {
	err := r.db.WithContext(ctx).Create(&eventPayment).Error
	return eventPayment, err
}

func (r *eventPaymentRepository) GetPaymentById(ctx context.Context, id uuid.UUID) (entity.EventPaymentTransaction, error) {
	var payment entity.EventPaymentTransaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	return payment, err
}

func (r *eventPaymentRepository) UpdatePayment(ctx context.Context, eventPayment entity.EventPaymentTransaction) (entity.EventPaymentTransaction, error) {
	err := r.db.WithContext(ctx).Save(&eventPayment).Error
	return eventPayment, err
}

func (r *eventPaymentRepository) GetPaymentByEventAndAccount(ctx context.Context, eventId uuid.UUID, accountId uuid.UUID) (entity.EventPaymentTransaction, error) {
	var payment entity.EventPaymentTransaction
	err := r.db.WithContext(ctx).Where("event_id = ? AND account_id = ?", eventId, accountId).First(&payment).Error
	return payment, err
}

func (r *eventPaymentRepository) GetByInvoiceId(ctx context.Context, invoiceId string) (entity.EventPaymentTransaction, error) {
	var payment entity.EventPaymentTransaction
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceId).First(&payment).Error
	return payment, err
}
