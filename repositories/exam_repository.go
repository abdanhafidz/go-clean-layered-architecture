package repositories

import (
	"context"

	entity "abdanhafidz.com/go-boilerplate/models/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ExamRepository interface {
	Create(ctx context.Context, e *entity.Exam) error
	Get(ctx context.Context, id uuid.UUID) (entity.Exam, error)
	GetBySlug(ctx context.Context, slug string) (entity.Exam, error)
	Update(ctx context.Context, e entity.Exam) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]entity.Exam, error)
	ListWithPagination(ctx context.Context, p entity.Pagination) ([]entity.Exam, int64, error)

	// Additional business need
	ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Exam, error)
}

type examRepository struct {
	db *gorm.DB
}

func NewExamRepository(db *gorm.DB) ExamRepository {
	return &examRepository{db}
}

// ================= CRUD =================

func (r *examRepository) Create(ctx context.Context, e *entity.Exam) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 1. Create Exam
	if err := tx.Create(&e).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Inject ExamId
	e.Configuration.ExamId = e.Id
	e.Proctoring.ExamId = e.Id

	// 3. Create Configuration
	if err := tx.Create(&e.Configuration).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 4. Create Proctoring
	if err := tx.Create(&e.Proctoring).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *examRepository) Get(ctx context.Context, id uuid.UUID) (entity.Exam, error) {
	var e entity.Exam
	err := r.db.WithContext(ctx).
		First(&e, "exam_id = ?", id).Error
	return e, err
}

func (r *examRepository) GetBySlug(ctx context.Context, slug string) (entity.Exam, error) {
	var e entity.Exam
	err := r.db.WithContext(ctx).
		Where("slug = ?", slug).
		First(&e).Error
	return e, err
}

func (r *examRepository) Update(ctx context.Context, e entity.Exam) error {
	return r.db.WithContext(ctx).
		Model(&entity.Exam{}).
		Where("exam_id = ?", e.Id).
		Updates(e).Error
}

func (r *examRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("exam_id = ?", id).
		Delete(&entity.Exam{}).Error
}

func (r *examRepository) List(ctx context.Context) ([]entity.Exam, error) {
	var list []entity.Exam
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *examRepository) ListWithPagination(ctx context.Context, p entity.Pagination) ([]entity.Exam, int64, error) {
	var list []entity.Exam
	var total int64

	db := r.db.WithContext(ctx).Model(&entity.Exam{})

	if p.Search != "" {
		db = db.Where("title ILIKE ?", "%"+p.Search+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if p.SortBy != "" {
		order := "ASC"
		if p.Order == "desc" {
			order = "DESC"
		}
		db = db.Order(p.SortBy + " " + order)
	} else {
		db = db.Order("created_at DESC")
	}

	err := db.Limit(p.Limit).Offset(p.Offset).Find(&list).Error
	return list, total, err
}

// =========== Business Specific ============

func (r *examRepository) ListByEvent(ctx context.Context, eventId uuid.UUID) ([]entity.Exam, error) {
	var exams []entity.Exam

	err := r.db.WithContext(ctx).
		Table("exam").
		Joins(`JOIN exam_event_assign ON exam_event_assign.exam_id = exam.id`).
		Where("exam_event_assign.event_id = ?", eventId).
		Find(&exams).Error

	return exams, err
}
