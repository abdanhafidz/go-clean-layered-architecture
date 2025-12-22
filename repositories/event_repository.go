package repositories

import (
	"context"
	"strings"
	"time"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventsRepository interface {
	Create(ctx context.Context, ev entity.Events) (entity.Events, error)
	Update(ctx context.Context, ev entity.Events) (entity.Events, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.Events, error)
	GetBySlug(ctx context.Context, slug string) (entity.Events, error)
	GetByCode(ctx context.Context, code string) (entity.Events, error)
	GetAllPaginate(ctx context.Context, p entity.Pagination) ([]entity.Events, int64, error)
	ListPublic(ctx context.Context, p *entity.Pagination) ([]entity.Events, int64, error)
	ListVisible(ctx context.Context, accountId uuid.UUID, p *entity.Pagination) ([]entity.Events, int64, error)
}

type eventsRepository struct{ db *gorm.DB }

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

func (r *eventsRepository) GetAllPaginate(ctx context.Context, p entity.Pagination) ([]entity.Events, int64, error) {
	var list []entity.Events
	q := r.db.WithContext(ctx).Model(&entity.Events{})
	if s := strings.TrimSpace(p.Search); s != "" {
		s = strings.Trim(s, "\"'")
		s = strings.ToLower(s)
		like := "%" + s + "%"
		q = q.Where("LOWER(title) LIKE ? OR LOWER(slug) LIKE ? OR LOWER(event_code) LIKE ?", like, like, like)
	}
	col := strings.ToLower(strings.TrimSpace(p.SortBy))
	ord := strings.ToLower(strings.TrimSpace(p.Order))
	if col == "" {
		col = "title"
	}
	if ord != "desc" {
		ord = "asc"
	}
	switch col {
	case "title", "start_event", "end_event", "overview", "is_public", "created_at":
		q = q.Order(col + " " + ord)
	default:
		q = q.Order("title " + ord)
	}
	if p.Limit > 0 {
		q = q.Limit(p.Limit)
	}
	if p.Offset > 0 {
		q = q.Offset(p.Offset)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *eventsRepository) ListPublic(ctx context.Context, p *entity.Pagination) ([]entity.Events, int64, error) {
	var list []entity.Events
	q := r.db.WithContext(ctx).Model(&entity.Events{}).Where("is_public = ?", true)
	if p != nil {
		if s := strings.TrimSpace(p.Search); s != "" {
			s = strings.Trim(s, "\"'")
			s = strings.ToLower(s)
			like := "%" + s + "%"
			q = q.Where("LOWER(title) LIKE ? OR LOWER(slug) LIKE ? OR LOWER(event_code) LIKE ?", like, like, like)
		}
		col := strings.ToLower(strings.TrimSpace(p.SortBy))
		ord := strings.ToLower(strings.TrimSpace(p.Order))
		if col == "" {
			col = "title"
		}
		if ord != "desc" {
			ord = "asc"
		}
		switch col {
		case "title", "slug", "start_event", "end_event", "event_code", "overview", "is_public", "created_at":
			q = q.Order(col + " " + ord)
		case "createdat":
			q = q.Order("created_at " + ord)
		case "id_event", "id":
			q = q.Order("id " + ord)
		default:
			q = q.Order("title " + ord)
		}
		if p.Limit > 0 {
			q = q.Limit(p.Limit)
		}
		if p.Offset > 0 {
			q = q.Offset(p.Offset)
		}
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *eventsRepository) ListVisible(ctx context.Context, accountId uuid.UUID, p *entity.Pagination) ([]entity.Events, int64, error) {
	var list []entity.Events
	sub := r.db.WithContext(ctx).Model(&entity.EventAssign{}).Select("event_id").Where("account_id = ?", accountId)
	q := r.db.WithContext(ctx).Model(&entity.Events{}).Where("is_public = ?", true).Or("id IN (?)", sub)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if p != nil {
		if p.RegisterStatus != nil {
			switch *p.RegisterStatus {
			case 1:
				q = q.Where("academy.id IN (?)", sub)
			case 0:
				q = q.Where("academy.id NOT IN (?)", sub)
			}
		}

		if p.Status != nil {
			now := time.Now()
			switch *p.Status {
			case entity.EventStatusUpcoming:
				q = q.Where("start_event > ?", now)
			case entity.EventStatusOngoing:
				q = q.Where("start_event <= ? AND end_event >= ?", now, now)
			case entity.EventStatusEnded:
				q = q.Where("end_event < ?", now)
			}
		}

		if s := strings.TrimSpace(p.Search); s != "" {
			s = strings.Trim(s, "\"'")
			s = strings.ToLower(s)
			like := "%" + s + "%"
			q = q.Where("LOWER(title) LIKE ? OR LOWER(slug) LIKE ? OR LOWER(event_code) LIKE ?", like, like, like)
		}
		if err := q.Count(&total).Error; err != nil {
			return nil, 0, err
		}

		col := strings.ToLower(strings.TrimSpace(p.SortBy))
		ord := strings.ToLower(strings.TrimSpace(p.Order))
		if col == "" {
			col = "title"
		}
		if ord != "desc" {
			ord = "asc"
		}
		switch col {
		case "title", "slug", "start_event", "end_event", "overview", "created_at":
			q = q.Order(col + " " + ord)
		case "createdat":
			q = q.Order("created_at " + ord)
		case "id_event", "id":
			q = q.Order("id " + ord)
		default:
			q = q.Order("title " + ord)
		}
		if err := q.Count(&total).Error; err != nil {
			return nil, 0, err
		}

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
	if len(list) == 0 {
		return list, total, nil
	}

	eventIDs := make([]uuid.UUID, len(list))
	for i, ev := range list {
		eventIDs[i] = ev.Id
	}

	var assigns []entity.EventAssign
	if err := r.db.WithContext(ctx).
		Where("account_id = ?", accountId).
		Where("event_id IN ?", eventIDs).
		Find(&assigns).Error; err != nil {
		return nil, 0, err
	}
	assignedMap := make(map[uuid.UUID]bool)
	for _, a := range assigns {
		assignedMap[a.EventId] = true
	}

	for i := range list {
		if assignedMap[list[i].Id] {
			list[i].RegisterStatus = 1
		} else {
			list[i].RegisterStatus = 0
		}
	}

	return list, total, nil
}
