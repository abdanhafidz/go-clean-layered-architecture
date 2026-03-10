package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademyPaymentRepository interface {
	CreatePayment(ctx context.Context, academyPayment entity.AcademyPaymentTransaction) (entity.AcademyPaymentTransaction, error)
	GetPaymentById(ctx context.Context, id uuid.UUID) (entity.AcademyPaymentTransaction, error)
	GetPaymentByAcademyAndAccount(ctx context.Context, academyId uuid.UUID, accountId uuid.UUID) (entity.AcademyPaymentTransaction, error)
	UpdatePayment(ctx context.Context, academyPayment entity.AcademyPaymentTransaction) (entity.AcademyPaymentTransaction, error)
	GetByInvoiceId(ctx context.Context, invoiceId string) (entity.AcademyPaymentTransaction, error)
}

type academyPaymentRepository struct {
	db *gorm.DB
}

func NewAcaddemyPaymentRepository(db *gorm.DB) AcademyPaymentRepository {
	return &academyPaymentRepository{db: db}
}

func (r *academyPaymentRepository) CreatePayment(ctx context.Context, academyPayment entity.AcademyPaymentTransaction) (entity.AcademyPaymentTransaction, error) {
	err := r.db.WithContext(ctx).Create(&academyPayment).Error
	return academyPayment, err
}

func (r *academyPaymentRepository) GetPaymentByAcademyAndAccount(ctx context.Context, academyId uuid.UUID, accountId uuid.UUID) (entity.AcademyPaymentTransaction, error) {
	var payment entity.AcademyPaymentTransaction
	err := r.db.WithContext(ctx).Where("academy_id = ? AND account_id = ?", academyId, accountId).First(&payment).Error
	return payment, err
}
func (r *academyPaymentRepository) GetPaymentById(ctx context.Context, id uuid.UUID) (entity.AcademyPaymentTransaction, error) {
	var payment entity.AcademyPaymentTransaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&payment).Error
	return payment, err
}

func (r *academyPaymentRepository) UpdatePayment(ctx context.Context, academyPayment entity.AcademyPaymentTransaction) (entity.AcademyPaymentTransaction, error) {
	err := r.db.WithContext(ctx).Save(&academyPayment).Error
	return academyPayment, err
}

func (r *academyPaymentRepository) GetByInvoiceId(ctx context.Context, invoiceId string) (entity.AcademyPaymentTransaction, error) {
	var payment entity.AcademyPaymentTransaction
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceId).First(&payment).Error
	return payment, err
}
