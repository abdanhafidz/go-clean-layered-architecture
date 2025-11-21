package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit  int
	Offset int
}

type EventsRepository interface {
	Create(ctx context.Context, ev entity.Events) (entity.Events, error)
	Update(ctx context.Context, ev entity.Events) (entity.Events, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Events, error)
	GetBySlug(ctx context.Context, slug string) (entity.Events, error)
	GetByCode(ctx context.Context, code string) (entity.Events, error)
	GetAllPaginate(ctx context.Context, p Pagination) ([]entity.Events, int64, error)
	ListPublic(ctx context.Context, p *Pagination) ([]entity.Events, int64, error)
}

type eventsRepository struct {
	db *gorm.DB
}

func NewEventsRepository(db *gorm.DB) EventsRepository {
	return &eventsRepository{db: db}
}

func (r *eventsRepository) Create(ctx context.Context, ev entity.Events) (entity.Events, error) {
	if err := r.db.WithContext(ctx).Create(&ev).Error; err != nil {
		return entity.Events{}, err
	}
	return ev, nil
}

func (r *eventsRepository) Update(ctx context.Context, ev entity.Events) (entity.Events, error) {
	if err := r.db.WithContext(ctx).Save(&ev).Error; err != nil {
		return entity.Events{}, err
	}
	return ev, nil
}

func (r *eventsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Events{}, "id = ?", id).Error
}

func (r *eventsRepository) GetByID(ctx context.Context, id uuid.UUID) (entity.Events, error) {
	var ev entity.Events
	if err := r.db.WithContext(ctx).First(&ev, "id = ?", id).Error; err != nil {
		return entity.Events{}, err
	}
	return ev, nil
}

func (r *eventsRepository) GetBySlug(ctx context.Context, slug string) (entity.Events, error) {
	var ev entity.Events
	if err := r.db.WithContext(ctx).First(&ev, "slug = ?", slug).Error; err != nil {
		return entity.Events{}, err
	}
	return ev, nil
}

func (r *eventsRepository) GetByCode(ctx context.Context, code string) (entity.Events, error) {
	var ev entity.Events
	if err := r.db.WithContext(ctx).First(&ev, "event_code = ?", code).Error; err != nil {
		return entity.Events{}, err
	}
	return ev, nil
}

func (r *eventsRepository) GetAllPaginate(ctx context.Context, p Pagination) ([]entity.Events, int64, error) {
	var list []entity.Events
	q := r.db.WithContext(ctx).Model(&entity.Events{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if p.Limit > 0 {
		q = q.Limit(p.Limit)
	}
	if p.Offset > 0 {
		q = q.Offset(p.Offset)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *eventsRepository) ListPublic(ctx context.Context, p *Pagination) ([]entity.Events, int64, error) {
	var list []entity.Events
	q := r.db.WithContext(ctx).Model(&entity.Events{}).Where("is_public = ?", true)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if p != nil {
		if p.Limit > 0 {
			q = q.Limit(p.Limit)
		}
		if p.Offset > 0 {
			q = q.Offset(p.Offset)
		}
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
