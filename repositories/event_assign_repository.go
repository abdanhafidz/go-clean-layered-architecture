package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventAssignRepository interface {
	Assign(ctx context.Context, assign entity.EventAssign) (entity.EventAssign, error)
	GetByEventAndAccount(ctx context.Context, eventId, accountId uuid.UUID) (entity.EventAssign, error)
	ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.EventAssign, error)
	ListByAccount(ctx context.Context, accountId uuid.UUID) ([]entity.EventAssign, error)
	Remove(ctx context.Context, id uuid.UUID) error
}

type eventAssignRepository struct {
	db *gorm.DB
}

func NewEventAssignRepository(db *gorm.DB) EventAssignRepository {
	return &eventAssignRepository{db: db}
}

func (r *eventAssignRepository) Assign(ctx context.Context, assign entity.EventAssign) (entity.EventAssign, error) {
	if err := r.db.WithContext(ctx).Create(&assign).Error; err != nil {
		return entity.EventAssign{}, err
	}
	return assign, nil
}

func (r *eventAssignRepository) GetByEventAndAccount(ctx context.Context, eventId, accountId uuid.UUID) (entity.EventAssign, error) {
	var rec entity.EventAssign
	if err := r.db.WithContext(ctx).
		Where("event_id = ? AND account_id = ?", eventId, accountId).
		First(&rec).Error; err != nil {
		return entity.EventAssign{}, err
	}
	return rec, nil
}

func (r *eventAssignRepository) ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.EventAssign, error) {
	var list []entity.EventAssign
	if err := r.db.WithContext(ctx).Where("event_id = ?", eventId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *eventAssignRepository) ListByAccount(ctx context.Context, accountId uuid.UUID) ([]entity.EventAssign, error) {
	var list []entity.EventAssign
	if err := r.db.WithContext(ctx).Where("account_id = ?", accountId).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *eventAssignRepository) Remove(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.EventAssign{}, "id = ?", id).Error
}
